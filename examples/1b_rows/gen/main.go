package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Station names for test data generation (100 unique stations)
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

func main() {
	numRows := flag.Int64("rows", 1_000_000_000, "Number of rows to generate")
	output := flag.String("output", "measurements.txt", "Output file path")
	flag.Parse()

	fmt.Printf("Generating %d rows to %s...\n", *numRows, *output)
	start := time.Now()

	if err := generateTestData(*output, *numRows); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	fmt.Printf("Done! Generated %d rows in %v\n", *numRows, elapsed)
	fmt.Printf("Speed: %.2f million rows/second\n", float64(*numRows)/elapsed.Seconds()/1_000_000)
}

func generateTestData(filename string, numRows int64) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Use large buffer for efficient writes (64KB)
	w := bufio.NewWriterSize(f, 64*1024)
	defer w.Flush()

	rng := rand.New(rand.NewSource(42)) // Fixed seed for reproducibility
	numStations := len(stationNames)

	// Pre-format temperature values to avoid fmt.Fprintf overhead
	// Temperature range: -99.9 to 99.9
	var buf [64]byte

	progressInterval := int64(10_000_000) // Report every 10M rows
	lastReport := time.Now()

	for i := int64(0); i < numRows; i++ {
		station := stationNames[rng.Intn(numStations)]
		// Generate temperature: -99.9 to 99.9
		temp := (rng.Float64() * 199.8) - 99.9

		// Write station name
		w.WriteString(station)
		w.WriteByte(';')

		// Format temperature with 1 decimal place manually (faster than fmt)
		n := formatFloat(buf[:], temp)
		w.Write(buf[:n])
		w.WriteByte('\n')

		// Progress reporting
		if (i+1)%progressInterval == 0 {
			now := time.Now()
			elapsed := now.Sub(lastReport)
			rate := float64(progressInterval) / elapsed.Seconds() / 1_000_000
			percent := float64(i+1) / float64(numRows) * 100
			fmt.Printf("Progress: %.1f%% (%d/%d rows) - %.2f M rows/sec\n",
				percent, i+1, numRows, rate)
			lastReport = now
		}
	}

	return w.Flush()
}

// formatFloat formats a float64 with 1 decimal place into buf, returns length
func formatFloat(buf []byte, f float64) int {
	neg := f < 0
	if neg {
		f = -f
	}

	// Round to 1 decimal place
	f = float64(int64(f*10+0.5)) / 10

	intPart := int64(f)
	decPart := int64((f-float64(intPart))*10 + 0.5)
	if decPart >= 10 {
		decPart = 0
		intPart++
	}

	// Build string backwards
	i := len(buf)

	// Decimal digit
	i--
	buf[i] = byte('0' + decPart)

	// Decimal point
	i--
	buf[i] = '.'

	// Integer part
	if intPart == 0 {
		i--
		buf[i] = '0'
	} else {
		for intPart > 0 {
			i--
			buf[i] = byte('0' + intPart%10)
			intPart /= 10
		}
	}

	// Negative sign
	if neg {
		i--
		buf[i] = '-'
	}

	// Move to beginning
	n := len(buf) - i
	copy(buf[:n], buf[i:])
	return n
}
