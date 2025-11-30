package csv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/num"
)

// fieldInfo caches reflection metadata for a single struct field
type fieldInfo struct {
	name       string
	index      int
	kind       reflect.Kind
	headerKey  string // lowercase name for header matching
	colIndex   int    // CSV column index (set after header parsing)
	colIndexOk bool   // whether colIndex has been set
}

// ReaderConfig provides configuration options for CSV parsing
type ReaderConfig struct {
	Delimiter byte
	HasHeader bool
}

// Reader is a high-performance CSV reader that deserializes records into structs.
// Uses bufio.Scanner with zero-copy byte parsing for maximum performance.
type Reader[T any] struct {
	scanner   *bufio.Scanner
	delimiter byte
	header    map[string]int

	// Cached reflection metadata - populated once
	recordType  reflect.Type
	fields      []fieldInfo
	initialized bool
	hasHeader   bool
}

// NewReaderFromReader creates a new CSV reader with default comma delimiter and header
func NewReaderFromReader[T any](r io.Reader) *Reader[T] {
	return &Reader[T]{
		scanner:   bufio.NewScanner(r),
		delimiter: ',',
		header:    map[string]int{},
		hasHeader: true,
	}
}

// NewReaderFromReaderWithDelimiter creates a new CSV reader with custom delimiter and header
func NewReaderFromReaderWithDelimiter[T any](r io.Reader, delimiter rune) *Reader[T] {
	return &Reader[T]{
		scanner:   bufio.NewScanner(r),
		delimiter: byte(delimiter),
		header:    map[string]int{},
		hasHeader: true,
	}
}

// NewReaderWithConfig creates a new CSV reader with full configuration options
func NewReaderWithConfig[T any](r io.Reader, config ReaderConfig) *Reader[T] {
	delimiter := config.Delimiter
	if delimiter == 0 {
		delimiter = ','
	}
	reader := &Reader[T]{
		delimiter: delimiter,
		header:    map[string]int{},
		hasHeader: config.HasHeader,
	}
	// Only create scanner if reader is provided (nil reader means parse-only mode)
	if r != nil {
		reader.scanner = bufio.NewScanner(r)
	}
	return reader
}

// initializeMetadata caches all reflection metadata once
func (r *Reader[T]) initializeMetadata() {
	if r.initialized {
		return
	}

	var record T
	r.recordType = reflect.TypeOf(record)
	numFields := r.recordType.NumField()
	r.fields = make([]fieldInfo, numFields)

	for i := 0; i < numFields; i++ {
		field := r.recordType.Field(i)
		r.fields[i] = fieldInfo{
			name:      field.Name,
			index:     i,
			kind:      field.Type.Kind(),
			headerKey: strings.ToLower(field.Name),
		}
	}

	r.initialized = true
}

func (r *Reader[T]) parseHeader() error {
	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return err
		}
		return io.EOF
	}

	line := r.scanner.Bytes()
	colIndex := 0
	start := 0

	for i := 0; i <= len(line); i++ {
		if i == len(line) || line[i] == r.delimiter {
			// Extract column name
			colName := string(bytes.TrimSpace(line[start:i]))
			r.header[colName] = colIndex
			colIndex++
			start = i + 1
		}
	}

	// Map header columns to field indices
	for i := range r.fields {
		if colIndex, ok := r.header[r.fields[i].headerKey]; ok {
			r.fields[i].colIndex = colIndex
			r.fields[i].colIndexOk = true
		}
	}

	return nil
}

// parseRecordBytes parses a line of bytes into struct T
func (r *Reader[T]) parseRecordBytes(line []byte) shepard.Result[T, error] {
	var record T
	recordValue := reflect.Indirect(reflect.ValueOf(&record))

	// For no-header mode, parse fields in a single pass using bytes.IndexByte
	if !r.hasHeader {
		pos := 0
		for i := range r.fields {
			fi := &r.fields[i]

			// Find next delimiter
			end := bytes.IndexByte(line[pos:], r.delimiter)
			var fieldBytes []byte
			if end == -1 {
				// Last field
				fieldBytes = line[pos:]
			} else {
				fieldBytes = line[pos : pos+end]
				pos = pos + end + 1
			}

			fieldValue := recordValue.Field(fi.index)
			switch fi.kind {
			case reflect.String:
				fieldValue.SetString(string(fieldBytes))
			case reflect.Float64:
				result := num.ParseFloat64Bytes(fieldBytes)
				if result.IsErr() {
					return shepard.Err[T, error](result.UnwrapErr())
				}
				fieldValue.SetFloat(result.Unwrap())
			case reflect.Float32:
				result := num.ParseFloat32Bytes(fieldBytes)
				if result.IsErr() {
					return shepard.Err[T, error](result.UnwrapErr())
				}
				fieldValue.SetFloat(float64(result.Unwrap()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result := num.ParseBytes[int64](fieldBytes)
				if result.IsErr() {
					return shepard.Err[T, error](result.UnwrapErr())
				}
				fieldValue.SetInt(result.Unwrap())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				result := num.ParseBytes[uint64](fieldBytes)
				if result.IsErr() {
					return shepard.Err[T, error](result.UnwrapErr())
				}
				fieldValue.SetUint(result.Unwrap())
			case reflect.Bool:
				s := unsafe.String(unsafe.SliceData(fieldBytes), len(fieldBytes))
				v, err := strconv.ParseBool(s)
				if err != nil {
					return shepard.Err[T, error](err)
				}
				fieldValue.SetBool(v)
			}
		}
		return shepard.Ok[T, error](record)
	}

	// For header mode, extract all fields first then map by header indices
	var fields [][]byte
	start := 0
	for i := 0; i <= len(line); i++ {
		if i == len(line) || line[i] == r.delimiter {
			fields = append(fields, line[start:i])
			start = i + 1
		}
	}

	for i := range r.fields {
		fi := &r.fields[i]
		if !fi.colIndexOk {
			return shepard.Err[T, error](fmt.Errorf("column %q not found in CSV header", fi.headerKey))
		}
		if fi.colIndex >= len(fields) {
			return shepard.Err[T, error](fmt.Errorf("column index %d out of range", fi.colIndex))
		}

		fieldBytes := fields[fi.colIndex]
		fieldValue := recordValue.Field(fi.index)

		switch fi.kind {
		case reflect.String:
			fieldValue.SetString(string(fieldBytes))
		case reflect.Float64:
			result := num.ParseFloat64Bytes(fieldBytes)
			if result.IsErr() {
				return shepard.Err[T, error](result.UnwrapErr())
			}
			fieldValue.SetFloat(result.Unwrap())
		case reflect.Float32:
			result := num.ParseFloat32Bytes(fieldBytes)
			if result.IsErr() {
				return shepard.Err[T, error](result.UnwrapErr())
			}
			fieldValue.SetFloat(float64(result.Unwrap()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result := num.ParseBytes[int64](fieldBytes)
			if result.IsErr() {
				return shepard.Err[T, error](result.UnwrapErr())
			}
			fieldValue.SetInt(result.Unwrap())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result := num.ParseBytes[uint64](fieldBytes)
			if result.IsErr() {
				return shepard.Err[T, error](result.UnwrapErr())
			}
			fieldValue.SetUint(result.Unwrap())
		case reflect.Bool:
			s := unsafe.String(unsafe.SliceData(fieldBytes), len(fieldBytes))
			v, err := strconv.ParseBool(s)
			if err != nil {
				return shepard.Err[T, error](err)
			}
			fieldValue.SetBool(v)
		}
	}

	return shepard.Ok[T, error](record)
}

// Deserialize reads and parses the next CSV record into type T
func (r *Reader[T]) Deserialize() shepard.Result[T, error] {
	// Initialize metadata once
	r.initializeMetadata()

	// Parse header on first call (if hasHeader is true)
	if r.hasHeader && len(r.header) == 0 {
		if err := r.parseHeader(); err != nil {
			return shepard.Err[T, error](err)
		}
	}

	// Read next line
	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return shepard.Err[T, error](err)
		}
		return shepard.Err[T, error](io.EOF)
	}

	return r.parseRecordBytes(r.scanner.Bytes())
}

// ParseLineBytes parses a single line of bytes directly into type T.
// This is useful for parsing lines from a scanner without needing the Reader's scanner.
// Note: Call this after creating the reader to ensure metadata is initialized.
func (r *Reader[T]) ParseLineBytes(line []byte) shepard.Result[T, error] {
	r.initializeMetadata()
	return r.parseRecordBytes(line)
}
