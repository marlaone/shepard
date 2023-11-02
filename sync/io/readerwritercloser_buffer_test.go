package io_test

import (
	"sync"
	"testing"
	"time"

	"github.com/marlaone/shepard/collections/slice"
	"github.com/marlaone/shepard/sync/io"
	"github.com/stretchr/testify/assert"
)

func TestReaderWriterCloser_Write(t *testing.T) {
	rwc := io.NewReaderWriterCloser()
	res := rwc.Write(slice.Init[byte](1, 2, 3))

	assert.Equal(t, 3, res.Unwrap())
}

func TestReaderWriterCloser_Read(t *testing.T) {

	rwc := io.NewReaderWriterCloser()
	rwc.Write(slice.Init[byte](1, 2, 3))

	buf := slice.WithCapacity[byte](3)

	res := rwc.Read(&buf)

	assert.Equal(t, 3, res.Unwrap())
}

func TestReaderWriterCloser_Read_Async_Write(t *testing.T) {
	var wg sync.WaitGroup

	rwc := io.NewReaderWriterCloser()

	buf := slice.WithCapacity[byte](3)

	go func() {
		time.Sleep(200 * time.Millisecond)
		rwc.Write(slice.Init[byte](1, 2, 3))
		wg.Done()
	}()
	wg.Add(1)

	res := rwc.Read(&buf)

	assert.Equal(t, 0, res.Unwrap())

	wg.Wait()

	res = rwc.Read(&buf)
	assert.Equal(t, 3, res.Unwrap())

	assert.Equal(t, slice.Init[byte](1, 2, 3), buf)
}

func TestReaderWriterCloser_Close(t *testing.T) {
	rwc := io.NewReaderWriterCloser()
	rwc.Close()
	s := slice.Init[byte](1, 2, 3)
	res := rwc.Write(s)

	assert.Equal(t, true, res.IsErr())
	assert.Equal(t, "can't write to closed buffer", res.UnwrapErr().Error())
}
