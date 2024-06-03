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

package main

import (
	"bufio"
	"fmt"
	"github.com/ctdk/morseudar/internal/morse"
	"github.com/ctdk/morseudar/internal/codegroups"
	"github.com/ctdk/morseudar/internal/copy-compare"
	"github.com/ctdk/morseudar/internal/morsestrings"
	"github.com/ctdk/morseudar/internal/stats"
	"github.com/ctdk/morseudar/internal/textblock"
	"github.com/ctdk/morseudar/internal/wordlists"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"
	"unicode"
)

const version = "0.0.1"

type Options struct {
	Version bool `short:"v" long:"version" description:"Print version info."`
	Wpm int `short:"w" long:"wpm" description:"Words per minute. Defaults to 10."`
	Farnsworth int `short:"o" long:"farnsworth" description:"Farnsworth timing. Words are sent at the speed given with -w/--wpm, but the spaces between words are sent at this WPM. For instance, -w 20 -o 10 would send words at 20 wpm, but spaced out as if they were sent at 10 wpm, giving you more time to process."`
	Frequency int `short:"f" long:"frequency" description:"Frequency in Hz for Morse beep. Defaults to 700."`
	Mode string `short:"m" long:"mode" description:"Mode to run morseudar under. Options include: text (requires -t/--text), randomline (also requires -t/--text), codegroups, codealnum, codenumbers, topwords, qcodes, chars. Defaults to topwords."`
	Text string `short:"t" long:"text" description:"Path to text file to load and use for copying testing. Required for 'text' mode."`
	SaveFile string `short:"s" long:"save" description:"Specify path to save file holding previous test results to help keep track of your progress."` // not ready yet
	TopWordNum int `short:"n" long:"top-word-num" description:"How many words from the top word list to include. Only relevant in topwords mode."`
	Qquestions bool `short:"q" long:"qcode-questions" description:"Include Q codes followed by a question mark (i.e. QRS and QRS?)."`
	Seq bool `short:"r" long:"sequential" description:"Send lines sequentially instead of randomly. Not relevant for the code group modes."`
	EntireBlock bool `short:"b" long:"entire-block" description:"Send the entire block of text at once, rather than line by line. Unsurprisingly, only relevant for -t/--text." hidden:"true"` // not ready
	PrintStats bool `short:"P" long:"print-stats" description:"Print out user statistics and exit."`
}

func main() {
	var opts = &Options{}
	parser := flags.NewParser(opts, flags.Default)
	parser.ShortDescription = fmt.Sprintf("A Morse code copying testing program, version %s.", version)

	if _, err := parser.Parse(); err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Println(err)
			os.Exit(1)
		}
	}

	if opts.Version {
		fmt.Printf("morseudar version %s, built with %s.\n", version, runtime.Version())
		os.Exit(0)
	}

	uStats, err := stats.Load(opts.SaveFile)
	if err != nil {
		log.Fatal(err)
	}

	if opts.PrintStats {
		// TODO: make this better later
		for _, st := range uStats.Summaries {
			fmt.Println(st)
		}
		os.Exit(0)
	}

	var mode morse.MorseMode

	switch strings.ToLower(opts.Mode) {
	case "codegroups":
		mode = morse.CodeGroup
	case "codealnum":
		mode = morse.CodeAlnum
	case "codenumbers":
		mode = morse.CodeNum
	case "topwords":
		mode = morse.TopWords
	case "qcodes":
		mode = morse.Qcode
	case "chars":
		mode = morse.MorseChar
	default:
		mode = morse.TextFile
	}


	// make a morse object!

	m, err := morse.New(mode, opts.Wpm, opts.Farnsworth, float64(opts.Frequency), opts.Seq, opts.EntireBlock, 0)
	if err != nil {
		log.Fatal(err)
	}

	// attach the Stone of Triumph
	switch mode {
	case morse.CodeGroup:
		// set up proper length options later
		m.TestingMaterial = codegroups.NewCodegroup(m.Src(), codegroups.Alpha, 0, 0)
	case morse.CodeAlnum:
		// set up proper length options later
		m.TestingMaterial = codegroups.NewCodegroup(m.Src(), codegroups.Alnum, 0, 0)
	case morse.CodeNum:
		// set up proper length options later
		m.TestingMaterial = codegroups.NewCodegroup(m.Src(), codegroups.Num, 0, 0)
	case morse.TopWords:
		m.TestingMaterial = wordlists.GetTopWords(opts.TopWordNum, m.Src())
	case morse.Qcode:
		m.TestingMaterial = wordlists.GetQCodes(opts.Qquestions, m.Src())
	case morse.MorseChar:
		m.TestingMaterial = wordlists.GetChars(m.Src())
	case morse.TextFile:
		tb := textblock.NewTextblock(m.Src())
		if err := tb.LoadFile(opts.Text); err != nil {
			log.Fatal(err)
		}
		m.TestingMaterial = tb
	}

	comp := compare.New()
	reader := bufio.NewReader(os.Stdin)
	l := 1
	answers := make(compare.AnswerBatch, 0)

	handleSignals(&answers)

	for {
		ml, _ := m.GetMorse()
MorseLoop:
		fmt.Printf("# %d\n", l)
		m.Send(ml)
		start := time.Now()
		var guess string

		tries := 0
		for {
			tries++
			fmt.Print("> ")
			guess, _ = reader.ReadString('\n')
			if guess == "" || guess == "\r\n" {
				fmt.Println("?")
				goto MorseLoop
			} else {
				break
			}
		}

		guess = strings.ToLower(strings.TrimSpace(guess))

		// Processing input may be better handled as a function or
		// method. TODO later.

		// commands start with ` since that character's not going to
		// come up as Morse code.
		if strings.HasPrefix(guess, "`") {
			switch guess {
			case "`quit", "`exit":
				fmt.Println("Saving and exiting...")
				perc, dur, tries := answers.Averages()
				sum := stats.NewSummary(time.Now(), mode, perc, dur, tries, l, opts.Wpm, opts.Farnsworth)
				fmt.Println(sum)
				uStats.Add(sum)
				if err = uStats.Save(); err != nil {
					log.Fatal(err)
				}
				os.Exit(0)
			case "`stats":
				fmt.Printf("Statistics for '%s':\n\n", uStats.Username)
				for _, st := range uStats.Summaries {
					fmt.Println(st)
				}
			case "`remind":
				// sort round
				i := 0
				k := make([]rune, len(morsestrings.Alphabet))
				for r := range morsestrings.Alphabet {
					k[i] = r
					i++
				}

				// This seems like it should be able to do in a
				// cleaner fashion, but I'm spacing on how
				// exactly.
				runeSort := func (i, j int) bool {
					if unicode.IsLetter(k[i]) {
						if unicode.IsNumber(k[j]) || unicode.IsPunct(k[j]) || unicode.IsSymbol(k[j]) {
							return true
						}
					} else if unicode.IsNumber(k[i]) {
						if unicode.IsLetter(k[j]) {
							return false
						} else if unicode.IsPunct(k[j]) || unicode.IsSymbol(k[j]) {
							return true
						}
					} else if unicode.IsPunct(k[i]) || unicode.IsSymbol(k[i]) {
						if unicode.IsLetter(k[j]) || unicode.IsNumber(k[j]) {
							return false
						}
					}

					// fallthrough, they're the same class
					// of character
					return k[i] < k[j]
				}
				sort.Slice(k, runeSort)

				tw := new(tabwriter.Writer)
				tw.Init(os.Stdout, 9, 8, 2, ' ', 0)
				i = 1
				for _, r := range k {
					if r == '0' || r == '"' {
						i = 1
						fmt.Fprint(tw, "\n\n")
					}
					fmt.Fprintf(tw, "%c| %s\t", r, morsestrings.Alphabet[r])
					if i % 5 == 0 {
						fmt.Fprintln(tw)
					}
					i++
				}
				fmt.Fprintln(tw)
				tw.Flush()
			default:
				fmt.Printf("Unknown command '%s'.\n", guess)
			}
			goto MorseLoop
		}

		ans := comp.Compare(ml.RawString(), guess, start, tries)
		fmt.Printf("'%s' was %.2f%% correct. Took %d tries over %s. Original: '%s'\n", guess, ans.Percentage * 100, ans.Tries, ans.Took.Round(time.Second / 100), ml.RawString())
		answers = append(answers, ans)
		l++
	}
	
}

func handleSignals(answers *compare.AnswerBatch) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for _ = range c {
			fmt.Println("Exiting...")
			// do cleanup later
			perc, dur, tries := answers.Averages()
			fmt.Printf("Averages: %.2f%% correct | %s | %.1f tries\n", perc * 100, dur.Round(time.Second / 100), tries)
			os.Exit(0)
		}
	}()
}
