// Copyright 2016 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"strings"
	"github.com/prometheus/common/model"
)

func ParseTime(s string) (time.Time, error) {
	if t, err := parseTime(s); err == nil {
		return t, nil
	}
	if elms := strings.Split(s, "now-"); len(elms) == 2 {
		// expect duration
		d, err := ParseDuration(elms[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("cannot parse %q to a valid duration", elms[1])
		}
		return time.Now().Add(-d), nil
	}
	return time.Time{}, fmt.Errorf("cannot parse %q to a valid timestamp", s)
}

// origin: https://github.com/prometheus/prometheus/blob/v2.2.1/web/api/v1/api.go#L798
func parseTime(s string) (time.Time, error) {
	if t, err := strconv.ParseFloat(s, 64); err == nil {
		s, ns := math.Modf(t)
		return time.Unix(int64(s), int64(ns*float64(time.Second))), nil
	}
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("cannot parse %q to a valid timestamp", s)
}

// origin: https://github.com/prometheus/prometheus/blob/v2.2.1/web/api/v1/api.go#L809
// modify: make the scope of `parseDuration` to public.
func ParseDuration(s string) (time.Duration, error) {
	if d, err := strconv.ParseFloat(s, 64); err == nil {
		ts := d * float64(time.Second)
		if ts > float64(math.MaxInt64) || ts < float64(math.MinInt64) {
			return 0, fmt.Errorf("cannot parse %q to a valid duration. It overflows int64", s)
		}
		return time.Duration(ts), nil
	}
	if d, err := model.ParseDuration(s); err == nil {
		return time.Duration(d), nil
	}
	return 0, fmt.Errorf("cannot parse %q to a valid duration", s)
}
