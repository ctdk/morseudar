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

package codegroups

import (
	"fmt"
	"github.com/ctdk/morseudar/internal/morserrors"
	"github.com/ctdk/morseudar/internal/morsestrings"
	"math/rand"
	"strings"
	"unsafe"
)

/* Borrowing this from what might be the single most detailed Stack Overflow
 * answer I've ever come across. Yowza!
 * https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
 */
const letterBytes = "abcdefghijklmnopqrstuvwxyz"
const alnumBytes = "abcdefghijklmnopqrstuvwxyz01234567890"

const (
	letterIdxBits = 5
	letterIdxMask = 1 << letterIdxBits - 1
	letterIdxMax = 63 / letterIdxBits
	alnumIdxBits = 6
	alnumIdxMask = 1 << alnumIdxBits - 1
	alnumIdxMax = 63 / alnumIdxBits
)

type CodeGroupType uint8
const (
	Alpha CodeGroupType = iota
	Alnum
	Num
)

const (
	CodeGroupLen = 5
	CodeGroupPer = 5
)

type Codegroup struct {
	src rand.Source
	groupType CodeGroupType
	codeGroupLen int
	codeGroupPer int
}

func NewCodegroup(src rand.Source, groupType CodeGroupType, codeGroupLen int, codeGroupPer int) *Codegroup {
	if codeGroupLen == 0 {
		codeGroupLen = CodeGroupLen
	}

	if codeGroupPer == 0 {
		codeGroupPer = CodeGroupPer
	}

	return &Codegroup{src: src, groupType: groupType, codeGroupLen: codeGroupLen, codeGroupPer: codeGroupPer}
}

func (cg *Codegroup) RandLetterCodegroupNum(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax
	// characters!
	for i, cache, remain := n-1, cg.src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = cg.src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func (cg *Codegroup) RandLetterCodegroupLineNum(n int, cgLen int) string {
	str := make([]string, n)
	for i := 0; i < n; i++ {
		str[i] = cg.RandLetterCodegroupNum(cgLen)
	}

	final := strings.Join(str, " ")
	return final
}

func (cg *Codegroup) RandAlnumCodegroupNum(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for alnumIdxMax
	// characters!
	for i, cache, remain := n-1, cg.src.Int63(), alnumIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = cg.src.Int63(), alnumIdxMax
		}
		if idx := int(cache & alnumIdxMask); idx < len(alnumBytes) {
			b[i] = alnumBytes[idx]
			i--
		}
		cache >>= alnumIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func (cg *Codegroup) RandAlnumCodegroupLineNum(n int, cgLen int) string {
	str := make([]string, n)
	for i := 0; i < n; i++ {
		str[i] = cg.RandAlnumCodegroupNum(cgLen)
	}

	final := strings.Join(str, " ")
	return final
}

func (cg *Codegroup) RandNumberCodegroupNum(n int) string {
	// really need some length checks...
	q := cg.src.Int63()
	qstr := fmt.Sprintf("%d", q)
	qlen := len(qstr)
	return qstr[qlen-n:qlen]
}

func (cg *Codegroup) RandNumberCodegroupLineNum(n int, cgLen int) string {
	str := make([]string, n)
	for i := 0; i < n; i++ {
		str[i] = cg.RandNumberCodegroupNum(cgLen)
	}
	final := strings.Join(str, " ")
	return final
}

// morselist interface functions

func (cg *Codegroup) NumLines() int {
	return 0
}

func (cg *Codegroup) GetAllLines() ([]morsestrings.MorseString, error) {
	return nil, morserrors.NotApplicable
}

func (cg *Codegroup) Reset() error {
	return morserrors.NotApplicable
}

func (cg *Codegroup) Seek(n int) error {
	return morserrors.NotApplicable
}

func (cg *Codegroup) RandomLine() (morsestrings.MorseString, error) {
	var str string

	switch cg.groupType {
	case Alpha:
		str = cg.RandLetterCodegroupLineNum(cg.codeGroupPer, cg.codeGroupLen)
	case Alnum:
		str = cg.RandAlnumCodegroupLineNum(cg.codeGroupPer, cg.codeGroupLen)
	case Num:
		str = cg.RandNumberCodegroupLineNum(cg.codeGroupPer, cg.codeGroupLen)
	}
	mStr := morsestrings.StringToMorse(str)
	return mStr, nil
}

func (cg *Codegroup) GetNextLine() (morsestrings.MorseString, error) {
	return cg.RandomLine()
}
