package timekit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSleep(t *testing.T) {
	testcases := []struct {
		Duration string
		seconds  int64
	}{
		{"1s", 1},
		{"2s", 2},
	}

	for _, tt := range testcases {
		start := time.Now()
		Sleep(tt.Duration)

		seconds := time.Now().Unix() - start.Unix()
		require.Equal(t, tt.seconds, seconds)
	}
}

func TestEqualUTC(t *testing.T) {
	testcases := []struct {
		t        time.Time
		s        string
		expected bool
	}{
		{UTC("2015-11-10T23:00:00Z"), "2015-11-10T23:00:00Z", true},
		{UTC("2015-12-10T23:00:00Z"), "2015-11-10T23:00:00Z", false},
	}

	for _, tt := range testcases {
		require.Equal(t, tt.expected, EqualUTC(tt.t, tt.s))
	}
}

func TestEqualString(t *testing.T) {
	testcases := []struct {
		t        time.Time
		format   string
		s        string
		expected bool
	}{
		{Parse("2006-01-02T15:04:05", "2018-10-10T20:04:30"), "2006-01-02T15:04:05", "2018-10-10T20:04:30", true},
		{Parse("2006-01-02T15:04:05", "2018-10-10T20:04:30"), "2006-01-02T15:04:05", "2018-10-11T20:04:30", false},
	}

	for _, tt := range testcases {
		require.Equal(t, tt.expected, EqualString(tt.t, tt.format, tt.s))
	}
}

func TestDuration_WrongFormat(t *testing.T) {
	testcase := []struct {
		s        string
		expected string
	}{
		{"wrong-format", "0s"},
		{"13s", "13s"},
	}

	for _, tt := range testcase {
		d := Duration(tt.s)
		require.Equal(t, tt.expected, d.String())
	}
}
