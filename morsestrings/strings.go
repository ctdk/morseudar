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

package morsestrings

import (
	"strings"
)

type MorseChar string
type MorseWord struct {
	word []MorseChar
	text string
	prosign bool
}

type MorseString []*MorseWord

var Alphabet = map[rune]MorseChar {
	'a': ".-",
	'b': "-...",
	'c': "-.-.",
	'd': "-..",
	'e': ".",
	'f': "..-.",
	'g': "--.",
	'h': "....",
	'i': "..",
	'j': ".---",
	'k': "-.-",
	'l': ".-..",
	'm': "--",
	'n': "-.",
	'o': "---",
	'p': ".--.",
	'q': "--.-",
	'r': ".-.",
	's': "...",
	't': "-",
	'u': "..-",
	'v': "...-",
	'w': ".--",
	'x': "-..-",
	'y': "-.--",
	'z': "--..",
	'1': ".----",
	'2': "..---",
	'3': "...--",
	'4': "....-",
	'5': ".....",
	'6': "-....",
	'7': "--...",
	'8': "---..",
	'9': "----.",
	'0': "-----",
	'.': ".-.-.-",
	',': "--..--",
	'?': "..--..",
	'=': "-...-",
	'/': "-..-.",
	':': "---...",
	'(': "-.--.",
	')': "-.--.-",
	'+': ".-.-.",
	'-': "-....-",
	'&': ".-...",
	'"': ".-..-.",
	'\'': ".----.",
	'@': ".--.-.",
}

const wordJoin = " / "

// hopefully this isn't overdoing keeping things private

func (mw *MorseWord) IsProsign() bool {
	return mw.prosign
}

func (mw *MorseWord) Chars() []MorseChar {
	return mw.word
}

func (mw *MorseWord) Len() int {
	return len(mw.word)
}

func (mw *MorseWord) String() string {
	return mw.text
}

func StringToMorse(str string) MorseString {
	str = strings.ToLower(str)
	toConv := strings.Split(str, " ")
	m := make([]*MorseWord, len(toConv))

	for i, s := range toConv {
		mw := new(MorseWord)

		// I'm lazy, sue me.
		if strings.HasPrefix(s, "~") && strings.HasSuffix(s, "~") && len(s) > 2 {
			mw.prosign = true
			s = s[1:len(s)-1]
		}
		mw.word = make([]MorseChar, 0, len(s))

		for _, c := range s {
			if mc, ok := Alphabet[c]; ok {
				// if the character isn't in the list, skip it
				mw.word = append(mw.word, mc)
			}
		}

		// Having this available will be very useful, it turns out.
		mw.text = s

		m[i] = mw
	}

	return MorseString(m)
}

// DotDashString spits out the encoded morse characters as dots and dashes
func (ms MorseString) DotDashString() string {
	str := make([]string, len(ms))

	for i, m := range ms {
		var joiner string
		if !m.prosign {
			joiner = " "
		} else {
			joiner = ""
		}

		wordAssemble := make([]string, len(m.word))
		for j, c := range m.word {
			wordAssemble[j] = string(c)
		}

		w := strings.Join(wordAssemble, joiner)
		str[i] = w
	}

	dotdash := strings.Join(str, wordJoin)

	return dotdash
}

// RawString returns the raw string that went into making this morse string.
// Mostly useful for testing.
func (ms MorseString) RawString() string {
	// short-circuit the common case where there's only one word in the
	// string.
	if len(ms) == 1 {
		return ms[0].text
	}
	strs := make([]string, len(ms))
	for i, v := range ms {
		strs[i] = v.text
	}

	return strings.Join(strs, " ")
}
