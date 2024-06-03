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
	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"time"
)

const (
	levReplaceCost = 1
	levInsertCost = 1
	levDeleteCost = 1
)

type Answer struct {
	Original string
	Response string
	Percentage float64
	Took time.Duration
	Tries int
}

type AnswerBatch []Answer

// avoiding constantly recreating the Levenshtein object
type Comparator struct {
	lev *metrics.Levenshtein
}

func New() *Comparator {
	// try Levenshein
	lev := metrics.NewLevenshtein()
	lev.CaseSensitive = false

	// fiddle with these as needed
	lev.ReplaceCost = levReplaceCost
	lev.InsertCost = levInsertCost
	lev.DeleteCost = levDeleteCost

	return &Comparator{lev: lev}
}

func CompareStrings(orig string, resp string) float64 {
	// try Levenshein
	lev := metrics.NewLevenshtein()
	lev.CaseSensitive = false

	// fiddle with these as needed
	lev.ReplaceCost = levReplaceCost
	lev.InsertCost = levInsertCost
	lev.DeleteCost = levDeleteCost

	sim := strutil.Similarity(orig, resp, lev)
	return sim
}

func (c *Comparator) Compare(orig string, resp string, start time.Time, tries int) Answer {
	sim := strutil.Similarity(orig, resp, c.lev)
	took := time.Since(start)
	ans := Answer{ Original: orig,
		Response: resp,
		Percentage: sim,
		Took: took,
		Tries: tries,
	}

	return ans
}

func (ab AnswerBatch) Averages() (float64, time.Duration, float64) {
	n := len(ab)
	if n == 0 {
		return 0, 0, 0
	}

	var percTotal float64
	var durTotal time.Duration
	var triesTotal int
	for _, ans := range ab {
		percTotal += ans.Percentage
		durTotal += ans.Took
		triesTotal += ans.Tries
	}

	avgPerc := percTotal / float64(n)
	avgDur := durTotal / time.Duration(n)
	avgTries := float64(triesTotal) / float64(n)

	return avgPerc, avgDur, avgTries
}
