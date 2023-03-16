package fuzzy

import (
	"github.com/herzrasen/hist/record"
	"sort"
	"strings"
	"unicode"
)

type WeightedRecord struct {
	record.Record
	Parts   []string
	Compact string
	Weight  int
}

type WeightedRecords []WeightedRecord

func Search(records []record.Record, input string) []WeightedRecord {
	var weightedRecs WeightedRecords
	inputOccurrences := countOccurrences(input)
	inputTuples := splitIntoChunks(input, 2)
	inputTriples := splitIntoChunks(input, 3)
	for _, r := range records {
		compact := strings.Join(strings.Fields(r.Command), "")
		occurrences := countOccurrences(compact)
		weight := weightOccurrences(occurrences, inputOccurrences)
		if weight > 0 {
			tuples := splitIntoChunks(compact, 2)
			triples := splitIntoChunks(compact, 3)
			weight += weightChunks(tuples, inputTuples) +
				weightChunks(triples, inputTriples) +
				int(r.Count)
		}
		weightedRecs = append(weightedRecs, WeightedRecord{
			Record:  r,
			Compact: compact,
			Weight:  weight,
		})
	}
	if len(input) == 0 {
		return weightedRecs
	}
	return weightedRecs.Sort()
}

func countOccurrences(input string) map[rune]uint32 {
	var occurrences = make(map[rune]uint32)
	for _, c := range []byte(input) {
		r := unicode.ToLower(rune(c))
		if _, ok := occurrences[r]; ok {
			occurrences[r]++
		} else {
			occurrences[r] = 1
		}
	}
	return occurrences
}

func weightOccurrences(occurrences map[rune]uint32, input map[rune]uint32) int {
	weight := 0
	for k, v := range input {
		runesInOccurrences := occurrences[k]
		switch {
		case runesInOccurrences == 0:
			return 0
		case runesInOccurrences > v:
			weight += int(v)
		case v >= runesInOccurrences:
			weight += int(runesInOccurrences)
		}
	}
	return weight
}

func weightChunks(tuples []string, input []string) int {
	weight := 0
	for _, t := range tuples {
		for _, i := range input {
			if t == i {
				weight += len(i)
			}
		}
	}
	return weight
}

func splitIntoChunks(input string, length int) []string {
	var result []string
	for i := 0; i < len(input); i += length {
		if i+1 < len(input) {
			result = append(result, input[i:i+2])
		} else {
			result = append(result, input[i:])
		}
	}
	return result
}

func (w *WeightedRecords) Sort() []WeightedRecord {
	var recs []WeightedRecord
	for _, wr := range *w {
		if wr.Weight > 0 {
			recs = append(recs, wr)
		}
	}
	sort.Slice(recs, func(i, j int) bool {
		left := recs[i]
		right := recs[j]
		if left.Weight == right.Weight {
			return left.LastUpdate.After(right.LastUpdate)
		}
		return left.Weight > right.Weight
	})
	return recs
}
