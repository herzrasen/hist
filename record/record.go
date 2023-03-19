package record

import (
	"github.com/fatih/color"
	"strings"
	"time"
	"unicode"
)

type Record struct {
	Id         int64
	Command    string
	LastUpdate time.Time
	Count      uint64
	Weight     uint64
}

type FormatOptions struct {
	NoLastUpdate bool
	NoCount      bool
	WithId       bool
}

func (r *Record) Format(options FormatOptions) string {
	buf := strings.Builder{}
	if !options.NoLastUpdate {
		buf.WriteString(color.GreenString("%s\t", r.LastUpdate.Format(time.RFC1123)))
	}
	if !options.NoCount {
		buf.WriteString(color.BlueString("%d\t", r.Count))
	}
	if options.WithId {
		buf.WriteString(color.YellowString("%d\t", r.Id))
	}
	buf.WriteString(r.Command)
	return buf.String()
}

func (r *Record) UpdateWeight(input string) {
	inputOccurrences := countOccurrences(input)
	inputTuples := splitIntoChunks(input, 2)
	inputTriples := splitIntoChunks(input, 3)
	compact := strings.Join(strings.Fields(r.Command), "")
	occurrences := countOccurrences(compact)
	r.Weight = weightOccurrences(occurrences, inputOccurrences)
	if r.Weight > 0 {
		tuples := splitIntoChunks(compact, 2)
		triples := splitIntoChunks(compact, 3)
		r.Weight += weightChunks(tuples, inputTuples) +
			weightChunks(triples, inputTriples) +
			r.Count
	}
}

func countOccurrences(input string) map[rune]uint64 {
	var occurrences = make(map[rune]uint64)
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

func weightOccurrences(occurrences map[rune]uint64, input map[rune]uint64) uint64 {
	var weight uint64 = 0
	for k, v := range input {
		runesInOccurrences := occurrences[k]
		switch {
		case runesInOccurrences == 0:
			return 0
		case runesInOccurrences > v:
			weight += v
		case v >= runesInOccurrences:
			weight += runesInOccurrences
		}
	}
	return weight
}

func weightChunks(tuples []string, input []string) uint64 {
	var weight uint64 = 0
	for _, t := range tuples {
		for _, i := range input {
			if t == i {
				weight += uint64(len(i))
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
