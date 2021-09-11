// Package exitcodes defines error code constants for when the application fails
package exitcodes

const (
	// SerialPortNotFound errors when we can't connect to the database
	SerialPortNotFound = 1

	// TLSError is when we can't connect to mqtt using TLS certs
	TLSError = 2

	// PostgresqlConnectionError happens when we can't connect to postgres
	PostgresqlConnectionError = 3
)
