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

// Statement to generate the word list

//go:generate sh -c "./gen-wordlists.pl top < ./google-10000-english-usa-no-swears.txt > topwords.go && go fmt topwords.go"
//go:generate sh -c "./gen-wordlists.pl qcode < ./q-code.txt > qcodes.go && go fmt qcodes.go"

package wordlists 

import (
	"github.com/ctdk/morse-copying/morse"
	"strings"
)

type Wordlist []morse.MorseString

// These words are ganked from `google-10000-english-usa-no-swears.txt`
// in https://github.com/first20hours/google-10000-english. Using the no-swear
// version because that seems best somehow.  

func GetTopWords(num int) Wordlist {
	allCnt := len(topWords)
	if num > allCnt || num == 0 {
		num = allCnt
	}
	wl := topWords[:num]

	return makeWordlistFromSlice(wl)
}

func GetQCodes(noQuestions bool) Wordlist {
	qLen := len(qcodeWords)
	if noQuestions {
		qLen /= 2
	}
	ql := qcodeWords[:qLen]

	return makeWordlistFromSlice(ql)
}

// MakewordList accepts either a string of text or a slice of strings and
// converts it into a Wordlist. NB: This is not what you want if you're looking
// to import a chunk of text to send line-per-line.
//
// TODO: Mke mo' betta
func MakeWordlist(text interface{}) Wordlist {
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

	return makeWordlistFromSlice(words)
}

func makeWordlistFromSlice(lines []string) Wordlist {
	wl := make([]morse.MorseString, len(lines))
	for i, v := range lines {
		wl[i] = morse.StringToMorse(v)
	}
	return Wordlist(wl)
}
