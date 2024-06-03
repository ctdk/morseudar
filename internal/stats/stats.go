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
	"encoding/gob"
	"fmt"
	"github.com/ctdk/morseudar/internal/morse"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"
)

const StatVersion = "0.1.0"

// A bit empty at the moment, but ought to become more detailed at a later time
// I think.
type UserStats struct {
	Username string
	Summaries []Summary
	Version string
	Created time.Time
	Updated time.Time
	saveFilePath string
}

type Summary struct {
	Date time.Time
	Mode morse.MorseMode
	AvgPerc float64
	AvgDur time.Duration
	AvgTries float64
	Count int
	Wpm int
	Farnsworth int
}

func NewSummary(date time.Time, mode morse.MorseMode, perc float64, dur time.Duration, tries float64, count int, wpm int, farns int) Summary {
	return Summary{Date: date, Mode: mode, AvgPerc: perc, AvgDur: dur, AvgTries: tries, Count: count, Wpm: wpm, Farnsworth: farns}
}

func (s Summary) String() string {
	str := fmt.Sprintf("- Date: %s\tMode: %s\tAvg %% Correct: %.2f%%\t Avg Dur: %s\tAvg Tries: %.2f\t WPM: %d\tFarnsworth: %d", s.Date, s.Mode, s.AvgPerc * 100, s.AvgDur.Round(time.Second / 100), s.AvgTries, s.Wpm, s.Farnsworth)
	return str
}

func New() *UserStats {
	t := time.Now()
	u := new(UserStats)
	if cu, err := user.Current(); err == nil {
		u.Username = cu.Username
	}
	u.Summaries = make([]Summary, 0)
	u.Version = StatVersion
	u.Created = t
	u.Updated = t
	return u
}

func (u *UserStats) Add(s Summary) {
	u.Summaries = append(u.Summaries, s)
	u.Updated = time.Now()
	return
}

func Load(s ...string) (*UserStats, error) {
	var saveFile string
	if len(s) > 0 && s[0] != "" {
		saveFile = s[0]
	} else {
		baseDir := defaultStatDir()
		saveFile = filepath.Join(baseDir, "user-stats")
	}

	fp, err := os.Open(saveFile)
	if err != nil {
		// If the file doesn't exist, it just means the data's never
		// been saved. Create a new UserStats object and send it back.
		if os.IsNotExist(err) {
			uNew := New()
			uNew.saveFilePath = saveFile
			return uNew, nil
		}
		return nil, err
	}
	u := new(UserStats)
	dec := gob.NewDecoder(fp)
	err = dec.Decode(&u)
	if err != nil {
		fp.Close()
		return nil, err
	}
	u.saveFilePath = saveFile
	return u, fp.Close()
}

func (u *UserStats) Save(s ...string) error {
	if len(s) > 0 && s[0] != "" {
		u.saveFilePath = s[0]
	}
	if err := os.Mkdir(filepath.Dir(u.saveFilePath), 0755); err != nil && !os.IsExist(err) {
		return err
	}

	fp, err := os.CreateTemp(filepath.Dir(u.saveFilePath), "user-stats")
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(fp)
	err = enc.Encode(u)
	if err != nil {
		fp.Close()
		return err
	}

	if err = fp.Close(); err != nil {
		return err
	}

	return os.Rename(fp.Name(), u.saveFilePath)
}

func defaultStatDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "morseudar")
	}
	return filepath.Join(os.Getenv("HOME"), ".morseudar")
}
