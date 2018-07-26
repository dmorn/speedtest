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
	ContentLength int64
}

func (r *Result) ElapsedTime() time.Duration {
	return r.End.Sub(r.Start)
}

func (r *Result) Bandwidth(u time.Duration) float64 {
	return 10.0
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
		ContentLength: resp.ContentLength,
	}, nil
}
