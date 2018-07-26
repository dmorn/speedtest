/*
Copyright (C) 2018 Daniel Morandini

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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
