package csv

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"sync"
)

// FileConfig provides configuration options for file-based CSV processing
type FileConfig struct {
	Delimiter byte
	HasHeader bool
	Parallel  bool // Enable parallel processing
	Workers   int  // Number of workers (default: runtime.NumCPU())
}

// chunk represents a file chunk for parallel processing
type chunk struct {
	Start int64
	End   int64
}


// createChunks divides the file into chunks for parallel processing
func createChunks(fileSize int64, numWorkers int) []chunk {
	chunks := make([]chunk, numWorkers)
	chunkSize := fileSize / int64(numWorkers)

	for i := 0; i < numWorkers; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = fileSize
		}
		chunks[i] = chunk{Start: start, End: end}
	}

	return chunks
}

// ProcessFile reads a CSV file and calls fn for each deserialized record.
// If Parallel is true, fn may be called concurrently from multiple goroutines.
// Returns error if file cannot be read or parsing fails.
func ProcessFile[T any](path string, config FileConfig, fn func(T)) error {
	return ProcessFileWithError(path, config, func(record T) error {
		fn(record)
		return nil
	})
}

// ProcessFileWithError allows the callback to return an error to stop processing.
func ProcessFileWithError[T any](path string, config FileConfig, fn func(T) error) error {
	if config.Parallel {
		return processFileParallel[T](path, config, fn)
	}
	return processFileSequential[T](path, config, fn)
}

func processFileSequential[T any](path string, config FileConfig, fn func(T) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := NewReaderWithConfig[T](file, ReaderConfig{
		Delimiter: config.Delimiter,
		HasHeader: config.HasHeader,
	})

	for {
		result := reader.Deserialize()
		if result.IsErr() {
			if result.UnwrapErr() == io.EOF {
				return nil
			}
			return result.UnwrapErr()
		}
		if err := fn(result.Unwrap()); err != nil {
			return err
		}
	}
}

func processFileParallel[T any](path string, config FileConfig, fn func(T) error) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	workers := config.Workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	chunks := createChunks(fileInfo.Size(), workers)

	var wg sync.WaitGroup
	errChan := make(chan error, workers)

	for _, c := range chunks {
		wg.Add(1)
		go func(start, end int64) {
			defer wg.Done()
			if err := processChunk[T](path, start, end, config, fn); err != nil {
				errChan <- err
			}
		}(c.Start, c.End)
	}

	wg.Wait()
	close(errChan)

	// Return first error if any
	for err := range errChan {
		return err
	}
	return nil
}

func processChunk[T any](path string, start, end int64, config FileConfig, fn func(T) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Seek to start position and track where we are
	currentPos := start
	if start > 0 {
		// Check if we're starting right after a newline (at a line boundary)
		// by reading the byte before our start position
		file.Seek(start-1, io.SeekStart)
		prevByte := make([]byte, 1)
		_, err := file.Read(prevByte)
		if err != nil {
			return err
		}

		if prevByte[0] != '\n' {
			// We're in the middle of a line - skip to the next newline
			reader := bufio.NewReader(file)
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return err
			}
			currentPos = start + int64(len(line))
			file.Seek(currentPos, io.SeekStart)
		} else {
			// We're exactly at a line boundary - no need to skip
			file.Seek(start, io.SeekStart)
		}
	}

	// Create our own scanner to control line reading
	scanner := bufio.NewScanner(file)

	// Create a CSV reader just for parsing (not scanning)
	// This caches reflection metadata for performance
	parser := NewReaderWithConfig[T](nil, ReaderConfig{
		Delimiter: config.Delimiter,
		HasHeader: false,
	})

	for scanner.Scan() {
		lineBytes := scanner.Bytes()
		lineLen := int64(len(lineBytes)) + 1 // +1 for newline

		// Check if this line started within our chunk
		// If the line starts at or after our end boundary, we're done
		if currentPos >= end {
			break
		}

		// Parse the line using the cached parser
		result := parser.ParseLineBytes(lineBytes)
		if result.IsErr() {
			return result.UnwrapErr()
		}
		if err := fn(result.Unwrap()); err != nil {
			return err
		}

		currentPos += lineLen
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// ProcessFileWithReducer provides zero-contention parallel processing with per-worker
// accumulators. Each worker has its own accumulator (created by init), processes records
// via reduce, and all accumulators are combined via merge at the end.
// This is the most performant option for aggregation workloads like 1BRC.
func ProcessFileWithReducer[T, A any](
	path string,
	config FileConfig,
	init func() A,
	reduce func(acc *A, record T),
	merge func(accumulators []A) A,
) (A, error) {
	var zero A

	if !config.Parallel {
		// Sequential mode: single accumulator
		acc := init()
		err := ProcessFileWithError(path, config, func(record T) error {
			reduce(&acc, record)
			return nil
		})
		if err != nil {
			return zero, err
		}
		return acc, nil
	}

	// Parallel mode: per-worker accumulators
	fileInfo, err := os.Stat(path)
	if err != nil {
		return zero, err
	}

	workers := config.Workers
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	chunks := createChunks(fileInfo.Size(), workers)
	accumulators := make([]A, len(chunks))

	var wg sync.WaitGroup
	errChan := make(chan error, workers)

	for i, c := range chunks {
		wg.Add(1)
		go func(idx int, start, end int64) {
			defer wg.Done()

			// Create worker-local accumulator
			acc := init()

			err := processChunk[T](path, start, end, config, func(record T) error {
				reduce(&acc, record)
				return nil
			})

			if err != nil {
				errChan <- err
				return
			}

			accumulators[idx] = acc
		}(i, c.Start, c.End)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		return zero, err
	}

	// Merge all accumulators
	return merge(accumulators), nil
}
