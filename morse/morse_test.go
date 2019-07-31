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
	"testing"
)

func TestStringToMorse(t *testing.T) {
	s := "foo bar"
	ms := StringToMorse(s)
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
	wpm := 15
	msg := "yo dawg, what's up? ~sk~"

	msg2 := "~cq~ ~cq~ ki7doz de us. ~cq~ ~cq~ ki7doz de us."

	ma, err := NewMorseAudio(freq, amplitudeLevel, wpm)
	if err != nil {
		t.Errorf("error creating MorseAudio: %s", err.Error())
	}
	defer ma.Destroy()

	morseMsg := StringToMorse(msg)
	err = ma.SendMessage(morseMsg)
	if err != nil {
		t.Errorf("error sending message: %s", err.Error())
	}

	morseMsg2 := StringToMorse(msg2)
	err = ma.SendMessage(morseMsg2)
	if err != nil {
		t.Errorf("error sending message2: %s", err.Error())
	}
}
