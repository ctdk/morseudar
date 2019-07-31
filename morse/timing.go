/*
 * Copyright (c) 2019, Jeremy Bingham (<jeremy@goiardi.gl>)
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

package morse

import (
	"fmt"
	"time"
)

// timing definitions and functions

const (
	mDit = 1
	mDash = 3
	mLetterSep = 3
	mWordSep = 7
	mParis = 50 // PARIS using 50 dits for wpm calculations.
)

type DitDuration time.Duration

func CalcDitDuration(wpm int) DitDuration {
	d := time.Minute / (mParis * time.Duration(wpm))
	return DitDuration(d) // needed?
}

func (d DitDuration) String() string {
	return fmt.Sprintf("%.3f", float64(d) / float64(time.Second))
}
