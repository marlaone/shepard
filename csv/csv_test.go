package csv_test

import (
	"io"
	"strings"
	"testing"

	"github.com/marlaone/shepard/csv"
	"github.com/stretchr/testify/assert"
)

type Item struct {
	Title     string
	Available bool
}

func TestReader_Deserialize(t *testing.T) {
	csvContent := `title,available
hello,true
world,false`

	r := csv.NewReaderFromReader[Item](strings.NewReader(csvContent))

	// Read first item
	item1 := r.Deserialize()
	assert.True(t, item1.IsOk(), "First item should deserialize successfully")
	assert.Equal(t, "hello", item1.Unwrap().Title)
	assert.Equal(t, true, item1.Unwrap().Available)

	// Read second item
	item2 := r.Deserialize()
	assert.True(t, item2.IsOk(), "Second item should deserialize successfully")
	assert.Equal(t, "world", item2.Unwrap().Title)
	assert.Equal(t, false, item2.Unwrap().Available)

	// Read past end - should get EOF error
	item3 := r.Deserialize()
	assert.True(t, item3.IsErr(), "Third read should return error (EOF)")
}

func TestReader_Deserialize_EOF(t *testing.T) {
	csvContent := `title,available
hello,true`

	r := csv.NewReaderFromReader[Item](strings.NewReader(csvContent))

	// Read first item
	item1 := r.Deserialize()
	assert.True(t, item1.IsOk())

	// Read past end
	item2 := r.Deserialize()
	assert.True(t, item2.IsErr())
	// The error should be io.EOF wrapped
	assert.Contains(t, item2.Err().Unwrap().Error(), io.EOF.Error())
}

type NumericItem struct {
	Name  string
	Count int
	Price float64
}

func TestReader_Deserialize_NumericTypes(t *testing.T) {
	csvContent := `name,count,price
apple,5,1.99
banana,10,0.50`

	r := csv.NewReaderFromReader[NumericItem](strings.NewReader(csvContent))

	// Read first item
	item1 := r.Deserialize()
	assert.True(t, item1.IsOk())
	assert.Equal(t, "apple", item1.Unwrap().Name)
	assert.Equal(t, 5, item1.Unwrap().Count)
	assert.InDelta(t, 1.99, item1.Unwrap().Price, 0.001)

	// Read second item
	item2 := r.Deserialize()
	assert.True(t, item2.IsOk())
	assert.Equal(t, "banana", item2.Unwrap().Name)
	assert.Equal(t, 10, item2.Unwrap().Count)
	assert.InDelta(t, 0.50, item2.Unwrap().Price, 0.001)
}

func TestReader_Deserialize_MissingColumn(t *testing.T) {
	csvContent := `title
hello`

	r := csv.NewReaderFromReader[Item](strings.NewReader(csvContent))

	// Should fail because "available" column is missing
	item := r.Deserialize()
	assert.True(t, item.IsErr())
}

func TestReader_Deserialize_InvalidBool(t *testing.T) {
	csvContent := `title,available
hello,notabool`

	r := csv.NewReaderFromReader[Item](strings.NewReader(csvContent))

	// Should fail because "notabool" is not a valid boolean
	item := r.Deserialize()
	assert.True(t, item.IsErr())
}
