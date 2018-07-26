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

package speedtest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	ID    string   // job identifier
	Urls  []string // list for download tasks
	Delay int      // expressed in milliseconds
	Sync  bool     // wether the download tasks should be performed sync or async
}

// Wait waits for the duration of j.Delay, coverted into a time duration expressed
// in milliseconds.
func (j *Job) Wait() {
	d, _ := time.ParseDuration(fmt.Sprintf("%dms", j.Delay))
	<-time.After(d)
}


type Client struct {
	*http.Client
}

func (c *Client) HandleJob(job *Job) {
	start := time.Now()
	fmt.Printf("[%v] handling job...\n", job.ID)
	if job.Sync {
		c.handleJobSync(job)
	} else {
		c.handleJobAsync(job)
	}

	fmt.Printf("[%v] done (%v).\n", job.ID, time.Now().Sub(start))
}

func (c *Client) handleJobSync(job *Job) {
	for _, v := range job.Urls {
		c.FetchAndDiscard(v)
	}
}

func (c *Client) handleJobAsync(job *Job) {
	var wg sync.WaitGroup
	for _, v := range job.Urls {
		wg.Add(1)

		go func(addr string) {
			defer wg.Done()
			c.FetchAndDiscard(addr)
		}(v)

		// wait before firing the next job
		job.Wait()
	}
	wg.Wait()
}

type Result struct {
	Start         time.Time
	End           time.Time
	ElapsedTime   time.Duration
	ContentLength int64
}

func (r *Result) Bandwidth() float64 {
	return float64(int64(r.ContentLength) / int64(r.ElapsedTime.Seconds()))
}

func (c *Client) FetchAndDiscard(addr string) (*Result, error) {
	start := time.Now()

	resp, err := c.Get(addr)
	if err != nil {
		return nil, err
	}
	end := time.Now()

	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	return &Result{
		Start:         start,
		End:           end,
		ElapsedTime:   end.Sub(start),
		ContentLength: resp.ContentLength,
	}, nil
}
