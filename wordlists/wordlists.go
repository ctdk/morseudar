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

	tw := make([]morse.MorseString, num)
	for i, w := range wl {
		tw[i] = morse.StringToMorse(w)
	}

	return Wordlist(tw)
}

func GetQCodes(noQuestions bool) Wordlist {
	qLen := len(qcodeWords)
	if noQuestions {
		qLen /= 2
	}
	ql := qcodeWords[:qLen]
	qc := make([]morse.MorseString, qLen)
	for i, q := range ql {
		qc[i] = morse.StringToMorse(q)
	}
	return Wordlist(qc)
}

func MakeWordlistFromText(text string) Wordlist {

}

func MakeWordlistFromSlice(lines []string) Wordlist {

}
