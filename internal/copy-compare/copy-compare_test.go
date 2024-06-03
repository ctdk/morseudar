/*
 * Copyright (c) 2022-2024, Jeremy Bingham (<jeremy@goiardi.gl>)
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

package compare

import (
	"math"
	"testing"
)

var ep float64 = 0.000001

func TestCopyCompare(t *testing.T) {
	orig := "foober"
	resp := "foo bar"
	z := CompareStrings(orig, resp)
	exp := 0.714286
	if !floatEq(z, exp) {
		t.Errorf("z is: %f, expected %f\n", z, exp)
	}
}

func TestCopyCompareSame(t *testing.T) {
	orig := "foober"
	resp := "foober"
	z := CompareStrings(orig, resp)
	exp := 1.0
	if !floatEq(z, exp) {
		t.Errorf("z is: %f, expected %f\n", z, exp)
	}
}

func TestCopyCompareLonger(t *testing.T) {
	orig := "fooberi sdfa foo moo goo ~sk~"
	resp := "foober"
	z := CompareStrings(orig, resp)
	exp := 0.206897
	if !floatEq(z, exp) {
		t.Errorf("z is: %f, expected %f\n", z, exp)
	}
}

func TestCopyCompareShorter(t *testing.T) {
	orig := "fee"
	resp := "fooberi sdfa foo moo goo ~sk~"
	z := CompareStrings(orig, resp)
	exp := 0.068966
	if !floatEq(z, exp) {
		t.Errorf("z is: %f, expected %f\n", z, exp)
	}
}

func TestCopyCompareAlmost(t *testing.T) {
	orig := "foobe narmi soogl"
	resp := "foobe narmi sooogr"
	z := CompareStrings(orig, resp)
	exp := 0.888889
	if !floatEq(z, exp) {
		t.Errorf("z is: %f, expected %f\n", z, exp)
	}
}

func floatEq(a float64, b float64) bool {
	if math.Abs(a - b) > ep {
		return false
	}
	return true
}
