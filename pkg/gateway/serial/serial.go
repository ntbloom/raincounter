// Package serial controls serial communication over USB
package serial

import (
	"os"
	"sync"
	"time"

	exitcodes2 "github.com/ntbloom/raincounter/pkg/common/exitcodes"

	messenger2 "github.com/ntbloom/raincounter/pkg/gateway/messenger"
	tlv2 "github.com/ntbloom/raincounter/pkg/gateway/tlv"

	"github.com/sirupsen/logrus"
)

// Serial communicates with a serial port
type Serial struct {
	port            string                // file descriptor of port
	maxPacketLen    int                   // how long you expect the packet to be
	timeout         time.Duration         // how long to wait for enumration
	data            []byte                // where the data lives
	file            *os.File              // file descriptor for the port
	kill            chan struct{}         // send a message to kill the serial loop
	messageReceived chan struct{}         // channel for waiting for message on serial port
	Messenger       *messenger2.Messenger // messenger object
	sync.Mutex
}

// NewConnection creates a new serial connection with a unix filename
func NewConnection(port string, maxPacketLen int, timeout time.Duration, msgr *messenger2.Messenger) (*Serial, error) {
	checkPortStatus(port, timeout)
	logrus.Infof("opening connection on `%s`", port)
	var data []byte

	// attempt to connect until timeout is exhausted

	file, err := os.Open(port)
	if err != nil {
		logrus.Errorf("problem opening port `%s`: %s", port, err)
		return nil, err
	}

	uart := &Serial{
		port,
		maxPacketLen,
		timeout,
		data,
		file,
		make(chan struct{}, 1),
		make(chan struct{}, 1),
		msgr,
		sync.Mutex{},
	}

	return uart, nil
}

// Start runs the main loop of listening on the serial port
func (serial *Serial) Start() {
	checkPortStatus(serial.port, serial.timeout)

	// send the first waiting command
	go serial.waitForMessage()

	// wait for next message or the kill signal
	for {
		select {
		case <-serial.kill:
			serial.close()
			return
		case <-serial.messageReceived:
			go serial.waitForMessage()
		}
	}
}

// Stop stops the main loop listening on the serial port
func (serial *Serial) Stop() {
	serial.kill <- struct{}{}
}

// waits for a message, then updates the main loop when it arrives
func (serial *Serial) waitForMessage() {
	serial.Lock()
	defer serial.Unlock()

	logrus.Tracef("waiting to read contents of `%s`", serial.port)
	packet := make([]byte, serial.maxPacketLen)

	_, err := serial.file.Read(packet)
	if err != nil {
		// connection to file was lost, attempt reconnection
		logrus.Infof("connection lost, attempting reconnection")
		checkPortStatus(serial.port, serial.timeout)
		_ = serial.reopenConnection()
	}
	logrus.Trace("a serial message arrived")
	serial.messageReceived <- struct{}{}

	tlvPacket, err := tlv2.NewTLV(packet)
	if err != nil {
		logrus.Errorf("unexpected TLV packet: %s", err)
	}
	msg, err := serial.Messenger.NewMessage(tlvPacket)
	if err != nil {
		logrus.Errorf("bad tlv packet, ignoring: %s", err)
	}
	serial.Messenger.Data <- msg
}

// reopens the serial connection if it gets broken
func (serial *Serial) reopenConnection() error {
	file, err := os.Open(serial.port)
	if err != nil {
		logrus.Debugf("port `%s` temporarily down: %s", serial.port, err)
		return err
	}
	serial.file = file
	return nil
}

// closes the serial port, to be used at the end of the program
func (serial *Serial) close() {
	logrus.Infof("closing serial port `%s`", serial.port)
	err := serial.file.Close()
	if err != nil {
		logrus.Errorf("problem closing `%s`: %s", serial.port, err)
	}
}

// make sure we can access the serial port else fail
func checkPortStatus(port string, timeout time.Duration) {
	logrus.Debugf("checking if `%s` exists", port)
	start := time.Now()
	for {
		_, err := os.Stat(port)
		if err == nil {
			logrus.Debugf("found port `%s`", port)
			return
		}
		logrus.Tracef("file `%s` doesn't exist on first look, re-checking for %s", port, timeout)
		if time.Since(start).Milliseconds() > timeout.Milliseconds() {
			handlePortFailure(port)
			return
		}
	}
}

// handles what to do when sensor is unresponsive
func handlePortFailure(port string) {
	logrus.Fatalf("unable to locate sensor at `%s`", port)

	// for now...
	os.Exit(exitcodes2.SerialPortNotFound)
}
