package speedtest_test

import (
	"testing"
	"time"

	"github.com/tecnoporto/speedtest"
)

func TestResult_Bandwidth(t *testing.T) {
	tt := []struct {
		res speedtest.Result
		u   time.Duration
		out float64
	}{}

	for _, v := range tt {
		r := v.res.Bandwidth(v.u)
		if r != v.out {
			t.Fatalf("wanted %.3f, found %.3f", v.out, r)
		}
	}
}
