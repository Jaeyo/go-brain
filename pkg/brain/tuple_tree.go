package brain

import (
	"fmt"
	"regexp"
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
					FrequencyWord: frequencyWord,
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

type FrequencyTuple struct {
	FrequencyWord int
	WordCharacter string
	Position      int
}

type FrequencyTuples []FrequencyTuple
