morseudar
=============

`morseudar` is a program that lets you practice copying Morse code played over a speaker on your computer. Currently, it is not strongly intended for learning the characters themselves but rather to practice copying lines of text.

Motivation
----------

There are many fine applications out there for learning Morse code, especially on mobile devices, but I found that there was a lack of apps that would work with a keyboard. While this may not be the case for everyone, I'm able to type faster than I can write and I can absolutely type faster with all of my fingers than I can with just my thumbs.

Features
--------

* Configurable WPM timing, defaulting to 10 wpm.
* Adjustable beep frequency, defaulting to 700 Hz.
* Optional Farnsworth timing. This means that while the words themselves are sent at one rate, the *spacing* between the words is sent as if it were a slower rate of words per minute.
* Different modes to choose from. Modes include the top X words in English, code groups, individual characters, Q codes, and text from arbitrary files.
* When your answer is compared to the original line sent, it's not an either/or comparison. Rather than missing one character absolutely derailing everything, you'll get partial credit for the answer.
* Session tatistics! At the end of a session, `morseudar` will print out a set of statistics on how you did, including average percentage correct, average time taken to answer, and the average number of tries you took to answer correctly.
* Statistics over time (in progress). Keep track of how you're doing over time.
* Command-line goodness. Instead of having a GUI, it happily runs in a terminal window and just does its job.

Usage
-----

	Usage:
	  morseudar [OPTIONS]

	Application Options:
	  -v, --version          Print version info.
	  -w, --wpm=             Words per minute. Defaults to 10.
	  -o, --farnsworth=      Farnsworth timing. Words are sent at the speed given
				 with -w/--wpm, but the spaces between words are sent
				 at this WPM. For instance, -w 20 -o 10 would send
				 words at 20 wpm, but spaced out as if they were sent
				 at 10 wpm, giving you more time to process.
	  -f, --frequency=       Frequency in Hz for Morse beep. Defaults to 700.
	  -m, --mode=            Mode to run morseudar under. Options include: text
				 (requires -t/--text), randomline (also requires
				 -t/--text), codegroups, codealnum, codenumbers,
				 topwords, qcodes, chars. Defaults to topwords.
	  -t, --text=            Path to text file to load and use for copying testing.
				 Required for 'text' mode.
	  -n, --top-word-num=    How many words from the top word list to include. Only
				 relevant in topwords mode.
	  -q, --qcode-questions  Include Q codes followed by a question mark (i.e. QRS
				 and QRS?).
	  -r, --sequential       Send lines sequentially instead of randomly. Not
				 relevant for the code group modes.

	Help Options:
	  -h, --help             Show this help message

TODO
----

See the TODO file.

BUGS
----

Surely they exist. Check the BUGS file.

LICENSE
-------

Copyright 2019-2024, Jeremy Bingham, under the terms of the Apache License Version 2.0.
