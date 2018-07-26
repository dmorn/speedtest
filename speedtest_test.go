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
	"fmt"
	"testing"
	"time"

	"github.com/tecnoporto/speedtest"
)

func TestResult_Bandwidth(t *testing.T) {
	tt := []struct {
		res *speedtest.Result
		out float64
	}{
		{res: &speedtest.Result{ElapsedTime: time.Second*3, ContentLength: 300}, out: 100.0},
	}

	for i, v := range tt {
		r := v.res.Bandwidth()
		if r != v.out {
			fmt.Printf("%+v\n", v.res)
			t.Fatalf("%d: wanted %.3f, found %.3f", i, v.out, r)
		}
	}
}
