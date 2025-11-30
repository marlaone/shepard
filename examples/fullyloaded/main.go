package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/marlaone/shepard"
	"github.com/marlaone/shepard/collections/slice"
	"github.com/marlaone/shepard/iter"
)

const (
	MaxDrivingMinutesPerDay = 360 // 6 hours
	MaxDays                 = 5
)

type Hotel struct {
	Distance int
	Rating   float64
}

type RouteData struct {
	TotalDistance int
	Hotels        slice.Slice[Hotel]
}

func parseHotels(filename string) shepard.Result[RouteData, error] {
	file, err := os.Open(filename)
	if err != nil {
		return shepard.Err[RouteData, error](err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Line 1: Number of hotels
	if !scanner.Scan() {
		return shepard.Err[RouteData, error](fmt.Errorf("failed to read hotel count"))
	}
	hotelCount, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return shepard.Err[RouteData, error](fmt.Errorf("invalid hotel count: %w", err))
	}

	// Line 2: Total travel time in minutes
	if !scanner.Scan() {
		return shepard.Err[RouteData, error](fmt.Errorf("failed to read total distance"))
	}
	totalDistance, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return shepard.Err[RouteData, error](fmt.Errorf("invalid total distance: %w", err))
	}

	hotels := slice.WithCapacity[Hotel](hotelCount)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return shepard.Err[RouteData, error](fmt.Errorf("invalid hotel line: %s", line))
		}

		distance, err := strconv.Atoi(parts[0])
		if err != nil {
			return shepard.Err[RouteData, error](fmt.Errorf("invalid distance: %w", err))
		}

		rating, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return shepard.Err[RouteData, error](fmt.Errorf("invalid rating: %w", err))
		}

		hotels.Push(Hotel{Distance: distance, Rating: rating})
	}

	if err := scanner.Err(); err != nil {
		return shepard.Err[RouteData, error](err)
	}

	return shepard.Ok[RouteData, error](RouteData{
		TotalDistance: totalDistance,
		Hotels:        hotels,
	})
}

// selectHotels attempts to select hotels with rating >= minRating
// Returns Some(selected hotels) if feasible, None otherwise
func selectHotels(hotels slice.Slice[Hotel], minRating float64, maxDayDistance, maxDays, totalDistance int) shepard.Option[slice.Slice[Hotel]] {
	// FilterLazy + CollectFilter - lazy evaluation, single allocation at collect time
	filteredSlice := slice.CollectFilter(hotels.Iter().FilterLazy(func(h *Hotel) bool {
		return h.Rating >= minRating
	}))

	if filteredSlice.Len() == 0 {
		return shepard.None[slice.Slice[Hotel]]()
	}

	maxNights := maxDays - 1 // 5 days = 4 nights = 4 hotel stays
	selected := slice.WithCapacity[Hotel](maxNights)
	currentPos := 0
	daysUsed := 0

	for currentPos < totalDistance && daysUsed < maxDays {
		maxReachable := currentPos + maxDayDistance

		// Can we reach the destination today?
		if maxReachable >= totalDistance {
			break
		}

		// Find the farthest hotel within reach (greedy)
		bestHotel := shepard.None[Hotel]()
		bestDistance := 0

		filteredSlice.Iter().Foreach(func(_ int, h Hotel) {
			if h.Distance > currentPos && h.Distance <= maxReachable {
				if h.Distance > bestDistance {
					bestDistance = h.Distance
					bestHotel = shepard.Some(h)
				}
			}
		})

		if bestHotel.IsNone() {
			// Can't reach any hotel - infeasible
			return shepard.None[slice.Slice[Hotel]]()
		}

		hotel := bestHotel.Unwrap()
		selected.Push(hotel)
		currentPos = hotel.Distance
		daysUsed++

		if selected.Len() > maxNights {
			// Too many nights needed
			return shepard.None[slice.Slice[Hotel]]()
		}
	}

	// Check if we can reach destination from last position
	if currentPos+maxDayDistance < totalDistance {
		return shepard.None[slice.Slice[Hotel]]()
	}

	return shepard.Some(selected)
}

// getAllRatings extracts all unique ratings from hotels, sorted descending
func getAllRatings(hotels slice.Slice[Hotel]) slice.Slice[float64] {
	// Map hotels to ratings, then deduplicate using Unique
	ratingsIter := iter.MapLazy(hotels.Iter(), func(h Hotel) float64 {
		return h.Rating
	})
	uniqueRatings := slice.CollectUniqueIter(iter.Unique(ratingsIter.Collect()))

	// Sort descending using library method
	slice.SortDesc(&uniqueRatings)

	return uniqueRatings
}

// findOptimalSelection uses binary search on ratings to find the maximum feasible minimum rating
func findOptimalSelection(hotels slice.Slice[Hotel], maxDayDistance, maxDays, totalDistance int) shepard.Option[slice.Slice[Hotel]] {
	ratings := getAllRatings(hotels)

	if ratings.Len() == 0 {
		return shepard.None[slice.Slice[Hotel]]()
	}

	// Binary search: find the highest rating where selection is still feasible
	// Ratings are sorted descending, so we want the earliest (lowest index) that works
	left := 0
	right := ratings.Len() - 1
	bestResult := shepard.None[slice.Slice[Hotel]]()

	for left <= right {
		mid := (left + right) / 2
		rating := ratings.Get(mid).Unwrap()

		result := selectHotels(hotels, *rating, maxDayDistance, maxDays, totalDistance)

		if result.IsSome() {
			// This rating works, try higher (lower index)
			bestResult = result
			right = mid - 1
		} else {
			// This rating doesn't work, try lower (higher index)
			left = mid + 1
		}
	}

	return bestResult
}

// getMinRating returns the minimum rating from selected hotels
func getMinRating(hotels slice.Slice[Hotel]) float64 {
	minRating := 5.0
	hotels.Iter().Foreach(func(_ int, h Hotel) {
		if h.Rating < minRating {
			minRating = h.Rating
		}
	})
	return minRating
}

func main() {
	fileFlag := flag.String("file", "hotels1.txt", "Hotel data file to load (e.g., hotels1.txt, hotels2.txt)")
	flag.Parse()

	result := parseHotels(*fileFlag)

	if result.IsErr() {
		fmt.Printf("Error: %v\n", result.UnwrapErr())
		os.Exit(1)
	}

	routeData := result.Unwrap()

	fmt.Printf("=== Hotel Selection Problem ===\n")
	fmt.Printf("Total travel time: %d minutes\n", routeData.TotalDistance)
	fmt.Printf("Hotels available: %d\n", routeData.Hotels.Len())
	fmt.Printf("Max driving per day: %d minutes (6 hours)\n", MaxDrivingMinutesPerDay)
	fmt.Printf("Max days: %d\n\n", MaxDays)

	selection := findOptimalSelection(
		routeData.Hotels,
		MaxDrivingMinutesPerDay,
		MaxDays,
		routeData.TotalDistance,
	)

	if selection.IsNone() {
		fmt.Println("No feasible solution found!")
		os.Exit(1)
	}

	selectedHotels := selection.Unwrap()
	minRating := getMinRating(selectedHotels)

	fmt.Printf("Optimal minimum rating: %.1f\n", minRating)
	fmt.Printf("Hotels selected: %d\n\n", selectedHotels.Len())

	fmt.Println("Selected hotels:")
	selectedHotels.Iter().Foreach(func(i int, h Hotel) {
		fmt.Printf("  Night %d: Hotel at %d min (rating %.1f)\n", i+1, h.Distance, h.Rating)
	})

	// Verify the solution
	fmt.Println("\nRoute verification:")
	prevPos := 0
	selectedHotels.Iter().Foreach(func(i int, h Hotel) {
		distance := h.Distance - prevPos
		fmt.Printf("  Day %d: Drive %d min (from %d to %d)\n", i+1, distance, prevPos, h.Distance)
		prevPos = h.Distance
	})
	finalDistance := routeData.TotalDistance - prevPos
	fmt.Printf("  Day %d: Drive %d min (from %d to %d - destination)\n",
		selectedHotels.Len()+1, finalDistance, prevPos, routeData.TotalDistance)
}
