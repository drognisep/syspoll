package main

import (
	"net/url"
	"time"
)

type System struct {
	Name          string      `json:"system"`
	CheckInterval string      `json:"interval"`
	FailedChecks  []time.Time `json:"-"`
	Http          *CheckHttp  `json:"http,omitempty"`
}

func (s *System) Interval() (time.Duration, error) {
	return time.ParseDuration(s.CheckInterval)
}

type CheckHttp struct {
	URL string `json:"url"`
}

func (h *CheckHttp) ToURL() (*url.URL, error) {
	return url.Parse(h.URL)
}
