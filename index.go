package main

import (
	"encoding/gob"
	"os"
)

// index is an inverted index. It maps tokens to document IDs.
type index map[string][]int

const SEARCH_INDEX = "index.gob"

// add adds documents to the index.
func (idx index) add(docs []document) {
	for _, doc := range docs {
		for _, token := range analyze(doc.Text) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				// Don't add same ID twice.
				continue
			}
			idx[token] = append(ids, doc.ID)
		}
	}
}

func saveIndex(idx index) error {
	file, err := os.Create(SEARCH_INDEX)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)

	if err = enc.Encode(idx); err != nil {
		return err
	}
	return nil
}

func loadIndex() (index, error) {
	var idx index
	file, err := os.Open(SEARCH_INDEX)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(file)
	if err = dec.Decode(&idx); err != nil {
		return nil, err
	}
	return idx, nil
}

// intersection returns the set intersection between a and b.
// a and b have to be sorted in ascending order and contain no duplicates.
func intersection(a []int, b []int) []int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}
	r := make([]int, 0, maxLen)
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			i++
		} else if a[i] > b[j] {
			j++
		} else {
			r = append(r, a[i])
			i++
			j++
		}
	}
	return r
}

// search queries the index for the given text.
func (idx index) search(text string) []int {
	var r []int
	for _, token := range analyze(text) {
		if ids, ok := idx[token]; ok {
			if r == nil {
				r = ids
			} else {
				r = intersection(r, ids)
			}
		} else {
			// Token doesn't exist.
			return nil
		}
	}
	return r
}
