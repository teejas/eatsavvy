package worker

import (
	"eatsavvy/pkg/places"
	"testing"
	"time"
)

func TestIsRestaurantOpen(t *testing.T) {
	tests := []struct {
		name          string
		openHours     []places.TimeRange
		currentDay    int
		currentHour   int
		currentMinute int
		want          bool
	}{
		{
			name: "[open] during regular hours same day",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:  1,
			currentHour: 14,
			want:        true,
		},
		{
			name: "[closed] before opening time",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:  1,
			currentHour: 8,
			want:        false,
		},
		{
			name: "[closed] after closing time",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:  1,
			currentHour: 23,
			want:        false,
		},
		{
			name: "[open] spans to next day - checking during late night",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 5, Hour: 18, Minute: 0},
					Close: places.TimePoint{Weekday: 6, Hour: 2, Minute: 0},
				},
			},
			currentDay:  5,
			currentHour: 23,
			want:        true,
		},
		{
			name: "[open] from previous day - still open next day",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 5, Hour: 18, Minute: 0},
					Close: places.TimePoint{Weekday: 6, Hour: 2, Minute: 0},
				},
			},
			currentDay:  6,
			currentHour: 1,
			want:        true,
		},
		{
			name: "[closed] on a different day",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:  2,
			currentHour: 14,
			want:        false,
		},
		{
			name:        "empty open hours",
			openHours:   []places.TimeRange{},
			currentDay:  1,
			currentHour: 14,
			want:        false,
		},
		{
			name: "[open] multiple time ranges - open in second range",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 14, Minute: 0},
				},
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 17, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:  1,
			currentHour: 19,
			want:        true,
		},
		{
			name: "[closed] multiple time ranges - closed between ranges",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 14, Minute: 0},
				},
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 17, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:  1,
			currentHour: 15,
			want:        false,
		},
		{
			name: "[open] exactly at opening time - restaurant just opened",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   9,
			currentMinute: 0,
			want:          true,
		},
		{
			name: "[closed] exactly at closing hour - should be closed",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:  1,
			currentHour: 22,
			want:        false,
		},
		{
			name: "[open] opens at midnight (hour 0)",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 0, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 6, Minute: 0},
				},
			},
			currentDay:  1,
			currentHour: 3,
			want:        true,
		},
		{
			name: "[open] late night bar - checking at 1 AM same day it opened",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 5, Hour: 18, Minute: 0},
					Close: places.TimePoint{Weekday: 6, Hour: 2, Minute: 0},
				},
			},
			currentDay:  5,
			currentHour: 23,
			want:        true,
		},
		{
			name: "[open] late night bar - checking next day before close",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 5, Hour: 18, Minute: 0},
					Close: places.TimePoint{Weekday: 6, Hour: 2, Minute: 0},
				},
			},
			currentDay:  6,
			currentHour: 1,
			want:        true,
		},
		{
			name: "[closed] late night bar - checking next day after close",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 5, Hour: 18, Minute: 0},
					Close: places.TimePoint{Weekday: 6, Hour: 2, Minute: 0},
				},
			},
			currentDay:  6,
			currentHour: 3,
			want:        false,
		},
		{
			name: "[closed] week wrap - Saturday night bar, checking Sunday after close",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 6, Hour: 20, Minute: 0},
					Close: places.TimePoint{Weekday: 0, Hour: 2, Minute: 0},
				},
			},
			currentDay:  0,
			currentHour: 3,
			want:        false,
		},
		{
			name: "[open] week wrap - Saturday night bar, checking Sunday before close",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 6, Hour: 20, Minute: 0},
					Close: places.TimePoint{Weekday: 0, Hour: 2, Minute: 0},
				},
			},
			currentDay:  0,
			currentHour: 1,
			want:        true,
		},
		// BUG TEST: Week wrap on the opening day itself
		{
			name: "[open] week wrap - Saturday night bar, checking SATURDAY before midnight",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 6, Hour: 20, Minute: 0},
					Close: places.TimePoint{Weekday: 0, Hour: 2, Minute: 0},
				},
			},
			currentDay:  6,
			currentHour: 23,
			want:        true,
		},
		// BUG TEST: Minute precision - opened 30 min ago
		{
			name: "[open] minute precision - opened 30 min ago same hour",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   9,
			currentMinute: 30,
			want:          true,
		},
		// BUG TEST: Minute precision - closes in 30 min
		{
			name: "[open] minute precision - closes in 30 min same hour",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 30},
				},
			},
			currentDay:    1,
			currentHour:   22,
			currentMinute: 15,
			want:          true,
		},
		// BUG TEST: Minute precision - closed 15 min ago
		{
			name: "[closed] minute precision - closed 15 min ago same hour",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   22,
			currentMinute: 15,
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRestaurantOpen(tt.openHours, tt.currentDay, tt.currentHour, tt.currentMinute)
			if got != tt.want {
				t.Errorf("isRestaurantOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCallbackTime(t *testing.T) {
	tests := []struct {
		name          string
		openHours     []places.TimeRange
		currentDay    int
		currentHour   int
		currentMinute int
		want          time.Duration
	}{
		{
			name: "opens later today",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 14, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   10,
			currentMinute: 0,
			want:          4*time.Hour + 30*time.Minute, // 4 hours + 30 min buffer
		},
		{
			name: "opens later today with minute offset",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 14, Minute: 30},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   10,
			currentMinute: 15,
			want:          4*time.Hour + 15*time.Minute + 30*time.Minute, // 4h + (30-15)m + 30m buffer
		},
		{
			name: "opens tomorrow",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 2, Hour: 9, Minute: 0},
					Close: places.TimePoint{Weekday: 2, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   20,
			currentMinute: 0,
			want:          (24-20+9)*time.Hour + 30*time.Minute, // 13 hours + 30 min buffer
		},
		{
			name: "opens tomorrow with minute offset",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 2, Hour: 9, Minute: 30},
					Close: places.TimePoint{Weekday: 2, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   20,
			currentMinute: 15,
			want:          (24-20+9)*time.Hour + 15*time.Minute + 30*time.Minute, // 13h + (30-15)m + 30m buffer
		},
		{
			name:          "empty open hours returns default 1 hour",
			openHours:     []places.TimeRange{},
			currentDay:    1,
			currentHour:   10,
			currentMinute: 0,
			want:          1 * time.Hour,
		},
		{
			name: "multiple ranges - picks first applicable",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 8, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 12, Minute: 0},
				},
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 17, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   14,
			currentMinute: 0,
			want:          3*time.Hour + 30*time.Minute, // 17-14 = 3 hours + 30 min buffer
		},
		{
			name: "wrap to next week - Tuesday to next Monday",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 8, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 12, Minute: 0},
				},
			},
			currentDay:    2,
			currentHour:   10,
			currentMinute: 0,
			// Days until next Monday: 7 - 2 + 1 = 6 days
			// Hours: 6*24 - 10 + 8 = 142 hours + 30 min buffer
			want: (6*24-10+8)*time.Hour + 30*time.Minute,
		},
		{
			name: "negative minute difference",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 14, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   10,
			currentMinute: 45,
			want:          4*time.Hour + (-45)*time.Minute + 30*time.Minute, // 4h - 45m + 30m = 3h45m
		},
		{
			name: "late night restaurant - next opening is tomorrow",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 3, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 3, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 4, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 4, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 5, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 5, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 6, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 6, Hour: 5, Minute: 30},
				},
			},
			currentDay:    5,
			currentHour:   10,
			currentMinute: 33,
			want:          (24-10+1)*time.Hour + (0-33)*time.Minute + 30*time.Minute, // 15h - 33m + 30m = 14h57m
		},
		{
			name: "same day but already passed - wait until next week",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 6, Hour: 1, Minute: 30},
					Close: places.TimePoint{Weekday: 6, Hour: 4, Minute: 30},
				},
			},
			currentDay:    6,  // Saturday
			currentHour:   10, // 10:45 UTC (after 4:30 close)
			currentMinute: 45,
			// Next opening is next Saturday at 1:30 AM (7 days later)
			// Hours: 7*24 - 10 + 1 = 159 hours + (30-45) minutes + 30 min buffer
			want: (7*24-10+1)*time.Hour + (30-45)*time.Minute + 30*time.Minute, // 159h + 15m = 159h15m
		},
		{
			name: "week wrap-around - Saturday to Sunday",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 0, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 0, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 1, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 2, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 2, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 3, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 3, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 4, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 4, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 5, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 5, Hour: 5, Minute: 30},
				},
				{
					Open:  places.TimePoint{Weekday: 6, Hour: 1, Minute: 0},
					Close: places.TimePoint{Weekday: 6, Hour: 5, Minute: 30},
				},
			},
			currentDay:    6,  // Saturday
			currentHour:   10, // 10:30 UTC (after 5:30 close)
			currentMinute: 30,
			// Next opening is Sunday (day 0) at 1:00 AM
			// Minutes until open: (24-10)*60 - 30 + 1*60 = 14*60 - 30 + 60 = 870 minutes = 14h30m + 30m buffer
			want: 14*time.Hour + 30*time.Minute + 30*time.Minute, // 15h
		},
		// BUG TEST: Same hour but opens later (minute precision)
		{
			name: "same hour opens later - minute precision",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 14, Minute: 30},
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   14,
			currentMinute: 15,
			// Opens in 15 minutes + 30 min buffer = 45 min
			want: 45 * time.Minute,
		},
		// BUG TEST: Multiple ranges - should find soonest (dinner before lunch in array)
		{
			name: "multiple ranges - finds soonest opening",
			openHours: []places.TimeRange{
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 11, Minute: 0}, // Lunch - already passed
					Close: places.TimePoint{Weekday: 1, Hour: 14, Minute: 0},
				},
				{
					Open:  places.TimePoint{Weekday: 1, Hour: 17, Minute: 0}, // Dinner - next opening
					Close: places.TimePoint{Weekday: 1, Hour: 22, Minute: 0},
				},
			},
			currentDay:    1,
			currentHour:   15,
			currentMinute: 0,
			// Dinner opens in 2 hours + 30 min buffer
			want: 2*time.Hour + 30*time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getCallbackTime(tt.openHours, tt.currentDay, tt.currentHour, tt.currentMinute)
			if got != tt.want {
				t.Errorf("getCallbackTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
