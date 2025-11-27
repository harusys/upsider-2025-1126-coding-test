package timeutil_test

import (
	"testing"

	"github.com/harusys/super-shiharai-kun/pkg/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestAsiaTokyo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		wantUnixtime int64
		wantTimezone string
	}{
		{
			name:         "2024-07-31 16:12:57 は、Unixtimeの 1722409977です",
			input:        "2024-07-31 16:12:57",
			wantUnixtime: 1722409977,
			wantTimezone: "Asia/Tokyo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := timeutil.AsiaTokyo(t, tt.input)
			assert.Equal(t, tt.wantUnixtime, got.Unix(), "Unixtimeが正しいこと")
			assert.Equal(t, tt.wantTimezone, got.Location().String(), "タイムゾーンが正しいこと")
		})
	}
}
