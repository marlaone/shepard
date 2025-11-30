package csv_test

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/marlaone/shepard/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestRecord struct {
	Name  string
	Value float64
}

func createTestFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.csv")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)
	return tmpFile
}

func TestProcessFile_Sequential(t *testing.T) {
	content := "Alice;10.5\nBob;20.3\nCharlie;15.7\n"
	tmpFile := createTestFile(t, content)

	var records []TestRecord
	err := csv.ProcessFile(tmpFile, csv.FileConfig{
		Delimiter: ';',
		HasHeader: false,
	}, func(r TestRecord) {
		records = append(records, r)
	})

	assert.NoError(t, err)
	assert.Len(t, records, 3)
	assert.Equal(t, "Alice", records[0].Name)
	assert.InDelta(t, 10.5, records[0].Value, 0.001)
	assert.Equal(t, "Bob", records[1].Name)
	assert.InDelta(t, 20.3, records[1].Value, 0.001)
	assert.Equal(t, "Charlie", records[2].Name)
	assert.InDelta(t, 15.7, records[2].Value, 0.001)
}

func TestProcessFile_WithHeader(t *testing.T) {
	content := "name,value\nAlice,10.5\nBob,20.3\n"
	tmpFile := createTestFile(t, content)

	var records []TestRecord
	err := csv.ProcessFile(tmpFile, csv.FileConfig{
		Delimiter: ',',
		HasHeader: true,
	}, func(r TestRecord) {
		records = append(records, r)
	})

	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, "Alice", records[0].Name)
	assert.Equal(t, "Bob", records[1].Name)
}

func TestProcessFileWithError_EarlyTermination(t *testing.T) {
	content := "Alice;10.5\nBob;20.3\nCharlie;15.7\n"
	tmpFile := createTestFile(t, content)

	var count int
	err := csv.ProcessFileWithError(tmpFile, csv.FileConfig{
		Delimiter: ';',
		HasHeader: false,
	}, func(r TestRecord) error {
		count++
		if count >= 2 {
			return assert.AnError
		}
		return nil
	})

	assert.Error(t, err)
	assert.Equal(t, 2, count)
}

func TestProcessFile_Parallel(t *testing.T) {
	// Create a larger file to test parallel processing
	var content string
	for i := 0; i < 1000; i++ {
		content += "Station;25.5\n"
	}
	tmpFile := createTestFile(t, content)

	var count atomic.Int64
	err := csv.ProcessFile(tmpFile, csv.FileConfig{
		Delimiter: ';',
		HasHeader: false,
		Parallel:  true,
		Workers:   4,
	}, func(r TestRecord) {
		count.Add(1)
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1000), count.Load())
}

func TestProcessFile_FileNotFound(t *testing.T) {
	err := csv.ProcessFile("/nonexistent/file.csv", csv.FileConfig{
		Delimiter: ',',
	}, func(r TestRecord) {})

	assert.Error(t, err)
}

func TestProcessFileWithReducer_Sequential(t *testing.T) {
	content := "Alice;10.0\nBob;20.0\nAlice;30.0\n"
	tmpFile := createTestFile(t, content)

	type Stats struct {
		Sum   float64
		Count int
	}

	result, err := csv.ProcessFileWithReducer(
		tmpFile,
		csv.FileConfig{
			Delimiter: ';',
			HasHeader: false,
			Parallel:  false,
		},
		// init
		func() map[string]*Stats {
			return make(map[string]*Stats)
		},
		// reduce
		func(acc *map[string]*Stats, r TestRecord) {
			m := *acc
			if m[r.Name] == nil {
				m[r.Name] = &Stats{}
			}
			m[r.Name].Sum += r.Value
			m[r.Name].Count++
		},
		// merge
		func(accumulators []map[string]*Stats) map[string]*Stats {
			final := make(map[string]*Stats)
			for _, acc := range accumulators {
				for k, v := range acc {
					if final[k] == nil {
						final[k] = &Stats{}
					}
					final[k].Sum += v.Sum
					final[k].Count += v.Count
				}
			}
			return final
		},
	)

	assert.NoError(t, err)
	assert.InDelta(t, 40.0, result["Alice"].Sum, 0.001)
	assert.Equal(t, 2, result["Alice"].Count)
	assert.InDelta(t, 20.0, result["Bob"].Sum, 0.001)
	assert.Equal(t, 1, result["Bob"].Count)
}

func TestProcessFileWithReducer_Parallel(t *testing.T) {
	// Create larger file for parallel testing
	var content string
	for i := 0; i < 1000; i++ {
		content += "Station;10.0\n"
	}
	tmpFile := createTestFile(t, content)

	type Stats struct {
		Sum   float64
		Count int
	}

	result, err := csv.ProcessFileWithReducer(
		tmpFile,
		csv.FileConfig{
			Delimiter: ';',
			HasHeader: false,
			Parallel:  true,
			Workers:   4,
		},
		// init
		func() map[string]*Stats {
			return make(map[string]*Stats)
		},
		// reduce
		func(acc *map[string]*Stats, r TestRecord) {
			m := *acc
			if m[r.Name] == nil {
				m[r.Name] = &Stats{}
			}
			m[r.Name].Sum += r.Value
			m[r.Name].Count++
		},
		// merge
		func(accumulators []map[string]*Stats) map[string]*Stats {
			final := make(map[string]*Stats)
			for _, acc := range accumulators {
				for k, v := range acc {
					if final[k] == nil {
						final[k] = &Stats{}
					}
					final[k].Sum += v.Sum
					final[k].Count += v.Count
				}
			}
			return final
		},
	)

	assert.NoError(t, err)
	assert.InDelta(t, 10000.0, result["Station"].Sum, 0.001)
	assert.Equal(t, 1000, result["Station"].Count)
}

func TestProcessFile_EmptyFile(t *testing.T) {
	tmpFile := createTestFile(t, "")

	var count int
	err := csv.ProcessFile(tmpFile, csv.FileConfig{
		Delimiter: ';',
		HasHeader: false,
	}, func(r TestRecord) {
		count++
	})

	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestProcessFile_DefaultWorkers(t *testing.T) {
	content := "Alice;10.5\nBob;20.3\n"
	tmpFile := createTestFile(t, content)

	var count atomic.Int64
	err := csv.ProcessFile(tmpFile, csv.FileConfig{
		Delimiter: ';',
		HasHeader: false,
		Parallel:  true,
		// Workers defaults to runtime.NumCPU()
	}, func(r TestRecord) {
		count.Add(1)
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count.Load())
}
