package memlog

import (
	"container/list"
	"sync"
)

const (
	allElements = -1
)

// memLog is a bounded linked list that is intended
// used as a mechanism for logging information
// in memory.  The log has a fixed length and
// supports automatically removing older entries
// as new entries are added to prevent unbounded growth.
//
// memLog is useful in microservice applications where
// it might be useful to maintain a subset of a log
// in memory for diagnostics.  For instance the last
// 1,000 lines of the output log could be saved in memory
// and made available via a restful API to simplify
// diagnostics.
//
//   - Storing a subset of output logs in memory for
//     easy access during debugging or diagnostics.
//
//   - Storing a short history of results from an operation
//     that can be reviewed at runtime.
//
// memLog is thread-safe
type memLog[T any] struct {
	lst    list.List
	max    int
	locker sync.Mutex
}

// NewMemLog returns a new, initialized instance of memlog
// that will not grow beyond the specified number of
// entries.  Once the log reaches the maximum number of
// entries, as new entries are added, the oldest entries
// are removed.
func New[T any](maxEntries int) *memLog[T] {
	return &memLog[T]{
		max: maxEntries,
	}
}

// Len returns the number of elements in
// the log
func (m *memLog[T]) Len() int {
	m.locker.Lock()
	defer m.locker.Unlock()
	return m.lst.Len()
}

// Append will add item to the log.  If the
// log has reached its maximum size the the oldest
// entry will be removed to make room for the new entry.
func (m *memLog[T]) Append(item T) {
	m.locker.Lock()
	defer m.locker.Unlock()

	m.lst.PushBack(item)
	if m.lst.Len() > m.max {
		m.lst.Remove(m.lst.Front())
	}
}

// Slice returns the contents of the log as a slice.
// The slice is ordered from oldest item to the newest
func (m *memLog[T]) Slice() (slice []T) {
	return m.SliceN(allElements)
}

// Clear will clear the current contents of the memLog
func (m *memLog[T]) Clear() {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.lst.Init()
}

// SliceN returns the last 'N' items
// from the log.
// The slice is ordered from oldest item to the newest
func (m *memLog[T]) SliceN(n int) (slice []T) {
	m.locker.Lock()
	defer m.locker.Unlock()

	len := m.lst.Len()

	if n <= allElements || n > len {
		n = len
	}

	return m.toSlice(n)
}

// toSlice creates a slice of the last 'n' elements
// of the log.
func (m *memLog[T]) toSlice(n int) (slice []T) {
	slice = make([]T, n)
	idx := n - 1

	// Walk the list 'backward', filling in the slice
	// from the last element to the zero element.  This
	// is more efficient than searching 'forward' when n < m.lst.Len()
	for e := m.lst.Back(); e != nil && idx >= 0; e = e.Prev() {
		slice[idx] = e.Value.(T)
		idx--
	}

	return slice
}

// toSlice will copy a range of elements in the linked
// list to a slice
func (m *memLog[T]) toSlicex(n int, len int) (slice []T) {
	first := len - n
	slice = make([]T, n)
	ptr := 0
	item := m.lst.Front()

	for i := 0; i < len; i++ {
		if i < first {
			item = item.Next()
			continue
		}
		slice[ptr] = item.Value.(T)
		item = item.Next()
		ptr++
	}

	return slice
}
