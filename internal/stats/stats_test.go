/*
 * Copyright (c) 2024, Jeremy Bingham (<jeremy@goiardi.gl>)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package stats

import (
	"github.com/ctdk/morseudar/internal/morse"
	"os"
	"testing"
	"time"
)

func TestSaveStatFile(t *testing.T) {
	f, err := os.CreateTemp("", "stat-test")
	if err != nil {
		t.Errorf("error creating test stat file: %s", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	u := New()
	err = u.Save(f.Name())
	if err != nil {
		t.Errorf("error saving file: %s", err)
	}
}

func TestLoadStatFile(t *testing.T) {
	f, err := os.CreateTemp("", "stat-test")
	if err != nil {
		t.Errorf("error creating test stat file: %s", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	u := New()
	origTime := u.Created
	err = u.Save(f.Name())
	if err != nil {
		t.Errorf("error saving file: %s", err)
	}

	u2, err := Load(f.Name())
	if err != nil {
		t.Errorf("error loading stat file: %s", err)
	}

	if !origTime.Equal(u2.Created) {
		t.Errorf("mismatch in stat file time: created should have been %s, got %s instead", origTime, u2.Created)
	}
}

func TestSummaries(t *testing.T) {
	f, err := os.CreateTemp("", "stat-test")
	if err != nil {
		t.Errorf("error creating test stat file: %s", err)
	}
	f.Close()
	defer os.Remove(f.Name())
	u := New()

	s1 := NewSummary(time.Now(), morse.TopWords, 76.4, time.Minute * 15, 2.4, 28, 10, 0)
	time.Sleep(5 * time.Second)
	s2 := NewSummary(time.Now(), morse.Qcode, 42.4, time.Minute * 5, 5.4, 10, 10, 0)
	time.Sleep(15 * time.Second)
	s3 := NewSummary(time.Now(), morse.TopWords, 96.4, time.Minute * 1, 1.1, 5, 10, 0)

	u.Add(s1)
	u.Add(s2)
	u.Add(s3)

	if len(u.Summaries) != 3 {
		t.Errorf("wrong length for summaries: should have been 3, got %d instead", len(u.Summaries))
	}
	if u.Created.Equal(u.Updated) {
		t.Errorf("user stats failed to update the updated time")
	}

	err = u.Save(f.Name())
	if err != nil {
		t.Errorf("error saving file: %s", err)
	}

	u2, err := Load(f.Name())
	if err != nil {
		t.Errorf("error loading stat file: %s", err)
	}
	if len(u2.Summaries) != 3 {
		t.Errorf("wrong length for loaded summaries: should have been 3, got %d instead", len(u2.Summaries))
	}
	if u2.Created.Equal(u2.Updated) {
		t.Errorf("loaded user stats failed to update the updated time")
	}

	if !u2.Created.Equal(u.Created) {
		t.Errorf("loaded user stat creation time not equal the original")
	}
	if !u2.Summaries[0].Date.Equal(s1.Date) {
		t.Errorf("s1 date not equal to u2.Summaries[0] date loaded from disk")
	}
}
