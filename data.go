package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

var _ json.Marshaler = (*FailureType)(nil)
var _ json.Unmarshaler = (*FailureType)(nil)

type FailureType int

const (
	Unknown FailureType = iota
	Down
	Error
)

func (f *FailureType) UnmarshalJSON(bytes []byte) error {
	if len(bytes) < 3 {
		return fmt.Errorf("invalid input: %v", bytes)
	}
	if string(bytes) == "null" {
		*f = Unknown
		return nil
	}
	runes := []rune(string(bytes))
	runes = runes[1 : len(runes)-1]
	s := string(runes)
	newF, ok := failureMap[s]
	if !ok {
		*f = Unknown
		return nil
	}
	*f = newF
	return nil
}

func (f *FailureType) MarshalJSON() ([]byte, error) {
	if f == nil {
		return []byte(fmt.Sprintf(`"%s"`, failureStrings[Unknown])), nil
	}
	if failureTypes[*f] {
		return []byte(fmt.Sprintf(`"%s"`, failureStrings[*f])), nil
	}
	return nil, errors.New("unrecognized failure type")
}

var failureTypes = map[FailureType]bool{Unknown: true, Down: true, Error: true}
var failureStrings = map[FailureType]string{
	Unknown: "unknown",
	Down:    "down",
	Error:   "error",
}
var failureMap = map[string]FailureType{
	"unknown": Unknown,
	"down":    Down,
	"error":   Error,
}

func (f FailureType) UiString() string {
	switch f {
	case Unknown:
		return "[gray]UNK"
	case Down:
		return "[red]DOWN"
	case Error:
		return "[red]ERROR"
	}
	return ""
}

type FailEvent struct {
	Type FailureType `json:"type"`
	Time time.Time   `json:"time"`
}

func DownFailure(t time.Time) FailEvent {
	return FailEvent{
		Down, t,
	}
}

func ErrorFailure(t time.Time) FailEvent {
	return FailEvent{
		Error, t,
	}
}

type System struct {
	Name          string      `json:"system"`
	CheckInterval string      `json:"interval"`
	FailedChecks  []FailEvent `json:"failedChecks,omitempty"`
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
	cp.FailedChecks = make([]FailEvent, len(s.FailedChecks))
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
