// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	numPollers     = 10
	pollInterval   = 2 * time.Second
	statusInterval = 2 * time.Second
	errTimeout     = 2 * time.Second
)

var urls = []string{
	"https://google.com",
	"https://vnexpress.net",
	"https://kenh14.vn",
	"https://facebook.com",
}

// State represent last-known status of an URL

type State struct {
	url    string
	status string
}

// Represent resource to be polled
type Resource struct {
	url      string
	errCount int
}

func logState(status map[string]string) {
	for k, v := range status {
		fmt.Println(k, v)
	}
}

func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)
	if err != nil {
		r.errCount += 1
		fmt.Println(err)
		return err.Error()
	}
	r.errCount = 0
	return resp.Status
}

func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(pollInterval)
	done <- r
}

func Poller(id int, in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	fmt.Println("Poller id =", id)
	for {
		r := <-in
		fmt.Printf("Poller id %+v, r.url = %+v\n", id, r.url)
		s := r.Poll()
		status <- State{url: r.url, status: s}
		out <- r
	}
}

// StateMonitor maintains a map that store the states of the URLs being
// polled, and prints the current state every updateInterval seconds

func StateMonitor(duration time.Duration) chan<- State {
	updates := make(chan State)
	urlStatus := make(map[string]string)
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ticker.C:
				logState(urlStatus)
			case s := <-updates:
				urlStatus[s.url] = s.status
			}
		}
	}()
	return updates

}

func main() {
	pending, complete := make(chan *Resource), make(chan *Resource)
	status := StateMonitor(statusInterval)

	// Launch some pollers go routine
	for i := 1; i <= numPollers; i++ {
		go Poller(i, pending, complete, status)
	}

	// send some resource to the pending queue
	go func() {
		for _, v := range urls {
			pending <- &Resource{url: v}
		}
	}()

	for r := range complete {
		go r.Sleep(pending)
	}
}
