package hw03frequencyanalysis

import (
	"sort"
	"strings"
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

	slice := strings.Fields(input)
	freqSlice := make([]wordFreq, 0)
	freqMap := make(map[string]int)

	for _, word := range slice {
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
		switch {
		case freqSlice[i].Freq > freqSlice[j].Freq:
			return true
		case freqSlice[i].Freq < freqSlice[j].Freq:
			return false
		default:
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
