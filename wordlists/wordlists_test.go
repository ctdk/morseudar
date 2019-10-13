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
	"testing"
)

func TestGetTopXWords(t *testing.T) {
	topLen := 3
	topwords := GetTopWords(topLen)
	if len(topwords) != topLen {
		t.Errorf("The top %d words should have %d elements, but got %d back instead.", topLen, topLen, len(topwords))
	}
	expectedWords := []string{"the", "of", "and"}
	for i := 0; i < len(topwords); i++ {
		if topwords[i].RawString() != expectedWords[i] {
			t.Errorf("Top word #%d should have been '%s', but was '%s'.", i, topwords[i].RawString(), expectedWords[i])
		}
	}
}
