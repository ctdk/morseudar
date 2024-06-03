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

import (
	"github.com/ctdk/morseudar/internal/audio"
	"github.com/ctdk/morseudar/internal/morsestrings"
	"testing"
	"strings"
)

const randSeed = 12345

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
	var freq float64 = 660
	wpm := 100 // we want this to actually finish someday 
	msg := "yo dawg ~sk~"

	msg2 := "~cq~ ~cq~ h3llo"

	ma, err := audio.NewMorseAudio(freq, wpm, 0)
	if err != nil {
		t.Errorf("error creating MorseAudio: %s", err.Error())
	}

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
	var freq float64 = 660
	wpm := 100
	msg := "hi"

	expectedMsg := ".... .."

	m, err := New(Default, wpm, 0, freq, false, randSeed)
	if err != nil {
		t.Errorf("error creating morse object: %s", err.Error())
	}

	// TODO: Test the actual object properties or something
}
