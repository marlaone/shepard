package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/marlaone/shepard"
)

type Reader[T any] struct {
	csv    *csv.Reader
	header map[string]int

	recordDescriber map[string]reflect.Type
}

func NewReaderFromReader[T any](r io.Reader) *Reader[T] {
	return &Reader[T]{
		csv:             csv.NewReader(r),
		header:          map[string]int{},
		recordDescriber: map[string]reflect.Type{},
	}
}

func (r *Reader[T]) parseHeader() error {
	records, err := r.csv.Read()
	if err != nil {
		return err
	}
	for i, v := range records {
		r.header[v] = i
	}
	return nil
}

func (r *Reader[T]) parseRecord(records []string) shepard.Result[T, error] {
	var record T

	recordType := reflect.TypeOf(record)
	recordValue := reflect.Indirect(reflect.ValueOf(&record))

	for i := 0; i < recordType.NumField(); i++ {
		field := recordType.Field(i)
		r.recordDescriber[field.Name] = field.Type
	}

	for fieldName, fieldType := range r.recordDescriber {
		// Look up header using lowercase field name to match CSV headers
		headerKey := strings.ToLower(fieldName)
		colIndex, ok := r.header[headerKey]
		if !ok {
			return shepard.Err[T, error](fmt.Errorf("column %q not found in CSV header", headerKey))
		}
		if colIndex >= len(records) {
			return shepard.Err[T, error](fmt.Errorf("column index %d out of range", colIndex))
		}

		value := strings.TrimSpace(records[colIndex])
		switch fieldType.Kind() {
		case reflect.String:
			recordValue.FieldByName(fieldName).SetString(value)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return shepard.Err[T, error](err)
			}
			recordValue.FieldByName(fieldName).SetBool(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return shepard.Err[T, error](err)
			}
			recordValue.FieldByName(fieldName).SetInt(v)
		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return shepard.Err[T, error](err)
			}
			recordValue.FieldByName(fieldName).SetFloat(v)
		}
	}

	return shepard.Ok[T, error](record)
}

func (r *Reader[T]) Deserialize() shepard.Result[T, error] {
	if len(r.header) == 0 {
		if err := r.parseHeader(); err != nil {
			return shepard.Err[T, error](err)
		}
	}
	records, err := r.csv.Read()
	if err != nil {
		return shepard.Err[T, error](err)
	}
	return r.parseRecord(records)
}
