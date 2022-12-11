package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type wordFreq struct {
	Word string
	Freq int
}

func Top10(input string) []string {
	if input == "" {
		var void []string
		return void
	}
	inputSR := []rune(input)
	space, _ := utf8.DecodeRune([]byte(` `))
	for index, rune := range inputSR {
		if unicode.IsControl(rune) {
			inputSR[index] = space
		}
	}

	input = string(inputSR)

	slice := strings.Split(input, " ")
	freqSlice := make([]wordFreq, 0)
	freqMap := make(map[string]int)

	for _, word := range slice {
		if word == "" {
			continue
		}
		if freqMap[word] == 0 {
			freqSlice = append(freqSlice, wordFreq{
				Word: word,
				Freq: 1,
			})
			freqMap[word] = len(freqSlice) - 1
		} else {
			freqSlice[freqMap[word]].Freq++
		}
	}

	sort.Slice(freqSlice, func(i, j int) bool {
		if freqSlice[i].Freq > freqSlice[j].Freq {
			return true
		} else if freqSlice[i].Freq < freqSlice[j].Freq {
			return false
		} else {
			if strings.Compare(freqSlice[i].Word, freqSlice[j].Word) == -1 {
				return true
			}
			return false
		}
	})

	resslice := make([]string, 10)
	for i := 0; i < 10; i++ {
		resslice[i] = freqSlice[i].Word
	}

	return resslice
}
