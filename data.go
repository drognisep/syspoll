package main

import (
	"net/url"
	"time"
)

type System struct {
	Name          string      `json:"system"`
	CheckInterval string      `json:"interval"`
	FailedChecks  []time.Time `json:"failedChecks,omitempty"`
	Http          *CheckHttp  `json:"http,omitempty"`
}

func (s *System) Copy() *System {
	if s == nil {
		return nil
	}
	cp := &System{
		Name:          s.Name,
		CheckInterval: s.CheckInterval,
		Http:          s.Http.Copy(),
	}
	cp.FailedChecks = make([]time.Time, len(s.FailedChecks))
	copy(cp.FailedChecks, s.FailedChecks)
	return cp
}

func (s *System) Interval() (time.Duration, error) {
	return time.ParseDuration(s.CheckInterval)
}

type CheckHttp struct {
	URL string `json:"url"`
}

func (h *CheckHttp) Copy() *CheckHttp {
	if h == nil {
		return nil
	}
	return &CheckHttp{
		URL: h.URL,
	}
}

func (h *CheckHttp) ToURL() (*url.URL, error) {
	return url.Parse(h.URL)
}
