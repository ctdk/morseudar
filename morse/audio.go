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

package morse

import (
	"fmt"
	"math"
	"github.com/gordonklaus/portaudio"
	"time"
)

// The square wave code is derived from a gist located at
// https://gist.github.com/zach-klippenstein/5eaa2c96c2620b9e9cdb.

const (
	defaultFrequency = 700
	defaultAmplitudeLevel = 5
	maxAmplitudeLevel = 10
	minAmplitudeLevel = 0
	audioQueueLen = 2044 // was 1022, going to try an experiment
)

type squareWave struct {
	frequency uint
	sampleRate float64
	period uint
	amplitude int32

	output chan<- int32
}

// the morse audio type should be able to toggle the audio stream on and off.

// TODO: also ensure cleanup on exit

// MorseAudio drops the beeps.
type MorseAudio struct {
	sq *squareWave
	wpm int
	ditDur time.Duration
	audioChan chan int32
	stream *portaudio.Stream
}

func newSquareWave(output chan<- int32, freq uint, sampleRate float64, amplitudeLevel uint8) (*squareWave, error) {
	period, err := calculatePeriod(freq, sampleRate)
	if err != nil {
		return nil, err
	}

	amp := calculateAmplitude(amplitudeLevel)
	fmt.Printf("amplitude level: %d amplitude: %d\n", amplitudeLevel, amp)

	sq := &squareWave{frequency: freq, sampleRate: sampleRate, output: output, amplitude: amp, period: period}

	return sq, nil
}

func calculateAmplitude(amplitudeLevel uint8) int32 {
	div := (math.MaxUint8 - amplitudeLevel) + 1
	amp := math.MaxInt32 / int32(div)
	return amp
}

func calculatePeriod(freq uint, sampleRate float64) (uint, error) {
	k := 1 / sampleRate;

	p := 2 * (1 / (float64(freq) * k))
	
	if p < 1 {
		err := fmt.Errorf("frequency %d is too high - the wave period would evaluate to 0 (and you wouldn't be hearing it anyway, and would drive everyone around you mad).", freq)
		return 0, err
	}
	if p > (sampleRate / 2) {
		err := fmt.Errorf("frequency %d is too low - the wave period would equal the sample rate (and how would you hear it anyway?)", freq)
		return 0, err
	}

	return uint(p), nil
}

func (sq *squareWave) generate() {
	// no period or amplitude channels for now, may add later
	// TODO: toggling this on or off might be useful down the road.
	for {
		for i := uint(0); i < sq.period / 2; i++ {
			sq.output <- sq.amplitude
		}
		for i := uint(0); i < sq.period / 2; i++ {
			sq.output <- 0
		}
	}
}

// can this be a non-exported function?
func audioCallbackFromChan(audioChan <-chan int32) func([]int32) {
	return func(out []int32) {
		for i := range out {
			out[i] = <-audioChan
		}
	}
}

func NewMorseAudio(freq uint, amplitudeLevel uint8, wpm int) (*MorseAudio, error) {
	ma := new(MorseAudio)

	// initialize the audio stream
	portaudio.Initialize()
	h, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}

	sampleRate := h.DefaultOutputDevice.DefaultSampleRate
	audioChan := make(chan int32, audioQueueLen)
	sq, err := newSquareWave(audioChan, freq, sampleRate, amplitudeLevel)
	if err != nil {
		return nil, err
	}

	go sq.generate()

	ma.sq = sq
	ma.audioChan = audioChan
	ma.wpm = wpm

	dur := CalcDitDuration(wpm)
	ma.ditDur = time.Duration(dur)

	lowLatencyParams := portaudio.LowLatencyParameters(nil, h.DefaultOutputDevice)

	stream, err := portaudio.OpenStream(lowLatencyParams, audioCallbackFromChan(ma.audioChan))
	if err != nil {
		return nil, err
	}
	ma.stream = stream

	return ma, nil
}

// don't forget to call this after you're done
func (ma *MorseAudio) Destroy() {
	portaudio.Terminate()
	ma.stream.Close()
}

func (ma *MorseAudio) SendMessage(ms MorseString) error {
	for _, mword := range ms {
		fmt.Printf("word: %s prosign: %v\n", mword.word, mword.prosign)
		lastChar := len(mword.word) - 1
		for i, char := range mword.word {
			fmt.Printf("char: %s\n", char)
			for _, r := range char {
				fmt.Printf("c: %c\n", r)
				var dur time.Duration
				switch r {
				case '.':
					dur = mDit
				case '-':
					dur = mDash
				default:
					return fmt.Errorf("This should never be able to happen, but somehow '%v' got passed in as a Morse beep!", char)
				}
				if err := ma.playBeep(dur); err != nil {
					return err
				}
			}

			if !mword.prosign && i != lastChar {
				fmt.Printf("char sep\n")
				time.Sleep(mLetterSep * ma.ditDur)
			}
		}
		fmt.Printf("word sep\n")
		time.Sleep(mWordSep * ma.ditDur)
	}

	// and done
	return nil
}

func (ma *MorseAudio) playBeep(dur time.Duration) error {
	if err := ma.stream.Start(); err != nil {
		return err
	}

	time.Sleep(dur * ma.ditDur)

	if err := ma.stream.Stop(); err != nil {
		return err
	}
	return nil
}
