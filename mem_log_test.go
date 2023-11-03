package memlog

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func Test_memlog_get_slice_when_empty(t *testing.T) {
	// given a memlog
	log := NewMemLog[string](10)

	// when a slice is created from an empty log
	slice := log.Slice()

	// then
	assert.Zero(t, len(slice))
	assert.Zero(t, log.Len())
}

func Test_memlog_get_length_when_size_is_less_than_max(t *testing.T) {
	// given a memlog
	log := NewMemLog[string](10)

	// when fewer than 'max' entries are added to the log
	log.Append("item #1")

	// then the length should match the expected number of entries
	assert.Equal(t, 1, log.Len())
}

func Test_memlog_get_length_when_log_is_at_max(t *testing.T) {
	// given a memlog
	log := NewMemLog[string](2)

	// when more than 'max' entries have been added
	log.Append("item #1")
	log.Append("item #2")
	log.Append("item #3")
	log.Append("item #4")

	// then the length should match the expected number of entries
	assert.Equal(t, 2, log.Len())
	assert.Equal(t, "item #3", log.Slice()[0])
}

func Test_memlog_get_last_n_entries(t *testing.T) {
	// given a memlog
	max := 20
	log := NewMemLog[string](max)

	// when more than 'max' entries have been added
	for i := 0; i < max; i++ {
		log.Append(fmt.Sprintf("item #%d", i))
	}

	// then the length should match the expected number of entries
	assert.Equal(t, "item #18", log.SliceN(2)[0])
}

func Test_memlog_list_memory(t *testing.T) {
	PrintMemUsage()

	size := 100000

	var l *MemLog[string]

	// given a memlog
	l = NewMemLog[string](size)

	// when an item is appended
	for i := 0; i < size; i++ {
		l.Append(fmt.Sprintf("this is a sample log entry that is probaby pretty typical in length %d", i))
	}

	// then
	slice := l.Slice()
	PrintMemUsage()
	//assert.False(t, true)
	assert.Equal(t, size, len(slice))
}

func Test_memlog_list_memory_pointers(t *testing.T) {
	PrintMemUsage()

	size := 100000

	// given a memlog
	l := NewMemLog[*string](size)

	// when an item is appended
	for i := 0; i < size; i++ {
		msg := fmt.Sprintf("this is a sample log entry that is probaby pretty typical in length %d", i)
		l.Append(&msg)
	}

	// then
	slice := l.Slice()
	PrintMemUsage()
	//assert.False(t, true)
	assert.Equal(t, size, len(slice))
}

func Benchmark_memlog_list_build_list(b *testing.B) {
	size := b.N

	// given a memlog
	l := NewMemLog[string](size)

	// when an item is appended
	for i := 0; i < b.N; i++ {
		msg := fmt.Sprintf("this is a sample log entry that is probaby pretty typical in length %d", i)
		l.Append(msg)
	}
}

func Benchmark_memlog_list_build_list_pointers(b *testing.B) {
	size := b.N

	// given a memlog
	l := NewMemLog[*string](size)

	// when an item is appended
	for i := 0; i < b.N; i++ {
		msg := fmt.Sprintf("this is a sample log entry that is probaby pretty typical in length %d", i)
		l.Append(&msg)
	}
}

func Benchmark_memlog_list_get_slice(b *testing.B) {
	// given a memlog
	l := NewMemLog[string](1000)

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("this is a sample log entry that is probaby pretty typical in length %d", i)
		l.Append(msg)
		slice := l.Slice()
		_ = slice
	}

	// when an item is appended
	for i := 0; i < b.N; i++ {
		slice := l.Slice()
		_ = slice
	}
}

func Benchmark_memlog_list_get_slice_pointers(b *testing.B) {
	// given a memlog
	l := NewMemLog[*string](1000)

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("this is a sample log entry that is probaby pretty typical in length %d", i)
		l.Append(&msg)
		slice := l.Slice()
		_ = slice
	}

	// when an item is appended
	for i := 0; i < b.N; i++ {
		slice := l.Slice()
		_ = slice
	}
}
