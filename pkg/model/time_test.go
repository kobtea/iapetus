package model

import (
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
	"strconv"
	"testing"
	"time"
)

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
			DurationCriteria{"<", model.Duration(24 * time.Hour)},
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
			DurationCriteria{">", model.Duration(time.Hour)},
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

	dc := DurationCriteria{">", model.Duration(10 * time.Hour)}
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
			DurationCriteria{">", model.Duration(24 * time.Hour)},
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
