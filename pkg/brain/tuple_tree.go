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

func getFrequencyVector(sentences []string, filter []string, delimiter []string, dataset string) (map[int][][]string, map[int][][]Tuple, map[int][][]int) {
	// map key (int): tokens count (`len(tokens`)
	// map value ([][]string): tokens list
	groupLen := make(map[int][]Tokens)
	// map key (string): columnar index for tokens
	// map value ([]string): tokens
	set := make(map[string]Tokens)
	lineID := 0

	for _, s := range sentences {
		s = applyFiltersAndDelimiters(s, filter, delimiter, dataset)
		s = strings.ReplaceAll(s, ",", ", ")
		s = regexp.MustCompile(" +").ReplaceAllString(s, " ")
		tokens := strings.Split(s, " ")
		tokens = append([]string{fmt.Sprint(lineID)}, tokens...)

		for i, token := range tokens {
			set[fmt.Sprint(i)] = append(set[fmt.Sprint(i)], token)
		}

		lena := len(tokens)
		groupLen[lena] = append(groupLen[lena], tokens) // first grouping: logs with the same length
		lineID++
	}

	tupleVector := make(map[int][][]Tuple)
	frequencyVector := make(map[int][][]int)
	maxTokenLength := maxKey(groupLen) // a: the biggest length of the log in this dataset
	// map key (string): `{i} {word}`, i: columnar index for tokens, word: token
	// map value (int): count
	freSet := make(map[string]int) // saving each word's frequency

	for i := 0; i < maxTokenLength; i++ {
		for _, word := range set[fmt.Sprint(i)] {
			key := fmt.Sprint(i) + " " + word
			freSet[key]++
		}
	}

	for key, group := range groupLen {
	
	}

	_ = tupleVector
	_ = frequencyVector
	return nil, nil, nil
	// TODO IMME
}

func applyFiltersAndDelimiters(s string, filter []string, delimiter []string, dataset string) string {
	for _, rgex := range filter {
		s = regexp.MustCompile(rgex).ReplaceAllString(s, "<*>")
	}
	for _, de := range delimiter {
		s = strings.ReplaceAll(s, de, "")
	}

	switch dataset {
	case "HealthApp":
		s = applyDatasetSpecificReplacements(s, []string{":", "=", "\\|"}, ": ", "= ", "| ")
	case "Android":
		s = applyDatasetSpecificReplacements(s, []string{"\\(", "\\)"}, "( ", ") ")
		s = applyDatasetSpecificReplacements(s, []string{":", "="}, ": ", "= ")
	case "HPC":
		s = applyDatasetSpecificReplacements(s, []string{"=", "-", ":"}, "= ", "- ", ": ")
	case "BGL":
		s = applyDatasetSpecificReplacements(s, []string{"=", "\\.\\.", "\\(", "\\)"}, "= ", ".. ", "( ", ") ")
	case "Hadoop":
		s = applyDatasetSpecificReplacements(s, []string{"_", ":", "=", "\\(", "\\)"}, "_ ", ": ", "= ", "( ", ") ")
	case "HDFS":
		s = applyDatasetSpecificReplacements(s, []string{":"}, ": ")
	case "Linux":
		s = applyDatasetSpecificReplacements(s, []string{"=", ":"}, "= ", ": ")
	case "Spark":
		s = applyDatasetSpecificReplacements(s, []string{":"}, ": ")
	case "Thunderbird":
		s = applyDatasetSpecificReplacements(s, []string{":", "="}, ": ", "= ")
	case "Windows":
		s = applyDatasetSpecificReplacements(s, []string{":", "=", "\\[", "\\]"}, ": ", "= ", "[ ", "] ")
	case "Zookeeper":
		s = applyDatasetSpecificReplacements(s, []string{":", "="}, ": ", "= ")
	}
	return s
}

func applyDatasetSpecificReplacements(s string, patterns []string, replacements ...string) string {
	for i, pattern := range patterns {
		s = regexp.MustCompile(pattern).ReplaceAllString(s, replacements[i])
	}
	return s
}

type Tuple struct {
	Frequency     int
	WordCharacter string
	Position      int
}

type Tokens []string
