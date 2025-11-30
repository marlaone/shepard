package main

import (
	"testing"
)

func BenchmarkHotelSelection(b *testing.B) {
	result := parseHotels("hotels5.txt")
	if result.IsErr() {
		b.Fatal(result.UnwrapErr())
	}
	routeData := result.Unwrap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findOptimalSelection(
			routeData.Hotels,
			MaxDrivingMinutesPerDay,
			MaxDays,
			routeData.TotalDistance,
		)
	}
}

func BenchmarkSelectHotels(b *testing.B) {
	result := parseHotels("hotels5.txt")
	if result.IsErr() {
		b.Fatal(result.UnwrapErr())
	}
	routeData := result.Unwrap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		selectHotels(
			routeData.Hotels,
			4.0,
			MaxDrivingMinutesPerDay,
			MaxDays,
			routeData.TotalDistance,
		)
	}
}

func BenchmarkGetAllRatings(b *testing.B) {
	result := parseHotels("hotels5.txt")
	if result.IsErr() {
		b.Fatal(result.UnwrapErr())
	}
	routeData := result.Unwrap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getAllRatings(routeData.Hotels)
	}
}
