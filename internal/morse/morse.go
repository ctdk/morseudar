/*
 * Copyright (c) 2019-2024, Jeremy Bingham (<jeremy@goiardi.gl>)
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

//go:generate stringer -type=MorseMode

import (
	"github.com/ctdk/morseudar/internal/audio"
	"github.com/ctdk/morseudar/internal/morsestrings"
	"math/rand"
	"time"
)

type MorseMode uint8

const (
	TextFile MorseMode = iota // Just play the supplied text
	CodeGroup // generate random 5 groups of 5 alphabetic strings
	CodeAlnum // ibid, but alphanumeric
	CodeNum // ibid, but numbers
	TopWords // play a word from the top words list
	Qcode // play a q code from the Q code list
	MorseChar // play a single character from the character list
	Koch // play a character from the Koch list. NB: will need work behind
	     // the scenes to work correctly, since it's not sequential but
	     // plays two chars until they're copied with 90% accuracy, then
	     // adding another.
)

const (
	defaultFrequency = 700
	defaultWPM = 10
)

type Morse struct {
	WPM int
	Farnsworth int
	Frequency float64
	Mode MorseMode
	Sequential bool
	EntireBlock bool
	TestingMaterial MorseList
	audio *audio.MorseAudio
	src rand.Source
	linesTested int
	percentage float64
}

type MorseList interface {
	NumLines() int
	RandomLine() (morsestrings.MorseString, error)
	GetNextLine() (morsestrings.MorseString, error)
	GetAllLines() ([]morsestrings.MorseString, error)
	Reset() error
	Seek(int) error
}

func New(mode MorseMode, wpm int, farn int, freq float64, seq bool, entire bool, randSeed int64) (*Morse, error) {
	m := new(Morse)
	m.Mode = mode

	if wpm == 0 {
		m.WPM = defaultWPM
	} else {
		m.WPM = wpm
	}

	// Ignore Farnsworth if Farnsworth is set yet somehow faster than wpm
	if farn != 0 && wpm > farn {
		m.Farnsworth = farn
	}

	if freq == 0 {
		m.Frequency = float64(defaultFrequency)
	} else {
		m.Frequency = freq
	}

	m.Sequential = seq
	m.EntireBlock = entire

	ma, err := audio.NewMorseAudio(m.Frequency, m.WPM, m.Farnsworth)
	if err != nil {
		return nil, err
	}
	m.audio = ma

	if randSeed == 0 {
		randSeed = time.Now().UnixNano()
	}
	m.src = rand.NewSource(randSeed)

	return m, nil
}

func (m *Morse) Send(ms morsestrings.MorseString) error {
	return m.audio.SendMessage(ms)
}

func (m *Morse) Src() rand.Source {
	return m.src
}

func (m *Morse) GetMorse() (morsestrings.MorseString, error) {
	if m.Sequential {
		return m.TestingMaterial.GetNextLine()
	}

	return m.TestingMaterial.RandomLine()
}
