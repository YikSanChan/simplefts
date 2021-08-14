package main

import (
	"compress/gzip"
	"encoding/gob"
	"encoding/xml"
	"os"
)

const DOCS = "docs.gob"

// document represents a Wikipedia abstract dump document.
type document struct {
	Title string `xml:"title"`
	URL   string `xml:"url"`
	Text  string `xml:"abstract"`
	ID    int
}

// loadDocuments loads a Wikipedia abstract dump and returns a slice of documents.
// Dump example: https://dumps.wikimedia.org/enwiki/latest/enwiki-latest-abstract1.xml.gz
func loadDocuments(path string) ([]document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	dec := xml.NewDecoder(gz)
	dump := struct {
		Documents []document `xml:"doc"`
	}{}
	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}
	docs := dump.Documents
	for i := range docs {
		docs[i].ID = i
	}
	return docs, nil
}

func saveDocs(docs []document) error {
	file, err := os.Create(DOCS)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)

	if err = enc.Encode(docs); err != nil {
		return err
	}
	return nil
}

func loadDocs() ([]document, error) {
	var docs []document
	file, err := os.Open(DOCS)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(file)
	if err = dec.Decode(&docs); err != nil {
		return nil, err
	}
	return docs, nil
}
