package transformer

import (
	"fmt"
	"strings"
)

// ReversibleStringBuilder is a string builder that allows the user to reverse
// a certain number of write instructions issued to it.
type ReversibleStringBuilder struct {
	builder      *strings.Builder
	buffer       []string
	internalSize int
}

// Init initializes internal fields within the string builder.
func (rsb *ReversibleStringBuilder) Init() {
	rsb.buffer = make([]string, 0)
	rsb.builder = &strings.Builder{}
}

// Len returns the length of the string in the reversible buffer plus the length
// of the string in the builder.
func (rsb *ReversibleStringBuilder) Len() int {
	return rsb.internalSize + rsb.builder.Len()
}

// Grow tells the internal string builder to grow by n characters.
func (rsb *ReversibleStringBuilder) Grow(n int) {
	rsb.builder.Grow(n)
}

// WriteString writes to the reversible buffer without writing to the string builder.
// Use Flush() commit changes to the string builder.
func (rsb *ReversibleStringBuilder) WriteString(s string) (int, error) {
	rsb.buffer = append(rsb.buffer, s)
	rsb.internalSize += len(s)

	return rsb.internalSize + rsb.builder.Len(), nil
}

// Reverse undos the last n WriteString operations. It can undo as many WriteString calls
// since the last call to Flush(). Passing clears the buffers.
func (rsb *ReversibleStringBuilder) Reverse(n int) {
	if n == -1 {
		rsb.buffer = rsb.buffer[:0]
		rsb.internalSize = 0

		return
	}

	numElements := len(rsb.buffer) - 1
	for i := numElements; i > numElements-n; i-- {
		rsb.internalSize -= len(rsb.buffer[i])
	}

	rsb.buffer = rsb.buffer[:len(rsb.buffer)-n]
}

// Flush writes everything in the buffer to the string builder. All potentially reversible
// operations will be saved to the string builer and can no longer be reversed.
func (rsb *ReversibleStringBuilder) Flush() error {
	for _, s := range rsb.buffer {
		_, err := rsb.builder.WriteString(s)
		if err != nil {
			return err
		}
	}

	rsb.buffer = rsb.buffer[:0]
	rsb.internalSize = 0

	return nil
}

// String flushes the buffer and returns the stored string.
func (rsb *ReversibleStringBuilder) String() string {
	if err := rsb.Flush(); err != nil {
		fmt.Println(err)
		return ""
	}

	return rsb.builder.String()
}
