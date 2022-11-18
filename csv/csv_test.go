package csv_test

import (
	"fmt"
	"github.com/marlaone/shepard/csv"
	"strings"
	"testing"
)

type Item struct {
	Title     string
	Available bool
}

func TestReader_Deserialize(t *testing.T) {
	csvContent := `title,available
"hello",true
"world", false`

	r := csv.NewReaderFromReader[Item](strings.NewReader(csvContent))
	for {
		item := r.Deserialize()
		if item.IsErr() {
			t.Fatalf("%v", item.Unwrap())
		}

		fmt.Println("item", item.Unwrap())
	}
}
