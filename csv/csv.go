package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/marlaone/shepard"
	"io"
	"reflect"
	"strconv"
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

	for n, t := range r.recordDescriber {
		value := records[r.header[n]]
		fmt.Println(n, value)
		switch t.Kind() {
		case reflect.String:
			recordValue.FieldByName(n).SetString(value)
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return shepard.Err[T, error](err)
			}
			recordValue.FieldByName(n).SetBool(v)
		}
	}

	return shepard.Err[T, error](errors.New("unimplemented"))
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
