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
	"github.com/ctdk/morse-copying/morsestrings"
	"testing"
	"strings"
)

const randSeed int64 = 12345

func TestStringToMorse(t *testing.T) {
	s := "foo bar"
	ms := morsestrings.StringToMorse(s)
	expected := "..-. --- --- / -... .- .-."

	if len(ms) != 2 {
		t.Errorf("MorseString object should have had len == 2, but was actually %d", len(ms))
	}

	dd := ms.DotDashString()

	if dd != expected {
		t.Errorf("Unexpected string representation of '%s': got '%s' instead", s, dd)
	}
}

func TestMorseSendMessage(t *testing.T) {
	var freq uint = 660
	var amplitudeLevel uint8 = 220
	wpm := 100 // we want this to actually finish someday 
	msg := "yo dawg ~sk~"

	msg2 := "~cq~ ~cq~ h3llo"

	ma, err := NewMorseAudio(freq, amplitudeLevel, wpm)
	if err != nil {
		t.Errorf("error creating MorseAudio: %s", err.Error())
	}
	defer ma.Destroy()

	morseMsg := morsestrings.StringToMorse(msg)
	err = ma.SendMessage(morseMsg)
	if err != nil {
		t.Errorf("error sending message: %s", err.Error())
	}

	morseMsg2 := morsestrings.StringToMorse(msg2)
	err = ma.SendMessage(morseMsg2)
	if err != nil {
		t.Errorf("error sending message2: %s", err.Error())
	}
}

func TestMorseObject(t *testing.T) {
	var freq uint = 660
	var amplitudeLevel uint8 = 220
	wpm := 100
	msg := "hi"

	expectedMsg := ".... .."

	m, err := New(Default, wpm, freq, amplitudeLevel, msg, randSeed)
	if err != nil {
		t.Errorf("error creating morse object: %s", err.Error())
	}

	// Whadda we got?
	if expectedMsg != m.Lines[0].DotDashString() {
		t.Errorf("expectedMsg test failed: wanted '%s', got '%s'", expectedMsg, m.Lines[0].DotDashString())
	}
	if err = m.SendLineNum(0); err != nil {
		t.Errorf("sending a mesage with the morse object failed: %s", err.Error())
	}
}

func TestLoadNewText(t *testing.T) {
	var freq uint = 660
	var amplitudeLevel uint8 = 220
	wpm := 100
	msg := "hi"
	newMsg := "yo yo"

	expectedMsg := ".... .."
	expectedNewMsg := "-.-- --- / -.-- ---"

	m, err := New(Default, wpm, freq, amplitudeLevel, msg, randSeed)
	if err != nil {
		t.Errorf("error creating morse object: %s", err.Error())
	}

	// Whadda we got?
	if expectedMsg != m.Lines[0].DotDashString() {
		t.Errorf("expectedMsg test failed: wanted '%s', got '%s'", expectedMsg, m.Lines[0].DotDashString())
	}

	err = m.LoadText(newMsg)
	if expectedNewMsg != m.Lines[0].DotDashString() {
		t.Errorf("expectedMsg test failed: wanted '%s', got '%s'", expectedMsg, m.Lines[0].DotDashString())
	}
}

func TestRandomLines(t *testing.T) {
	var freq uint = 660
	var amplitudeLevel uint8 = 220
	wpm := 100
	rt := []string{"hello", "howdy", "dawg"}

	m, err := New(Default, wpm, freq, amplitudeLevel, strings.Join(rt, "\n"), randSeed)
	if err != nil {
		t.Errorf("error creating morse object: %s", err.Error())
	}

	// call random lines, see what we get back.
	expectedLineNums := []int{2, 1, 1, 1, 2, 2, 0, 0, 1 ,0}
	for i, line := range expectedLineNums {
		ln, err := m.RandomLineNum()
		if err != nil {
			t.Errorf("somehow got an error getting a random line number: %s", err.Error())
		}
		if line != ln {
			t.Errorf("RandomLineNum call #%d failed - expected line %d, got %d", i, line, ln)
		}
	}

	// reset the seed!
	m.SetSeed(randSeed)

	expectedStrIdx := []int{2, 1, 1}
	for i, idx := range expectedStrIdx {
		ms, err := m.RandomLine()
		if err != nil {
			t.Errorf("err getting random line: %s", err.Error())
		}
		if rt[idx] != ms[0].String() {
			t.Errorf("Round #%d of getting random lines of MorseStrings returned a word did not match the expected value - expected '%s', got '%s'", i, rt[idx], ms[0].String())
		}
	}
}
