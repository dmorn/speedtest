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

package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/tecnoporto/speedtest"
)

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

	var jobs []*speedtest.Job
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
	client := &speedtest.Client{
		&http.Client{Transport: t},
	}

	start := time.Now()
	fmt.Printf("Start time: %v\n\n", start)

	for _, j := range jobs {
		client.HandleJob(j)
	}

	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Printf("\nEnd time: %v, elapsed: %v\n", end, elapsed)
}
