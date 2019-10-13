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

//go:generate sh -c "./gen-topwords.pl < ./google-10000-english-usa-no-swears.txt > wordlist.go && go fmt wordlist.go"

package topwords

import (
	"github.com/ctdk/morse-copying/morse"
)

// These words are ganked from `google-10000-english-usa-no-swears.txt`
// in https://github.com/first20hours/google-10000-english. Using the no-swear
// version because that seems best somehow.  

func GetTopWords(num int) []morse.MorseString {
	allCnt := len(topWds)
	if num > allCnt || num == 0 {
		num = allCnt
	}
	wl := topWds[:num]

	tw := make([]morse.MorseString, num)
	for i, w := range wl {
		tw[i] = morse.StringToMorse(w)
	}

	return tw
}
