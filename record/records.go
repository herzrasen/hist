package record

import "sort"

type Records []Record

func (r Records) Sort() {
	sort.Slice(r, func(i, j int) bool {
		left := r[i]
		right := r[j]
		if left.Weight == right.Weight {
			return left.LastUpdate.After(right.LastUpdate)
		}
		return left.Weight > right.Weight
	})
}

func (r Records) Search(input string) []Record {
	var newRecs Records
	for _, rec := range r {
		rec.UpdateWeight(input)
		if rec.Weight > 0 || input == "" {
			newRecs = append(newRecs, rec)
		}
	}
	newRecs.Sort()
	return newRecs
}
