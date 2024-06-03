/*
 * Copyright (c) 2024, Jeremy Bingham (<jeremy@goiardi.gl>)
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
	"errors"
	"github.com/ctdk/morseudar/internal/morserrors"
	"github.com/ctdk/morseudar/internal/morsestrings"
	"io"
	"math/rand"
	"os"
	"strings"
)

type Textblock struct {
	lines []morsestrings.MorseString
	numLines int
	pos int
	randNum *rand.Rand
}

func NewTextblock(src rand.Source) *Textblock{
	tb := new(Textblock)
	tb.randNum = rand.New(src)
	return tb
}

func (tb *Textblock) LoadFile(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	txt, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return tb.processTextFile(string(txt))
}

func (tb *Textblock) processTextFile(text string) error {
	if len(text) == 0 {
		return errors.New("Can not load empty text.")
	}
	textLines := strings.Split(text, "\n")
	tb.lines = make([]morsestrings.MorseString, 0, len(textLines))
	for _, l := range textLines {
		l = strings.TrimSpace(l)
		if len(l) > 0 {
			ml := morsestrings.StringToMorse(l)
			tb.lines = append(tb.lines, ml)
		}
	}

	tb.numLines = len(tb.lines)
	tb.pos = 0 // make sure pos is set back to 0
	return nil
}

func (tb *Textblock) GetNextLine() (morsestrings.MorseString, error) {
	if tb.numLines == 0 {
		return nil, morserrors.NoText
	}
	if tb.pos >= tb.numLines {
		return nil, morserrors.EOF
	}
	tbret := tb.lines[tb.pos]
	tb.pos++
	return tbret, nil
}

func (tb *Textblock) GetAllLines() ([]morsestrings.MorseString, error) {
	if tb.numLines == 0 {
		return nil, morserrors.NoText
	}
	return tb.lines, nil
}

func (tb *Textblock) NumLines() int {
	return tb.numLines
}

func (tb *Textblock) RandomLine() (morsestrings.MorseString, error) {
	if tb.numLines == 0 {
		return nil, morserrors.NoText
	}
	n := tb.randNum.Intn(tb.numLines)
	tb.pos = n
	return tb.lines[n], nil
}

func (tb *Textblock) Reset() error {
	if tb.numLines == 0 {
		return morserrors.NoText
	}
	tb.pos = 0
	return nil
}

func (tb *Textblock) Seek(n int) error {
	absn := n
	if absn < 0 {
		absn = -absn
	}

	if absn > tb.numLines {
		return morserrors.OutOfRange
	}
	if n < 0 {
		n = tb.numLines + n
	}
	tb.pos = n
	return nil
}
