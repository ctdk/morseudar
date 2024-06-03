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

// Statement to generate the word list

//go:generate go run gen-wordlists.go top google-10000-english-usa-no-swears.txt topwords.go
//go:generate go run gen-wordlists.go qcode q-code.txt qcodes.go
//go:generate go run gen-wordlists.go char chars.txt chars.go
//go:generate go run gen-wordlists.go koch koch.txt koch.go

package wordlists 

import (
	"github.com/ctdk/morseudar/internal/morserrors"
	"github.com/ctdk/morseudar/internal/morsestrings"
	"math/rand"
	"strings"
)

type Wordlist struct {
	wl []morsestrings.MorseString
	numLines int
	pos int
	rand *rand.Rand
}

// These words are ganked from `google-10000-english-usa-no-swears.txt`
// in https://github.com/first20hours/google-10000-english. Using the no-swear
// version because that seems best somehow.  

func GetTopWords(num int, src rand.Source) *Wordlist {
	allCnt := len(topWords)
	if num > allCnt || num == 0 {
		num = allCnt
	}
	wl := topWords[:num]

	return makeWordlistFromSlice(wl, src)
}

func GetQCodes(noQuestions bool, src rand.Source) *Wordlist {
	qLen := len(qcodeWords)
	if noQuestions {
		qLen /= 2
	}
	ql := qcodeWords[:qLen]

	return makeWordlistFromSlice(ql, src)
}

// do alphabet only, num only, etc. versions later
func GetChars(src rand.Source) *Wordlist {
	return makeWordlistFromSlice(charWords, src)
}

// MakeWordlist accepts either a string of text or a slice of strings and
// converts it into a Wordlist. NB: This is not what you want if you're looking
// to import a chunk of text to send line-per-line.
//
// TODO: Make mo' betta
func MakeWordlist(text interface{}, src rand.Source) *Wordlist {
	var words []string

	// Temporarily take the easy way out with invalid input
	switch t := text.(type) {
	case string:
		// Hmm!
		words = strings.Fields(t)
	case []string:
		for _, line := range t {
			// I suspect it's not significantly cheaper to test if
			// the line has any spaces in it before splitting it.
			lSplit := strings.Fields(line)
			words = append(words, lSplit...)
		}
	}

	return makeWordlistFromSlice(words, src)
}

func makeWordlistFromSlice(lines []string, src rand.Source) *Wordlist {
	wl := make([]morsestrings.MorseString, len(lines))
	for i, v := range lines {
		wl[i] = morsestrings.StringToMorse(v)
	}
	w := new(Wordlist)
	w.wl = wl
	w.rand = rand.New(src)
	w.numLines = len(wl)
	w.pos = 0

	return w
}

func (w *Wordlist) GetNextLine() (morsestrings.MorseString, error) {
	if w.pos >= w.numLines {
		return nil, morserrors.EOF
	}
	wret := w.wl[w.pos]
	w.pos++
	return wret, nil
}

func (w *Wordlist) GetAllLines() ([]morsestrings.MorseString, error) {
	return w.wl, nil
}

func (w *Wordlist) NumLines() int {
	return w.numLines
}

func (w *Wordlist) RandomLine() (morsestrings.MorseString, error) {
	n := w.rand.Intn(w.numLines)
	return w.wl[n], nil
}

func (w *Wordlist) Reset() error {
	w.pos = 0
	return nil
}

func (w *Wordlist) Seek(n int) error {
	absn := n
	if absn < 0 {
		absn = -absn
	}

	if absn > w.numLines {
		return morserrors.OutOfRange
	}
	if n < 0 {
		n = w.numLines + n
	}
	w.pos = n
	return nil
}
