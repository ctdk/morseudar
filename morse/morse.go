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
	//"fmt"
	//"time"
)

type MorseMode uint8

const (
	Default MorseMode = iota // Just play the supplied text
	RandomLine // play a random line from the text
	CodeGroup // generate random 5 groups of 5 alphanumeric strings
)

const CodeGroupLen = 5
const CodeGroupPer = 5

type Morse struct {
	WPM   uint8
	Frequency uint
	Lines []MorseString
	Mode MorseMode
}

func New(mode MorseMode, wpm uint8, text string) (*Morse, error) {
	// m := new(Morse)

	return nil, nil
}
