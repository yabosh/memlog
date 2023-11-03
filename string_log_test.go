package memlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_write_string_log_without_crlf(t *testing.T) {
	sl := NewStringLog(100)
	sl.Write([]byte("Test message"))
	assert.Equal(t, "Test message", sl.Buffer.Slice()[0])
}

func Test_write_string_log_with_lf(t *testing.T) {
	sl := NewStringLog(100)
	sl.Write([]byte("Test message\n"))
	assert.Equal(t, "Test message", sl.Buffer.Slice()[0])
}

func Test_write_string_log_with_cr(t *testing.T) {
	sl := NewStringLog(100)
	sl.Write([]byte("Test message\r"))
	assert.Equal(t, "Test message", sl.Buffer.Slice()[0])
}

func Test_write_string_log_multiple_lines(t *testing.T) {
	sl := NewStringLog(100)
	sl.Write([]byte("Test message 1\r"))
	sl.Write([]byte("Test message 2\r"))
	sl.Write([]byte("Test message 3\r"))
	assert.Equal(t, "Test message 1", sl.Buffer.Slice()[0])
	assert.Equal(t, "Test message 2", sl.Buffer.Slice()[1])
	assert.Equal(t, "Test message 3", sl.Buffer.Slice()[2])
}
