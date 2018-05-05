package model

import (
	"gopkg.in/yaml.v2"
	"strconv"
	"testing"
	"time"
)

func TestNewTimeCriteria(t *testing.T) {
	tests := []struct {
		s     string
		v     TimeCriteria
		error bool
	}{
		{
			"hoge",
			TimeCriteria{},
			true,
		},
		{
			"<",
			TimeCriteria{},
			true,
		},
		{
			"- 1000",
			TimeCriteria{},
			true,
		},
		{
			"< 1000",
			TimeCriteria{"<", time.Unix(1000, 0)},
			false,
		},
	}
	for _, test := range tests {
		v, err := NewTimeCriteria(test.s)
		if test.error && err == nil {
			t.Errorf("expect error, but don't")
		}
		if v != test.v {
			t.Errorf("expect %v, but got %v", test.v, v)
		}
	}
}

func TestTimeCriteria_IsZero(t *testing.T) {
	tests := []struct {
		c TimeCriteria
		v bool
	}{
		{
			TimeCriteria{},
			true,
		},
		{
			TimeCriteria{">", time.Unix(1000, 0)},
			false,
		},
	}
	for _, test := range tests {
		if test.c.IsZero() != test.v {
			t.Errorf("expect %v, but don't", test.v)
		}
	}
}

func TestTimeCriteria_Satisfy(t *testing.T) {
	tests := []struct {
		ts time.Time
		v  bool
	}{
		{
			time.Unix(2000, 0),
			false,
		},
		{
			time.Unix(500, 0),
			true,
		},
	}
	tc := TimeCriteria{"<", time.Unix(1000, 0)}
	for _, test := range tests {
		if tc.Satisfy(test.ts) != test.v {
			t.Errorf("expect %v, but don't", test.v)
		}
	}
}

func TestTimeCriteria_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		s     string
		v     TimeCriteria
		error bool
	}{
		{
			"hoge",
			TimeCriteria{},
			true,
		},
		{
			"< 1000",
			TimeCriteria{"<", time.Unix(1000, 0)},
			false,
		},
	}
	for _, test := range tests {
		var v TimeCriteria
		e := yaml.Unmarshal([]byte(strconv.Quote(test.s)), &v)
		if test.error && e == nil {
			t.Errorf("expect error, but don't")
		}
		if v != test.v {
			t.Errorf("expect %v, but got %v", test.v, v)
		}
	}
}

func TestNewDurationCriteria(t *testing.T) {
	tests := []struct {
		s     string
		v     DurationCriteria
		error bool
	}{
		{
			"1",
			DurationCriteria{},
			true,
		},
		{
			"<",
			DurationCriteria{},
			true,
		},
		{
			"~ 1d",
			DurationCriteria{},
			true,
		},
		{
			"< 1d",
			DurationCriteria{"<", 24 * time.Hour},
			false,
		},
	}
	for _, test := range tests {
		v, err := NewDurationCriteria(test.s)
		if test.error && err == nil {
			t.Errorf("expect error, but don't")
		}
		if v != test.v {
			t.Errorf("expect %v, but got %v", test.v, v)
		}
	}
}

func TestDurationCriteria_IsZero(t *testing.T) {
	tests := []struct {
		c DurationCriteria
		v bool
	}{
		{
			DurationCriteria{},
			true,
		},
		{
			DurationCriteria{">", time.Hour},
			false,
		},
	}
	for _, test := range tests {
		if test.c.IsZero() != test.v {
			t.Errorf("expect %v, but don't", test.v)
		}
	}
}

func TestDurationCriteria_Satisfy(t *testing.T) {
	tests := []struct {
		start time.Time
		end   time.Time
		v     bool
	}{
		{
			time.Unix(0, 0),
			time.Unix(0, 0).Add(2 * time.Hour),
			false,
		},
		{
			time.Unix(0, 0).Add(2 * time.Hour),
			time.Unix(0, 0),
			false,
		},
		{
			time.Unix(0, 0),
			time.Unix(0, 0).Add(20 * time.Hour),
			true,
		},
	}

	dc := DurationCriteria{">", 10 * time.Hour}
	for _, test := range tests {
		if dc.Satisfy(test.start, test.end) != test.v {
			t.Errorf("expect %v, but don't", test.v)
		}
	}
}

func TestDurationCriteria_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		s     string
		v     DurationCriteria
		error bool
	}{
		{
			"10",
			DurationCriteria{},
			true,
		},
		{
			"> 1d",
			DurationCriteria{">", 24 * time.Hour},
			false,
		},
	}
	for _, test := range tests {
		var v DurationCriteria
		e := yaml.Unmarshal([]byte(strconv.Quote(test.s)), &v)
		if test.error && e == nil {
			t.Errorf("expect error, but don't")
		}
		if v != test.v {
			t.Errorf("expect %v, but got %v", test.v, v)
		}
	}
}
