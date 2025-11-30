package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/marlaone/shepard/collections/hashmap"
)

// Station names for test data generation
var stationNames = []string{
	"Hamburg", "Berlin", "Munich", "Frankfurt", "Cologne", "Stuttgart", "Düsseldorf",
	"Leipzig", "Dortmund", "Essen", "Bremen", "Dresden", "Hanover", "Nuremberg",
	"Duisburg", "Bochum", "Wuppertal", "Bielefeld", "Bonn", "Münster",
	"Paris", "Lyon", "Marseille", "Toulouse", "Nice", "Nantes", "Strasbourg",
	"London", "Birmingham", "Manchester", "Leeds", "Glasgow", "Liverpool", "Bristol",
	"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia",
	"Tokyo", "Osaka", "Kyoto", "Yokohama", "Nagoya", "Sapporo", "Fukuoka",
	"Sydney", "Melbourne", "Brisbane", "Perth", "Adelaide", "Canberra",
	"Moscow", "Saint Petersburg", "Novosibirsk", "Yekaterinburg", "Kazan",
	"Beijing", "Shanghai", "Guangzhou", "Shenzhen", "Chengdu", "Hangzhou",
	"São Paulo", "Rio de Janeiro", "Brasília", "Salvador", "Fortaleza",
	"Mumbai", "Delhi", "Bangalore", "Hyderabad", "Chennai", "Kolkata",
	"Cairo", "Lagos", "Johannesburg", "Nairobi", "Casablanca", "Accra",
	"Dubai", "Singapore", "Hong Kong", "Bangkok", "Jakarta", "Manila",
	"Toronto", "Vancouver", "Montreal", "Calgary", "Ottawa", "Edmonton",
	"Mexico City", "Guadalajara", "Monterrey", "Puebla", "Tijuana",
	"Buenos Aires", "Santiago", "Lima", "Bogotá", "Caracas", "Quito",
}

// generateTestData creates a test file with the specified number of rows
func generateTestData(filename string, numRows int) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	rng := rand.New(rand.NewSource(42)) // Fixed seed for reproducibility

	for i := 0; i < numRows; i++ {
		station := stationNames[rng.Intn(len(stationNames))]
		// Temperature range: -99.9 to 99.9
		temp := (rng.Float64() * 199.8) - 99.9
		fmt.Fprintf(w, "%s;%.1f\n", station, temp)
	}

	return nil
}

// TestGenerateTestData generates test data files of various sizes
func TestGenerateTestData(t *testing.T) {
	sizes := map[string]int{
		"testdata/10k.csv":  10_000,
		"testdata/100k.csv": 100_000,
		"testdata/1m.csv":   1_000_000,
		"testdata/10m.csv":  10_000_000,
	}

	for filename, size := range sizes {
		t.Run(filename, func(t *testing.T) {
			if err := generateTestData(filename, size); err != nil {
				t.Fatalf("Failed to generate %s: %v", filename, err)
			}
			t.Logf("Generated %s with %d rows", filename, size)
		})
	}
}

// ensureTestData creates test data if it doesn't exist
func ensureTestData(t testing.TB, filename string, numRows int) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := generateTestData(filename, numRows); err != nil {
			t.Fatalf("Failed to generate test data: %v", err)
		}
	}
}

// ============================================================================
// Full Pipeline Benchmarks
// ============================================================================

func BenchmarkPipeline_10K(b *testing.B) {
	ensureTestData(b, "testdata/10k.csv", 10_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/10k.csv", runtime.NumCPU())
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
	}
}

func BenchmarkPipeline_100K(b *testing.B) {
	ensureTestData(b, "testdata/100k.csv", 100_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/100k.csv", runtime.NumCPU())
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
	}
}

func BenchmarkPipeline_1M(b *testing.B) {
	ensureTestData(b, "testdata/1m.csv", 1_000_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/1m.csv", runtime.NumCPU())
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
	}
}

func BenchmarkPipeline_10M(b *testing.B) {
	ensureTestData(b, "testdata/10m.csv", 10_000_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/10m.csv", runtime.NumCPU())
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
	}
}

// ============================================================================
// Single Worker vs Parallel Comparison
// ============================================================================

func BenchmarkPipeline_1M_SingleWorker(b *testing.B) {
	ensureTestData(b, "testdata/1m.csv", 1_000_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/1m.csv", 1)
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
	}
}

func BenchmarkPipeline_1M_2Workers(b *testing.B) {
	ensureTestData(b, "testdata/1m.csv", 1_000_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/1m.csv", 2)
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
	}
}

func BenchmarkPipeline_1M_4Workers(b *testing.B) {
	ensureTestData(b, "testdata/1m.csv", 1_000_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/1m.csv", 4)
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
	}
}

func BenchmarkPipeline_1M_8Workers(b *testing.B) {
	ensureTestData(b, "testdata/1m.csv", 1_000_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/1m.csv", 8)
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
	}
}

// ============================================================================
// Component Benchmarks: HashMap Entry API
// ============================================================================

func BenchmarkHashMap_Entry_10K(b *testing.B) {
	stations := make([]string, 10000)
	rng := rand.New(rand.NewSource(42))
	for i := range stations {
		stations[i] = stationNames[rng.Intn(len(stationNames))]
	}
	temps := make([]float64, 10000)
	for i := range temps {
		temps[i] = (rng.Float64() * 199.8) - 99.9
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		localMap := hashmap.WithCapacity[string, *StationStats](100)

		for j := 0; j < len(stations); j++ {
			station := stations[j]
			temp := temps[j]

			entry := localMap.Entry(station)
			entry.AndModify(func(s **StationStats) {
				stats := *s
				if temp < stats.Min {
					stats.Min = temp
				}
				if temp > stats.Max {
					stats.Max = temp
				}
				stats.Sum += temp
				stats.Count++
			}).OrInsertWith(func() *StationStats {
				return &StationStats{
					Min:   temp,
					Max:   temp,
					Sum:   temp,
					Count: 1,
				}
			})
		}
	}
}

// ============================================================================
// Native Go Comparison Benchmarks
// ============================================================================

func BenchmarkNativePipeline_1M(b *testing.B) {
	ensureTestData(b, "testdata/1m.csv", 1_000_000)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := processFileNative("testdata/1m.csv")
		if result != nil {
			b.Fatal(result)
		}
	}
}

// processFileNative is a native Go implementation for comparison
func processFileNative(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	stats := make(map[string]*StationStats, 10000)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse line
		idx := -1
		for i := 0; i < len(line); i++ {
			if line[i] == ';' {
				idx = i
				break
			}
		}
		if idx == -1 {
			continue
		}

		station := line[:idx]
		var temp float64
		fmt.Sscanf(line[idx+1:], "%f", &temp)

		if s, ok := stats[station]; ok {
			if temp < s.Min {
				s.Min = temp
			}
			if temp > s.Max {
				s.Max = temp
			}
			s.Sum += temp
			s.Count++
		} else {
			stats[station] = &StationStats{
				Min:   temp,
				Max:   temp,
				Sum:   temp,
				Count: 1,
			}
		}
	}

	return scanner.Err()
}

func BenchmarkNativeHashMap_Entry_10K(b *testing.B) {
	stations := make([]string, 10000)
	rng := rand.New(rand.NewSource(42))
	for i := range stations {
		stations[i] = stationNames[rng.Intn(len(stationNames))]
	}
	temps := make([]float64, 10000)
	for i := range temps {
		temps[i] = (rng.Float64() * 199.8) - 99.9
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		stats := make(map[string]*StationStats, 100)

		for j := 0; j < len(stations); j++ {
			station := stations[j]
			temp := temps[j]

			if s, ok := stats[station]; ok {
				if temp < s.Min {
					s.Min = temp
				}
				if temp > s.Max {
					s.Max = temp
				}
				s.Sum += temp
				s.Count++
			} else {
				stats[station] = &StationStats{
					Min:   temp,
					Max:   temp,
					Sum:   temp,
					Count: 1,
				}
			}
		}
	}
}

// ============================================================================
// Memory Profile Benchmarks
// ============================================================================

func BenchmarkProfile_Memory_1M(b *testing.B) {
	ensureTestData(b, "testdata/1m.csv", 1_000_000)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := processFile("testdata/1m.csv", runtime.NumCPU())
		if result.IsErr() {
			b.Fatal(result.UnwrapErr())
		}
		// Format output to test full pipeline
		_ = formatOutput(result.Unwrap())
	}
}

// ============================================================================
// Line Parsing Benchmark
// ============================================================================

func BenchmarkProcessLine_10K(b *testing.B) {
	lines := make([][]byte, 10000)
	rng := rand.New(rand.NewSource(42))
	for i := range lines {
		station := stationNames[rng.Intn(len(stationNames))]
		temp := (rng.Float64() * 199.8) - 99.9
		lines[i] = []byte(fmt.Sprintf("%s;%.1f", station, temp))
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		localMap := hashmap.WithCapacity[string, *StationStats](100)
		for _, line := range lines {
			processLineBytes(line, &localMap)
		}
	}
}
