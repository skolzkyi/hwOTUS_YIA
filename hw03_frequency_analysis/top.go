package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`[!,.;:?)"(]`)

type wordFreq struct {
	Word string
	Freq int
}

func Top10(input string) []string {
	var resLength int
	if input == "" {
		var void []string
		return void
	}

	slice := strings.Fields(input)
	freqSlice := make([]wordFreq, len(slice))
	freqMap := make(map[string]int)

	i := 0
	for _, rawWord := range slice {
		word := reg.ReplaceAllString(strings.ToLower(rawWord), "")
		if word == "-" {
			continue
		}
		_, ok := freqMap[word]
		if !ok {
			freqSlice[i] = wordFreq{
				Word: word,
				Freq: 1,
			}
			freqMap[word] = i
			i++
		} else {
			freqSlice[freqMap[word]].Freq++
		}
	}
	freqSlice = freqSlice[:i+1]
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
	if len(freqSlice) >= 10 {
		resLength = 10
	} else {
		resLength = len(freqSlice) - 1
	}
	resslice := make([]string, resLength)
	for i := 0; i < resLength; i++ {
		resslice[i] = freqSlice[i].Word
	}

	return resslice
}
