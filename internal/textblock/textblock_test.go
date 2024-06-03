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

package textblock

import (
	"math/rand"
	"os"
	"testing"
	"strings"
)

const randSeed = 12345
var src = rand.NewSource(randSeed)

func TestLoadNewText(t *testing.T) {
	msg := "hi"
	msg1 := "yo yo"

	expectedMsg := ".... .."
	expectedNewMsg := "-.-- --- / -.-- ---"

	tb := NewTextblock(src)

	f, err := os.CreateTemp("", "msg-test")
	if err != nil {
		t.Errorf("error creating temp file: %v\n", err)
	}
	if _, err := f.Write([]byte(msg)); err != nil {
		t.Errorf("error writing file: %v\n", err)
	}
	defer os.Remove(f.Name())

	f1, err := os.CreateTemp("", "msg1-test")
	if err != nil {
		t.Errorf("error creating temp file: %v\n", err)
	}
	if _, err := f1.Write([]byte(msg1)); err != nil {
		t.Errorf("error writing file: %v\n", err)
	}
	defer os.Remove(f1.Name())

	if lerr := tb.LoadFile(f.Name()); err != nil {
		t.Errorf("error loading msg file: %v\n", lerr)
	}

	// Whadda we got?
	mFirst, _ := tb.GetNextLine()
	if expectedMsg != mFirst.DotDashString() {
		t.Errorf("expectedMsg test failed: wanted '%s', got '%s'", expectedMsg, mFirst.DotDashString())
	}

	if lerr := tb.LoadFile(f1.Name()); err != nil {
		t.Errorf("error loading msg1 file: %v\n", lerr)
	}
	m1First, _ := tb.GetNextLine()
	if expectedNewMsg != m1First.DotDashString() {
		t.Errorf("expectedMsg test failed: wanted '%s', got '%s'", expectedMsg, m1First.DotDashString())
	}
}

func TestRandomLines(t *testing.T) {
	rt := []string{"hello", "howdy", "dawg"}

	tb := NewTextblock(src)

	f, err := os.CreateTemp("", "rand-line-test")
	if err != nil {
		t.Errorf("error creating temp file: %v\n", err)
	}
	if _, err := f.Write([]byte(strings.Join(rt, "\n"))); err != nil {
		t.Errorf("error writing file: %v\n", err)
	}

	defer os.Remove(f.Name())

	if lerr := tb.LoadFile(f.Name()); err != nil {
		t.Errorf("error loading rand test file: %v\n", lerr)
	}

	expectedStrIdx := []int{2, 1, 1}
	for i, idx := range expectedStrIdx {
		ms, err := tb.RandomLine()
		if err != nil {
			t.Errorf("err getting random line: %s", err.Error())
		}
		if rt[idx] != ms[0].String() {
			t.Errorf("Round #%d of getting random lines of MorseStrings returned a word did not match the expected value - expected '%s', got '%s'", i, rt[idx], ms[0].String())
		}
	}
}
