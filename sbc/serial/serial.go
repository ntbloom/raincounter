// Package serial controls serial communication over USB
package serial

import (
	"os"
	"sync"
	"time"

	"github.com/ntbloom/raincounter/sbc/messenger"

	"github.com/ntbloom/raincounter/common/exitcodes"
	"github.com/ntbloom/raincounter/sbc/tlv"

	"github.com/sirupsen/logrus"
)

// Serial communicates with a serial port
type Serial struct {
	port              string               // file descriptor of port
	maxPacketLen      int                  //how long you expect the packet to be
	timeout           time.Duration        // how long to wait for enumration
	data              []byte               //where the data lives
	file              *os.File             // file descriptor for the port
	Kill              chan bool            // send a message to kill the serial loop
	waitingForMessage chan bool            // flag for waiting for message on serial port
	Messenger         *messenger.Messenger //messenger object
	sync.RWMutex
}

// NewConnection creates a new serial connection with a unix filename
func NewConnection(port string, maxPacketLen int, timeout time.Duration, msgr *messenger.Messenger) (*Serial, error) {
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
		make(chan bool, 1),
		make(chan bool, 1),
		msgr,
		sync.RWMutex{},
	}

	return uart, nil
}

// Close closes the serial connection
func (serial *Serial) Close() {
	logrus.Infof("closing serial port `%s`", serial.port)
	err := serial.file.Close()
	if err != nil {
		logrus.Errorf("problem closing `%s`: %s", serial.port, err)
	}
}

// Loop reads the file contents
func (serial *Serial) Loop() {
	checkPortStatus(serial.port, serial.timeout)

	// send the waiting loop
	go serial.waitForMessage()
	logrus.Tracef("reading contents of `%s`", serial.port)
	for {
		select {
		case <-serial.Kill:
			serial.Close()
			return
		case <-serial.waitingForMessage:
			go serial.waitForMessage()
		}
	}
}

func (serial *Serial) waitForMessage() {
	serial.Lock()
	defer serial.Unlock()

	packet := make([]byte, serial.maxPacketLen)
	_, err := serial.file.Read(packet)
	serial.waitingForMessage <- true
	if err != nil {
		// connection to file was lost, attempt reconnection
		logrus.Infof("connection lost, attempting reconnection")
		checkPortStatus(serial.port, serial.timeout)
		_ = serial.reopenConnection()
	}

	tlvPacket, err := tlv.NewTLV(packet)
	if err != nil {
		logrus.Errorf("unexpected TLV packet: %s", err)
	}
	msg, err := serial.Messenger.NewMessage(tlvPacket)
	if err != nil {
		logrus.Errorf("bad tlv packet, ignoring: %s", err)
	}
	serial.Messenger.Data <- msg
}

// HandlePortFailure handles what to do when sensor is unresponsive
func HandlePortFailure(port string) {
	logrus.Fatalf("unable to locate sensor at `%s`", port)

	// for now...
	os.Exit(exitcodes.SerialPortNotFound)
}

// checks if a port is open
// doesn't use a *Serial receiver since we use it before creating *Serial object
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
			HandlePortFailure(port)
			return
		}
	}
}

func (serial *Serial) reopenConnection() error {
	file, err := os.Open(serial.port)
	if err != nil {
		logrus.Debugf("port `%s` temporarily down: %s", serial.port, err)
		return err
	}
	serial.file = file
	return nil
}
