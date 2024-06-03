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

package wordlists

import (
	"math/rand"
	"testing"
)

const randSeed = 12345
var src = rand.NewSource(randSeed)
var randNum = rand.New(src)

// NB: these are going to be horribly broken for the time being

func TestGetTopXWords(t *testing.T) {
	topLen := 3
	topwords := GetTopWords(topLen, randNum)
	if topwords.NumLines() != topLen {
		t.Errorf("The top %d words should have %d elements, but got %d back instead.", topLen, topLen, topwords.NumLines())
	}
	expectedWords := []string{"the", "of", "and"}
	for i := 0; i < topwords.NumLines(); i++ {
		k, _ := topwords.GetNextLine()
		if k.RawString() != expectedWords[i] {
			t.Errorf("Top word #%d should have been '%s', but was '%s'.", i, k.RawString(), expectedWords[i])
		}
	}
}

func TestQCodes(t *testing.T) {
	qcodeLen := 98 // we happen to know this
	allQ := GetQCodes(false, randNum)
	if allQ.NumLines() != qcodeLen {
		t.Errorf("Somehow the number of Q codes was not %d, got %d.", qcodeLen, allQ.NumLines())
	}
	q1, _ := allQ.GetNextLine()
	if q1.RawString() != "qni" {
		t.Errorf("The first q code should have been 'qni', but got '%s'.", q1.RawString())
	}
	allQ.Seek(-1)
	qLast, _ := allQ.GetNextLine()
	if qLast.RawString() != "quf?" {
		t.Errorf("The last q code should have been 'quf?', but got '%s'.", qLast.RawString())
	}

	noQLen := qcodeLen / 2
	noQQ := GetQCodes(true, randNum)
	if noQQ.NumLines() != noQLen {
		t.Errorf("Somehow the number of Q codes excluding questions was not %d, got %d.", noQLen, noQQ.NumLines())
	}
	nq1, _ := noQQ.GetNextLine()
	if nq1.RawString() != "qni" {
		t.Errorf("The first q code excluding questions should have been 'qni', but got '%s'.", nq1.RawString())
	}
	noQQ.Seek(-1)
	nqLast, _ := noQQ.GetNextLine()
	if nqLast.RawString() != "quf" {
		t.Errorf("The last q code excluding questions should have been 'quf', but got '%s'.", nqLast.RawString())
	}
}
