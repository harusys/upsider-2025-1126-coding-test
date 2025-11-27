package ctxutiltest

import (
	"testing"
	"time"
)

const (
	// JSTOffset is the UTC offset for Asia/Tokyo timezone (UTC+9)
	jstOffset = 9 * time.Hour
)

func (p *TestContextProvider) Now() time.Time {
	if p.CurrentTime == nil {
		return time.Now() //nolint:forbidigo // 空の時刻を返すことで問題を見落とすリスクを考慮して、常になんらかの時刻を返したい
	}

	return *p.CurrentTime
}

func (p *TestContextProvider) SetAsiaTokyo(t *testing.T, input string) time.Time {
	t.Helper()

	now := parse(t, input, jstOffset, "Asia/Tokyo")
	p.CurrentTime = &now

	return now
}

func parse(t *testing.T, input string, diff time.Duration, location string) time.Time {
	t.Helper()

	loc, err := time.LoadLocation(location)
	if err != nil {
		t.Fatalf("Failed to load location: %v", err)
	}

	format := "2006-01-02 15:04:05"

	tm, err := time.Parse(format, input)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}

	return tm.Add(-diff).In(loc)
}
