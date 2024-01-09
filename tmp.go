/*
Copyright 2024 The gpt4batch Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gpt4batch

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// MonthDayYear returns the current month and day in the format MMDD.
func MonthDayYear() string {
	now := time.Now()
	return fmt.Sprintf("%04d%02d%02d", now.Year(), now.Month(), now.Day())
}

// HourMinute returns the current hour and minute in the format HHMM.
func HourMinute() string {
	now := time.Now()
	ts := now.Format("15:04")
	return ts[:2] + ts[3:]
}

// TempFileName returns a temporary file name in the format
func TempFileName(name string) string {
	md := MonthDayYear()
	ts := HourMinute()

	name = filepath.Base(name)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	ps := make([]string, 0, 3)
	ps = append(ps, md)
	ps = append(ps, ts)
	ps = append(ps, name)
	js := strings.Join(ps, "-")
	return fmt.Sprintf(".%s.jsonl", js)
}
