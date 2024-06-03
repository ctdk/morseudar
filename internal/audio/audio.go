/*
 * Copyright (c) 2019-2024, Jeremy Bingham (<jeremy@goiardi.gl>)
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
	"fmt"
	"github.com/ctdk/morseudar/internal/morsestrings"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/generators"
	"github.com/gopxl/beep/speaker"
	"time"
)
const wordAvg = 5

var sampleRate beep.SampleRate = 48000

// the morse audio type should be able to toggle the audio stream on and off.

// TODO: also ensure cleanup on exit

// MorseAudio drops the beeps.
type MorseAudio struct {
	wpm int
	farn int
	ditDur time.Duration
	farnDitDur time.Duration
	silence beep.Streamer
	dit *beep.Buffer
	dah *beep.Buffer
	sr beep.SampleRate
}

func NewMorseAudio(freq float64, wpm int, farn int) (*MorseAudio, error) {
	ma := new(MorseAudio)

	ma.wpm = wpm

	dur := calcDitDuration(wpm)
	ma.ditDur = time.Duration(dur)

	if farn != 0 {
		fDur := calcDitDuration(farn)
		ma.farnDitDur = time.Duration(fDur)
	}

	// TODO: allow for different wave generators
	silence := generators.Silence(-1)
	ma.silence = silence

	// dit buffer
	ditStr, ditf, err := PreCalcSine(sampleRate, freq, ma.Dit())
	if err != nil {
		return nil, err
	}
	ditBuf := beep.NewBuffer(ditf)
	ditBuf.Append(ditStr)

	dahStr, dahf, err := PreCalcSine(sampleRate, freq, ma.Dash())
	if err != nil {
		return nil, err
	}
	dahBuf := beep.NewBuffer(dahf)
	dahBuf.Append(dahStr)

	ma.dit = ditBuf
	ma.dah = dahBuf

	ma.sr = sampleRate
	
	speaker.Init(ma.sr, int(ma.sr / 10))

	return ma, nil
}

func (ma *MorseAudio) SendMessage(ms morsestrings.MorseString) error {
	ch := make(chan struct{})
	morseSend := make([]beep.Streamer, 0, len(ms) * wordAvg * 2)

	for _, mword := range ms {
		lastChar := mword.Len() - 1
		mw := make([]beep.Streamer, 0, mword.Len() * 2)
		for i, char := range mword.Chars() {
			for _, r := range char {
				switch r {
				case '.':
					mw = append(mw, ma.dit.Streamer(0, ma.dit.Len()))
				case '-':
					mw = append(mw, ma.dah.Streamer(0, ma.dah.Len()))
				default:
					return fmt.Errorf("This should never be able to happen, but somehow '%v' got passed in as a Morse beep!", char)
				}
				mw = append(mw, ma.Silence(ma.Dit()))
			}

			if !mword.IsProsign() && i != lastChar {
				mw = append(mw, ma.Silence(ma.LetterSep()))
			}
		}
		mw = append(mw, ma.Silence(ma.WordSep()))
		morseSend = append(morseSend, mw...)
	}

	morseSend = append(morseSend, beep.Callback(func(){
		ch <- struct{}{}
	}))

	speaker.Play(beep.Seq(morseSend...))
	<-ch

	// and done
	return nil
}

func (ma *MorseAudio) Silence(dur time.Duration) beep.Streamer {
	sampDur := ma.sr.N(dur)
	return beep.Take(sampDur, ma.silence)
}

func (ma *MorseAudio) Dit() time.Duration {
	return ma.ditDur * mDit
}

func (ma *MorseAudio) Dash() time.Duration {
	return ma.ditDur * mDash
}

func (ma *MorseAudio) LetterSep() time.Duration {
	return ma.ditDur * mLetterSep
}

func (ma *MorseAudio) WordSep() time.Duration {
	if ma.farnDitDur != 0 {
		return ma.farnDitDur * mWordSep
	}

	return ma.ditDur * mWordSep
}
