package time

import (
	"github.com/peter-mount/go-script/packages"
	"time"
)

func init() {
	packages.Register("time", &Time{
		January:     time.January,
		February:    time.February,
		March:       time.March,
		April:       time.April,
		May:         time.May,
		June:        time.June,
		July:        time.July,
		August:      time.August,
		September:   time.September,
		October:     time.October,
		November:    time.November,
		December:    time.December,
		Sunday:      time.Sunday,
		Monday:      time.Monday,
		Tuesday:     time.Tuesday,
		Wednesday:   time.Wednesday,
		Thursday:    time.Thursday,
		Friday:      time.Friday,
		Saturday:    time.Saturday,
		Nanosecond:  time.Nanosecond,
		Microsecond: time.Microsecond,
		Millisecond: time.Millisecond,
		Second:      time.Second,
		Minute:      time.Minute,
		Hour:        time.Hour,
		Layout:      time.Layout,
		ANSIC:       time.ANSIC,
		UnixDate:    time.UnixDate,
		RubyDate:    time.RubyDate,
		RFC822:      time.RFC822,
		RFC822Z:     time.RFC822Z,
		RFC850:      time.RFC850,
		RFC1123:     time.RFC1123,
		RFC1123Z:    time.RFC1123Z,
		RFC3339:     time.RFC3339,
		RFC3339Nano: time.RFC3339Nano,
		Kitchen:     time.Kitchen,
		Stamp:       time.Stamp,
		StampMilli:  time.StampMilli,
		StampMicro:  time.StampMicro,
		StampNano:   time.StampNano,
	})
	time.Now()
}

type Time struct {
	January     time.Month
	February    time.Month
	March       time.Month
	April       time.Month
	May         time.Month
	June        time.Month
	July        time.Month
	August      time.Month
	September   time.Month
	October     time.Month
	November    time.Month
	December    time.Month
	Sunday      time.Weekday
	Monday      time.Weekday
	Tuesday     time.Weekday
	Wednesday   time.Weekday
	Thursday    time.Weekday
	Friday      time.Weekday
	Saturday    time.Weekday
	Nanosecond  time.Duration
	Microsecond time.Duration
	Millisecond time.Duration
	Second      time.Duration
	Minute      time.Duration
	Hour        time.Duration
	Layout      string
	ANSIC       string
	UnixDate    string
	RubyDate    string
	RFC822      string
	RFC822Z     string
	RFC850      string
	RFC1123     string
	RFC1123Z    string
	RFC3339     string
	RFC3339Nano string
	Kitchen     string
	Stamp       string
	StampMilli  string
	StampMicro  string
	StampNano   string
}

func (_ Time) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (_ Time) Until(t time.Time) time.Duration {
	return time.Until(t)
}

func (_ Time) Now() time.Time {
	return time.Now()
}

func (_ Time) Unix(sec, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func (_ Time) UnixMilli(msec int64) time.Time {
	return time.UnixMilli(msec)
}

func (_ Time) UnixMicro(usec int64) time.Time {
	return time.UnixMicro(usec)
}

func (_ Time) Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, loc)
}

func (_ Time) LoadLocation(name string) (*time.Location, error) {
	return time.LoadLocation(name)
}

func (_ Time) Parse(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func (_ Time) ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, value, loc)
}

func (_ Time) ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}
