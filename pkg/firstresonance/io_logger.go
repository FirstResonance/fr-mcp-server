package firstresonance

import (
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

// IOLogger wraps an io.Reader and io.Writer to log all reads and writes
type IOLogger struct {
	reader io.Reader
	writer io.Writer
	logger *log.Logger
}

// NewIOLogger creates a new IOLogger
func NewIOLogger(reader io.Reader, writer io.Writer, logger *log.Logger) *IOLogger {
	return &IOLogger{
		reader: reader,
		writer: writer,
		logger: logger,
	}
}

// Read implements io.Reader
func (l *IOLogger) Read(p []byte) (n int, err error) {
	n, err = l.reader.Read(p)
	if n > 0 {
		l.logger.WithFields(log.Fields{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"type":      "read",
			"size":      n,
		}).Debug(string(p[:n]))
	}
	return
}

// Write implements io.Writer
func (l *IOLogger) Write(p []byte) (n int, err error) {
	n, err = l.writer.Write(p)
	if n > 0 {
		l.logger.WithFields(log.Fields{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"type":      "write",
			"size":      n,
		}).Debug(string(p[:n]))
	}
	return
}
