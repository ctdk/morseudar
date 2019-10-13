package topwords

import (
	"github.com/ctdk/morse-copying/morse"
)

// These words are ganked from `google-10000-english-usa-no-swears.txt`
// in https://github.com/first20hours/google-10000-english. Using the no-swear
// version because that seems best somehow.  

type TopWordList []morse.MorseString

//go:generate sh -c "./gen-topwords.pl < ./google-10000-english-usa-no-swears.txt > wordlist.go && go fmt wordlist.go"

func GetTopWords(num int) TopWordList {
	allCnt := len(topWds)
	if num > allCnt || num == 0 {
		num = allCnt
	}
	wl := topWds[:num]

	tw := make([]morse.MorseString, num)
	for i, w := range wl {
		tw[i] = morse.StringToMorse(w)
	}

	return TopWordList(tw)
}
