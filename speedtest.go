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

// Job describes a speedtest job.
type Job struct {
	ID    string   `json:"id"`    // job identifier
	Urls  []string `json:"urls"`  // list for download tasks
	Delay int      `json:"delay"` // expressed in milliseconds
}

// Wait waits for the duration of j.Delay, coverted into a time duration expressed
// in milliseconds.
func (j *Job) Wait() {
	d, _ := time.ParseDuration(fmt.Sprintf("%dms", j.Delay))
	<-time.After(d)
}

// Client is a wrapper around an http client.
type Client struct {
	*http.Client
}

// HandleJob calls FetchAndDiscard for each address in the url list, printing the
// results to standard output.
func (c *Client) HandleJob(job *Job) {
	start := time.Now()
	fmt.Printf("[%v] Handling job...\n", job.ID)

	var wg sync.WaitGroup
	for _, v := range job.Urls {
		wg.Add(1)

		go func(addr string) {
			defer wg.Done()

			res, err := c.FetchAndDiscard(addr)
			if err != nil {
				fmt.Printf("[%v] error: %v\n", job.ID, err)
				return
			}

			fmt.Printf("[%v] Downloaded: %v (elapsed: %v, bandwidth: %v)\n", job.ID, addr, res.ElapsedTime, res.Bandwidth().String())
		}(v)

		// wait before firing the next job
		job.Wait()
	}
	wg.Wait()

	fmt.Printf("[%v] Done (%v).\n", job.ID, time.Now().Sub(start))
}

// Result contains information about a download task.
type Result struct {
	Start         time.Time
	End           time.Time
	ElapsedTime   time.Duration
	ContentLength int64
}

type bandwidth float64

// Bandwidth returns the number of bytes transferred per second.
func (r *Result) Bandwidth() bandwidth {
	return bandwidth(float64(r.ContentLength) / float64(r.ElapsedTime.Seconds()))
}

func (b bandwidth) String() string {
	return fmt.Sprintf("%.2fmb/s", b)
}

// FetchAndDiscard performs a GET request, returns an error if the request is
// not successful, otherwise returns a Result, containing metrics about the
// request performed.
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
