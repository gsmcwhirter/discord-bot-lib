package util

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/rs/xid"
)

// ReadBody TODOC
func ReadBody(resp io.Reader, estSize int) (bytesRead int, data []byte, err error) {
	data = make([]byte, estSize)
	buffer := make([]byte, estSize)
	var readSize int
	for err == nil {
		readSize, err = resp.Read(buffer)

		if readSize+bytesRead > estSize {
			newData := make([]byte, len(data)*2)
			for i := 0; i < len(data); i++ {
				newData[i] = data[i]
			}
			data = newData
		}

		for i := 0; i < readSize; i++ {
			data[bytesRead+i] = buffer[i]
		}
		bytesRead += readSize

		if err != nil {
			break
		}
	}

	if err == io.EOF {
		err = nil
	}

	data = data[:bytesRead]

	return
}

// CheckDefer TODOC
func CheckDefer(fs ...func() error) {
	for i := len(fs) - 1; i >= 0; i-- {
		if err := fs[i](); err != nil {
			fmt.Fprintf(os.Stderr, "Error in defer: %s\n", err) // nolint: errcheck,gas
		}
	}
}

// ContextKey TODOC
type ContextKey string

// NewRequestContext TODOC
func NewRequestContext() context.Context {
	return NewRequestContextFrom(context.Background())
}

// NewRequestContextFrom TODOC
func NewRequestContextFrom(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKey("request_id"), GenerateRequestID())
}

// GenerateRequestID TODOC
func GenerateRequestID() string {
	return xid.New().String()
}
