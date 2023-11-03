package memlog

import "strings"

// StringLog is used to write an internal list of
// strings to a MemLog[T] structure.
type StringLog struct {
	Buffer *MemLog[string]
}

// NewStringLog returns a StringLog initialized
// with a maximum of size entries.
func NewStringLog(size int) *StringLog {
	return &StringLog{
		Buffer: NewMemLog[string](size),
	}
}

// Write provides an implentation of the io.Writer
// interface that writes the output from the stream
// into a set of strings inside a MemLog buffer
func (s *StringLog) Write(p []byte) (n int, err error) {
	s.Buffer.Append(strings.Trim(string(p), "\r\n"))
	return len(p), nil
}
