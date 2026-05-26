package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Calendar holds the holiday set for a venue (e.g. "USD_NYC", "EUR_TARGET2", "BRL_BRASILIA").
type Calendar struct {
	id       string
	holidays map[string]struct{} // keyed by yyyy-mm-dd UTC
}

// NewCalendar constructs a Calendar from an id + holiday list.
// Dates are normalised to UTC date-only.
func NewCalendar(id string, holidays []time.Time) (*Calendar, error) {
	id = strings.ToUpper(strings.TrimSpace(id))
	if id == "" {
		return nil, fmt.Errorf("%w: calendar id required", ErrInvalidInput)
	}
	c := &Calendar{id: id, holidays: make(map[string]struct{}, len(holidays))}
	for _, h := range holidays {
		c.holidays[keyOf(h)] = struct{}{}
	}
	return c, nil
}

func keyOf(t time.Time) string {
	return t.UTC().Format("2006-01-02")
}

func (c *Calendar) ID() string { return c.id }

// IsHoliday reports whether `day` (date-only, UTC) is a holiday in this calendar.
func (c *Calendar) IsHoliday(day time.Time) bool {
	_, ok := c.holidays[keyOf(day)]
	return ok
}

// IsBusinessDay returns true when `day` is Mon-Fri AND not a holiday.
func (c *Calendar) IsBusinessDay(day time.Time) bool {
	wd := day.UTC().Weekday()
	if wd == time.Saturday || wd == time.Sunday {
		return false
	}
	return !c.IsHoliday(day)
}

// NextBusinessDay walks forward from `day` returning the first business day strictly after it.
func (c *Calendar) NextBusinessDay(day time.Time) time.Time {
	next := day.AddDate(0, 0, 1)
	for !c.IsBusinessDay(next) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

// AddBusinessDays returns the date `n` business days after `day` (n > 0).
// Negative n returns date `|n|` business days before.
func (c *Calendar) AddBusinessDays(day time.Time, n int) time.Time {
	if n == 0 {
		return day
	}
	step := 1
	if n < 0 {
		step = -1
		n = -n
	}
	cur := day
	for i := 0; i < n; i++ {
		cur = cur.AddDate(0, 0, step)
		for !c.IsBusinessDay(cur) {
			cur = cur.AddDate(0, 0, step)
		}
	}
	return cur
}

// HolidaysSorted returns the holiday set in ascending order — useful for export/seed.
func (c *Calendar) HolidaysSorted() []time.Time {
	out := make([]time.Time, 0, len(c.holidays))
	for k := range c.holidays {
		t, _ := time.Parse("2006-01-02", k)
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Before(out[j]) })
	return out
}
