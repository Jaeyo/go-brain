package brain

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type TupleTree struct {
	SortedTupleVector       [][]int
	WordCombinations        [][][]int
	WordCombinationsReverse [][][]string
	TupleVector             [][][]int
	GroupLen                []int
}

func NewTupleTree(
	sortedTupleVector [][]int,
	wordCombinations [][][]int,
	wordCombinationsReverse [][][]string,
	tupleVector [][][]int,
	groupLen []int,
) *TupleTree {
	return &TupleTree{
		SortedTupleVector:       sortedTupleVector,
		WordCombinations:        wordCombinations,
		WordCombinationsReverse: wordCombinationsReverse,
		TupleVector:             tupleVector,
		GroupLen:                groupLen,
	}
}

func getFrequencyVector(sentences []string, filter []string, delimiter []string, dataset string) (map[int][][]string, map[int][]FrequencyTuples, map[int][][]int) {
	groupLen := make(map[int][][]string)
	set := make(map[string][]string)
	lineID := 0

	for _, sentence := range sentences {
		for _, rgex := range filter {
			re := regexp.MustCompile(rgex)
			sentence = re.ReplaceAllString(sentence, "<*>")
		}
		for _, de := range delimiter {
			re := regexp.MustCompile(de)
			sentence = re.ReplaceAllString(sentence, "")
		}
		if dataset == "HealthApp" {
			sentence = strings.ReplaceAll(sentence, ":", ": ")
			sentence = strings.ReplaceAll(sentence, "=", "= ")
			sentence = strings.ReplaceAll(sentence, "|", "| ")
		}
		if dataset == "Android" {
			sentence = strings.ReplaceAll(sentence, "(", "( ")
			sentence = strings.ReplaceAll(sentence, ")", ") ")
			sentence = strings.ReplaceAll(sentence, ":", ": ")
			sentence = strings.ReplaceAll(sentence, "=", "= ")
		}
		if dataset == "HPC" {
			sentence = strings.ReplaceAll(sentence, "=", "= ")
			sentence = strings.ReplaceAll(sentence, "-", "- ")
			sentence = strings.ReplaceAll(sentence, ":", ": ")
		}
		if dataset == "BGL" {
			sentence = strings.ReplaceAll(sentence, "=", "= ")
			sentence = strings.ReplaceAll(sentence, "..", ".. ")
			sentence = strings.ReplaceAll(sentence, "(", "( ")
			sentence = strings.ReplaceAll(sentence, ")", ") ")
		}
		if dataset == "Hadoop" {
			sentence = strings.ReplaceAll(sentence, "_", "_ ")
			sentence = strings.ReplaceAll(sentence, ":", ": ")
			sentence = strings.ReplaceAll(sentence, "=", "= ")
			sentence = strings.ReplaceAll(sentence, "(", "( ")
			sentence = strings.ReplaceAll(sentence, ")", ") ")
		}
		if dataset == "HDFS" {
			sentence = strings.ReplaceAll(sentence, ":", ": ")
		}
		if dataset == "Linux" {
			sentence = strings.ReplaceAll(sentence, "=", "= ")
			sentence = strings.ReplaceAll(sentence, ":", ": ")
		}
		if dataset == "Spark" {
			sentence = strings.ReplaceAll(sentence, ":", ": ")
		}
		if dataset == "Thunderbird" {
			sentence = strings.ReplaceAll(sentence, ":", ": ")
			sentence = strings.ReplaceAll(sentence, "=", "= ")
		}
		if dataset == "Windows" {
			sentence = strings.ReplaceAll(sentence, ":", ": ")
			sentence = strings.ReplaceAll(sentence, "=", "= ")
			sentence = strings.ReplaceAll(sentence, "[", "[ ")
			sentence = strings.ReplaceAll(sentence, "]", "] ")
		}
		if dataset == "Zookeeper" {
			sentence = strings.ReplaceAll(sentence, ":", ": ")
			sentence = strings.ReplaceAll(sentence, "=", "= ")
		}
		sentence = strings.ReplaceAll(sentence, ",", ", ")
		sentence = regexp.MustCompile(" +").ReplaceAllString(sentence, " ")
		words := strings.Split(sentence, " ")
		words = append([]string{fmt.Sprint(lineID)}, words...)
		length := 0
		for _, token := range words {
			set[fmt.Sprint(length)] = append(set[fmt.Sprint(length)], token)
			length++
		}
		wordCnt := len(words)
		groupLen[wordCnt] = append(groupLen[wordCnt], words)
		lineID++
	}

	tupleVector := make(map[int][]FrequencyTuples)
	frequencyVector := make(map[int][][]int)
	maxWordCnt := 0
	for keyWordCnt := range groupLen {
		if keyWordCnt > maxWordCnt {
			maxWordCnt = keyWordCnt
		}
	}
	i := 0
	freSet := make(map[string]int)
	for i < maxWordCnt {
		for _, word := range set[fmt.Sprint(i)] {
			wordKey := fmt.Sprintf("%d %s", i, word)
			freSet[wordKey]++
		}
		i++
	}

	for keyWordCnt, groupWords := range groupLen {
		for _, word := range groupWords {
			position := 0
			frequencyTuples := FrequencyTuples{}
			frequencyCommon := []int{}
			skipLineID := 1
			for _, wordCharacter := range word {
				if skipLineID == 1 {
					skipLineID = 0
					continue
				}
				frequencyWord := freSet[fmt.Sprintf("%d %s", position+1, wordCharacter)]
				tuple := FrequencyTuple{
					Frequency:     frequencyWord,
					WordCharacter: wordCharacter,
					Position:      position,
				}
				frequencyTuples = append(frequencyTuples, tuple)
				frequencyCommon = append(frequencyCommon, frequencyWord)
				position++
			}
			tupleVector[keyWordCnt] = append(tupleVector[keyWordCnt], frequencyTuples)
			frequencyVector[keyWordCnt] = append(frequencyVector[keyWordCnt], frequencyCommon)
		}
	}

	return groupLen, tupleVector, frequencyVector
}

func tupleGenerate(groupLen map[int][][]string, tupleVector map[int][]FrequencyTuples, frequencyVector map[int][][]int) (map[int][]FrequencyTuples, map[int][][]WordCounts, map[int][][]WordCounts) {
	sortedTupleVector := make(map[int][]FrequencyTuples)
	wordCombinations := make(map[int][][]WordCounts)
	wordCombinationsReverse := make(map[int][][]WordCounts)

	for key := range groupLen {
		rootSet := map[string]struct{}{"": {}}
		for _, fre := range tupleVector[key] {
			sortedFreReverse := fre.SortReverseByFrequency()
			rootSet[sortedFreReverse[0].WordCharacter] = struct{}{}
			sortedTupleVector[key] = append(sortedTupleVector[key], sortedFreReverse)
		}

		for _, fc := range frequencyVector[key] {
			number := make(map[int]int)
			for _, freq := range fc {
				number[freq]++
			}

			result := make([]WordCount, 0, len(number))
			for k, v := range number {
				result = append(result, WordCount{Word: fmt.Sprint(k), Count: v})
			}

			sortedResult := make([]WordCount, len(result))
			copy(sortedResult, result)
			sort.Slice(sortedResult, func(i, j int) bool {
				return sortedResult[i].Count > sortedResult[j].Count
			})

			sortedFre := make([]WordCount, len(result))
			copy(sortedFre, result)
			sort.Slice(sortedFre, func(i, j int) bool {
				return sortedFre[i].Word > sortedFre[j].Word
			})

			// TODO IMME

		}
	}

	// TODO IMME

}

type FrequencyTuple struct {
	Frequency     int
	WordCharacter string
	Position      int
}

type FrequencyTuples []FrequencyTuple

func (f FrequencyTuples) SortReverseByFrequency() FrequencyTuples {
	tuples := f.clone()
	sort.Slice(tuples, func(i, j int) bool {
		return tuples[i].Frequency > tuples[j].Frequency
	})
	return tuples
}

func (f FrequencyTuples) clone() FrequencyTuples {
	tuples := make(FrequencyTuples, len(f))
	for i, item := range f {
		tuples[i] = item
	}
	return tuples
}

type WordCount struct {
	Word  string
	Count int
}

type WordCounts []WordCount

type Counter[T comparable] struct {
	counts map[T]int
}

func NewCounter[T comparable](items ...T) *Counter[T] {
	counter := &Counter[T]{counts: make(map[T]int)}

	for _, item := range items {
		counter.counts[item]++
	}

	return counter
}

func (c *Counter[T]) MostCommon() []struct {
	Item  T
	Count int
} {
	type kv struct {
		Item  T
		Count int
	}

	var freqList []kv
	for item, count := range c.counts {
		freqList = append(freqList, kv{item, count})
	}

	sort.Slice(freqList, func(i, j int) bool {
		return freqList[i].Count > freqList[j].Count
	})

	result := make([]struct {
		Item  T
		Count int
	}, len(freqList))
	for i, kv := range freqList {
		result[i] = struct {
			Item  T
			Count int
		}{
			Item:  kv.Item,
			Count: kv.Count,
		}
	}

	return result
}
