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
	"errors"
	"fmt"
	"math/rand"
	"strings"
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
	WPM   int
	Frequency uint
	Amplitude uint8
	Lines []MorseString
	Mode MorseMode
	audio *MorseAudio
	numLines int32
}

func New(mode MorseMode, wpm int, freq uint, amplitude uint8, text string) (*Morse, error) {
	m := new(Morse)
	m.Mode = mode

	if wpm == 0 {
		m.WPM = defaultWPM
	} else {
		m.WPM = wpm
	}

	if freq == 0 {
		m.Frequency = defaultFrequency
	} else {
		m.Frequency = freq
	}

	if amplitude == 0 {
		m.Amplitude = defaultAmplitudeLevel
	} else {
		m.Amplitude = amplitude
	}

	if text != "" {
		if tErr := m.LoadText(text); tErr != nil {
			return nil, tErr
		}
	}

	ma, err := NewMorseAudio(m.Frequency, m.Amplitude, m.WPM)
	if err != nil {
		return nil, err
	}
	m.audio = ma

	return m, nil
}

func (m *Morse) LoadText(text string) error {
	if len(text) == 0 {
		return errors.New("Can not load empty text.")
	}
	textLines := strings.Split(text, "\n")
	m.Lines = make([]MorseString, 0, len(textLines))
	for _, l := range textLines {
		l = strings.TrimSpace(l)
		if len(l) > 0 {
			ml := StringToMorse(l)
			m.Lines = append(m.Lines, ml)
		}
	}

	m.numLines = len(m.Lines)
	return nil
}

func (m *Morse) RandomLine() error {
	if m.numLines == 0 {
		return errors.New("no text has been loaded!")
	}
	i := rand.Int31n(m.Numblines)
	return m.Send(i)
}

func (m *Morse) Send(lineNum int) error {
	switch {
	case lineNum < 0:
		return errors.New("line number cannot be negative")
	case lineNum >= m.numLines:
		return fmt.Errorf("line number '%d' is out of range", lineNum)
	}

	if err := m.audio.SendMessage(m.Lines[lineNum]); err != nil {
		return err
	}
	return nil
}
