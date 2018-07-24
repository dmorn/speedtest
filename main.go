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
	ID    string
	Urls  []string
	Delay time.Duration
	Sync  bool
}

var input = flag.String("job", "job.json", "input job file formatted in json")
var verbose = flag.Bool("verbose", false, "enables verbose mode")

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
			return url.Parse("socks5://localhost:1080")
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
		<-time.After(job.Delay)
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
