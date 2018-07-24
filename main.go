/*
Copyright (C) 2018 Morandini Daniel

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

package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type Job struct {
	ID    string // job identifier
	Urls  []string // list for download tasks
	Delay int // expressed in milliseconds
	Sync  bool // wether the download tasks should be performed sync or async
}

// Wait waits for the duration of j.Delay, coverted into a time duration expressed
// in milliseconds.
func (j *Job) Wait() {
	d, _ := time.ParseDuration(fmt.Sprintf("%dms", j.Delay))
	<-time.After(d)
}

var proxyAddr = flag.String("proxy", "", "optional proxy address")
var input = flag.String("job", "job.json", "input job file formatted in json")

func main() {
	flag.Parse()
	flag.Usage()

	fmt.Printf("\nParsing input file %v...\n", *input)
	file, err := os.Open(*input)
	if err != nil {
		panic(err)
	}

	var jobs []*Job
	if err = json.NewDecoder(file).Decode(&jobs); err != nil {
		panic(err)
	}
	fmt.Printf("Jobs count: %v\n", len(jobs))

	fmt.Printf("Test is starting...\n\n")

	t := &http.Transport{
		Proxy: func(*http.Request) (*url.URL, error) {
			if *proxyAddr != "" {
				return url.Parse(*proxyAddr)
			}
			return nil, nil
		},
		DisableCompression: false,
		DisableKeepAlives:  true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &Client{
		&http.Client{Transport: t},
	}

	start := time.Now()
	fmt.Printf("Start time: %v\n", start)

	for _, j := range jobs {
		client.HandleJob(j)
	}

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Printf("\nEnd time: %v, elapsed: %v\n", end, elapsed)
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

func (c *Client) FetchAndDiscard(addr string) {
	start := time.Now()
	log.Printf("[%v] fetching...\n", addr)

	resp, err := c.Get(addr)
	if err != nil {
		fmt.Printf("FetchAndDiscard(%v): %v\n", addr, err)
		return
	}

	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	d := time.Now().Sub(start)
	log.Printf("[%v] done in %dns (%v).\n", addr, d.Nanoseconds(), d)
}
