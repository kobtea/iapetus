package model

import (
	"fmt"
	"github.com/kobtea/iapetus/pkg/util"
	"strings"
	"time"
)

var operators = []string{"<", ">"}

type TimeCriteria struct {
	Op   string
	Time time.Time
}

func NewTimeCriteria(s string) (TimeCriteria, error) {
	elms := strings.Split(strings.TrimSpace(s), " ")
	if len(elms) != 2 {
		return TimeCriteria{}, fmt.Errorf("invalid time criteria: expects `<op> <time>`")
	}
	// check operand
	t, err := util.ParseTime(elms[1])
	if err != nil {
		return TimeCriteria{}, fmt.Errorf("invalid time criteria: %s", elms[1])
	}
	// check operator
	for _, op := range operators {
		if op == elms[0] {
			return TimeCriteria{Op: elms[0], Time: t}, nil
		}
	}
	return TimeCriteria{}, fmt.Errorf("invalid time criteria: %s", elms[0])
}

func (t *TimeCriteria) IsZero() bool {
	return *t == TimeCriteria{}
}

func (t *TimeCriteria) Satisfy(ts time.Time) bool {
	switch t.Op {
	case "<":
		return ts.Before(t.Time)
	case ">":
		return ts.After(t.Time)
	default:
		panic(fmt.Sprintf("invalid time criteria: %s", t.Op))
	}
}

func (t *TimeCriteria) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return fmt.Errorf("failed to parse time criteria")
	}
	tc, err := NewTimeCriteria(s)
	if err != nil {
		return fmt.Errorf("failed to parse time criteria")
	}
	*t = tc
	return nil
}

type DurationCriteria struct {
	Op       string
	Duration time.Duration
}

func NewDurationCriteria(s string) (DurationCriteria, error) {
	elms := strings.Split(strings.TrimSpace(s), " ")
	if len(elms) != 2 {
		return DurationCriteria{}, fmt.Errorf("invalid duration criteria: expects `<op> <duration>`")
	}
	// check operand
	d, err := util.ParseDuration(elms[1])
	if err != nil {
		return DurationCriteria{}, fmt.Errorf("invalid duration criteria: %s", elms[1])
	}
	// check operator
	for _, op := range operators {
		if op == elms[0] {
			return DurationCriteria{Op: elms[0], Duration: d}, nil
		}
	}
	return DurationCriteria{}, fmt.Errorf("invalid duration criteria: %s", elms[0])
}

func (d *DurationCriteria) IsZero() bool {
	return *d == DurationCriteria{}
}

func (d *DurationCriteria) Satisfy(start, end time.Time) bool {
	sub := end.Sub(start)
	switch d.Op {
	case "<":
		return sub < d.Duration
	case ">":
		return sub > d.Duration
	default:
		panic(fmt.Sprintf("invalid duration criteria: %s", d.Op))
	}
}

func (d *DurationCriteria) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return fmt.Errorf("failed to parse duration criteria")
	}
	dd, err := NewDurationCriteria(s)
	if err != nil {
		return fmt.Errorf("failed to parse duration criteria")
	}
	*d = dd
	return nil
}
