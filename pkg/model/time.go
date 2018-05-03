package model

import (
	"fmt"
	"github.com/prometheus/common/model"
	"strings"
	"time"
)

var operators = []string{"<", ">"}

type DurationCriteria struct {
	Op       string
	Duration model.Duration
}

func NewDurationCriteria(s string) (DurationCriteria, error) {
	elms := strings.Split(strings.TrimSpace(s), " ")
	if len(elms) != 2 {
		return DurationCriteria{}, fmt.Errorf("invalid duration criteria: expects `<op> <duration>`")
	}
	// check operand
	d, err := model.ParseDuration(elms[1])
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
	sub := model.Duration(end.Sub(start))
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
