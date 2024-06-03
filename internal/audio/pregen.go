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

package audio

import (
	"errors"
	"github.com/gopxl/beep"
	"math"
	"time"
)

type preCalcSine struct {
	pos int
	samples [][2]float64
	freq float64 
	sr beep.SampleRate
	dur time.Duration
	err error
}

const decayLenPercentage = 5
const decayFactor float64 = 0.90

func PreCalcSine(sr beep.SampleRate, freq float64, dur time.Duration) (beep.StreamSeekCloser, beep.Format, error) {
	dt := freq / float64(sr)
	if dt > 1.0/2.0 {
		return nil, beep.Format{}, errors.New("sample rate must be at least two times greater than the frequency")
	}

	// how many samples, then?
	sampleLen := sr.N(dur)
	samples := make([][2]float64, sampleLen)

	var t float64 = 0

	decayLen := sampleLen * decayLenPercentage / 100
	decayStep := 0

	for i := 0; i < sampleLen; i++ {
		var n [2]float64
		v := math.Sin(t * 2.0 * math.Pi)

		// damp down the end a bit
		if i > sampleLen - decayLen {
			v *= math.Pow(decayFactor, float64(decayStep))
			decayStep++
		}

		n[0] = v
		n[1] = v

		samples[i] = n
		_, t = math.Modf(t + dt)
	}

	pc := new(preCalcSine)
	pc.pos = 0
	pc.samples = samples
	pc.freq = freq
	pc.sr = sr
	pc.dur = dur
	f := beep.Format{
		SampleRate: sr,
		NumChannels: 2,
		Precision: 4, // hopefully enough?
	}
	return pc, f, nil
}

func (pc *preCalcSine) Err() error {
	return pc.err
}

func (pc *preCalcSine) Stream(samples [][2]float64) (n int, ok bool) {
	sampleLen := len(pc.samples)

	for i := range samples {
		if pc.pos >= sampleLen {
			return 0, false
		}

		samples[i] = pc.samples[pc.pos]
		pc.pos++
		n++
	}
	ok = true

	return n, ok
}

func (pc *preCalcSine) Len() int {
	return len(pc.samples)
}

func (pc *preCalcSine) Position() int {
	return pc.pos
}

func (pc *preCalcSine) Seek(p int) error {
	if p > len(pc.samples) {
		return errors.New("cannot seek position past end of stream")
	}
	pc.pos = p
	return nil
}

func (pc *preCalcSine) Close() error {
	// might not be strictly useful
	// but set it back to 0 I guess. Maybe it should free things up though.
	pc.pos = 0
	return nil
}
