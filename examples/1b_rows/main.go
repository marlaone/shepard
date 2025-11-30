package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/marlaone/shepard/collections/hashmap"
	"github.com/marlaone/shepard/csv"
)

// Measurement represents a single CSV row: station;temperature
type Measurement struct {
	Station     string
	Temperature float64
}

// StationStats aggregates statistics for a weather station
type StationStats struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int64
}

func main() {
	fileFlag := flag.String("file", "measurements.txt", "Input file path")
	cpuprofile := flag.String("cpuprofile", "", "Write CPU profile to file")
	memprofile := flag.String("memprofile", "", "Write memory profile to file")
	workers := flag.Int("workers", runtime.NumCPU(), "Number of worker goroutines")
	flag.Parse()

	// CPU profiling
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Error starting CPU profile: %v\n", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	// Process file using csv.ProcessFileWithReducer
	result, err := csv.ProcessFileWithReducer(
		*fileFlag,
		csv.FileConfig{
			Delimiter: ';',
			HasHeader: false,
			Parallel:  true,
			Workers:   *workers,
		},
		// init: create per-worker hashmap
		func() hashmap.HashMap[string, *StationStats] {
			return hashmap.WithCapacity[string, *StationStats](10000)
		},
		// reduce: process each record into worker-local map
		func(localMap *hashmap.HashMap[string, *StationStats], m Measurement) {
			entry := localMap.Entry(m.Station)
			entry.AndModify(func(s **StationStats) {
				stats := *s
				if m.Temperature < stats.Min {
					stats.Min = m.Temperature
				}
				if m.Temperature > stats.Max {
					stats.Max = m.Temperature
				}
				stats.Sum += m.Temperature
				stats.Count++
			}).OrInsertWith(func() *StationStats {
				return &StationStats{
					Min:   m.Temperature,
					Max:   m.Temperature,
					Sum:   m.Temperature,
					Count: 1,
				}
			})
		},
		// merge: combine all worker maps into final result
		mergeResults,
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Output result
	fmt.Println(formatOutput(result))

	// Memory profiling
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating memory profile: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing memory profile: %v\n", err)
			os.Exit(1)
		}
	}
}

// mergeResults merges all worker HashMaps into a single result
func mergeResults(results []hashmap.HashMap[string, *StationStats]) hashmap.HashMap[string, *StationStats] {
	final := hashmap.WithCapacity[string, *StationStats](10000)

	for _, workerMap := range results {
		// Use IterUnordered for faster iteration (no sorting needed during merge)
		workerMap.IterUnordered().Foreach(func(_ int, pair hashmap.Pair[string, *StationStats]) {
			entry := final.Entry(pair.Key)
			entry.AndModify(func(existing **StationStats) {
				stats := *existing
				workerStats := pair.Value
				if workerStats.Min < stats.Min {
					stats.Min = workerStats.Min
				}
				if workerStats.Max > stats.Max {
					stats.Max = workerStats.Max
				}
				stats.Sum += workerStats.Sum
				stats.Count += workerStats.Count
			}).OrInsertWith(func() *StationStats {
				// Copy the stats to avoid sharing pointers
				return &StationStats{
					Min:   pair.Value.Min,
					Max:   pair.Value.Max,
					Sum:   pair.Value.Sum,
					Count: pair.Value.Count,
				}
			})
		})
	}

	return final
}

// formatOutput generates the sorted output string
func formatOutput(stats hashmap.HashMap[string, *StationStats]) string {
	var sb strings.Builder
	sb.WriteString("{")

	first := true
	// Keys() returns sorted keys automatically
	stats.Keys().Foreach(func(_ int, station string) {
		opt := stats.Get(station)
		if opt.IsNone() {
			return
		}
		// Get() returns Option[**StationStats] since our value type is *StationStats
		s := *opt.Unwrap()
		mean := s.Sum / float64(s.Count)

		if !first {
			sb.WriteString(", ")
		}
		first = false

		fmt.Fprintf(&sb, "%s=%.1f/%.1f/%.1f", station, s.Min, mean, s.Max)
	})

	sb.WriteString("}")
	return sb.String()
}
