package util

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestUtilTimeParseDuration(t *testing.T) {
	test := is.New(t)
	for _, testCase := range []struct {
		representation string
		object         time.Duration
	}{
		{"1h", time.Duration(time.Hour)},
		{"1m", time.Duration(time.Minute)},
		{"1s", time.Duration(time.Second)},
		{"1ms", time.Duration(time.Millisecond)},
		{"1µs", time.Duration(time.Microsecond)},
		{"1us", time.Duration(time.Microsecond)},
		{"1ns", time.Duration(time.Nanosecond)},
		{"4.000000001s", time.Duration(4*time.Second + time.Nanosecond)},
		{"1h0m4.000000001s", time.Duration(time.Hour + 4*time.Second + time.Nanosecond)},
		{"1h1m0.01s", time.Duration(61*time.Minute + 10*time.Millisecond)},
		{"1h1m0.123456789s", time.Duration(61*time.Minute + 123456789*time.Nanosecond)},
		{"1.00002ms", time.Duration(time.Millisecond + 20*time.Nanosecond)},
		{"1.00000002s", time.Duration(time.Second + 20*time.Nanosecond)},
		{"693ns", time.Duration(693 * time.Nanosecond)},
		{"10s1us693ns", time.Duration(10*time.Second + time.Microsecond + 693*time.Nanosecond)},
		{"1ms1ns", time.Duration(time.Millisecond + 1*time.Nanosecond)},
		{"1s20ns", time.Duration(time.Second + 20*time.Nanosecond)},
		{"60h8ms", time.Duration(60*time.Hour + 8*time.Millisecond)},
		{"96h63s", time.Duration(96*time.Hour + 63*time.Second)},
		{"2d3s96ns", time.Duration(2*time.Hour*dayHours + 3*time.Second + 96*time.Nanosecond)},
		{"1w2d3s96ns", time.Duration(time.Hour*weekHours + 2*time.Hour*dayHours + 3*time.Second + 96*time.Nanosecond)},
		{"1w2d3s3µs96ns", time.Duration(time.Hour*weekHours + 2*time.Hour*dayHours + 3*time.Second + 3*time.Microsecond + 96*time.Nanosecond)},
		{"1M2d3s3µs96ns", time.Duration(time.Hour*monthHours + 2*time.Hour*dayHours + 3*time.Second + 3*time.Microsecond + 96*time.Nanosecond)},
		{"3y1M3d10s1us693ns", time.Duration(3*time.Hour*yearHours + time.Hour*monthHours + 3*time.Hour*dayHours + +10*time.Second + time.Microsecond + 693*time.Nanosecond)},
	} {
		duration, err := ParseDuration(testCase.representation)
		test.NoErr(err)
		test.Equal(testCase.object, duration)
	}
}

func TestUtilTimeParseDurationString(t *testing.T) {
	test := is.New(t)
	for _, duration := range []time.Duration{
		time.Duration(time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Minute + time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Second + time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Millisecond + time.Microsecond + time.Nanosecond),
		time.Duration(time.Microsecond + time.Nanosecond),
		time.Duration(time.Nanosecond),
		time.Duration(time.Hour + time.Minute + time.Second + time.Millisecond + time.Microsecond),
		time.Duration(time.Minute + time.Second + time.Millisecond + time.Microsecond),
		time.Duration(time.Second + time.Millisecond + time.Microsecond),
		time.Duration(time.Millisecond + time.Microsecond),
		time.Duration(time.Microsecond),
		time.Duration(time.Hour + time.Minute + time.Second + time.Millisecond),
		time.Duration(time.Minute + time.Second + time.Millisecond),
		time.Duration(time.Second + time.Millisecond),
		time.Duration(time.Millisecond),
		time.Duration(time.Hour + time.Minute + time.Second),
		time.Duration(time.Minute + time.Second),
		time.Duration(time.Second),
		time.Duration(time.Hour + time.Minute),
		time.Duration(time.Minute),
		time.Duration(time.Hour),
		time.Duration(time.Millisecond + time.Nanosecond),
		time.Duration(1001001 * time.Microsecond),
		time.Duration(1001 * time.Nanosecond),
		time.Duration(61 * time.Minute),
		time.Duration(61 * time.Second),
		time.Duration(time.Microsecond + 16*time.Nanosecond),
	} {
		durationFromString, _ := ParseDuration(duration.String())
		test.Equal(duration.String(), durationFromString.String())
	}
}
