package main

import (
	"flag"
	"log"
	"time"
)

func main() {
	var dumpPath, query string
	var rebuildDocs, rebuildIndex bool

	flag.StringVar(&dumpPath, "p", "enwiki-latest-abstract1.xml.gz", "wiki abstract dump path")
	flag.StringVar(&query, "q", "Small wild cat", "search query")
	flag.BoolVar(&rebuildDocs, "rd", false, "rebuild docs or not")
	flag.BoolVar(&rebuildIndex, "ri", false, "rebuild search index or not")
	flag.Parse()

	log.Println("Starting simplefts")

	var err error
	var idx index
	var docs []document

	if rebuildDocs {
		start := time.Now()
		docs, err = loadDocuments(dumpPath)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Loaded %d documents in %v", len(docs), time.Since(start))
		saveDocs(docs)
		log.Println("Saved built docs")
	} else {
		docs, err = loadDocs()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Loaded prebuilt docs")
	}

	if rebuildIndex {
		start := time.Now()
		idx := make(index)
		idx.add(docs)
		log.Printf("Indexed %d documents in %v", len(docs), time.Since(start))
		saveIndex(idx)
		log.Println("Saved built index")
	} else {
		idx, err = loadIndex()
		if err != nil {
			panic(err)
		}
		log.Println("Loaded prebuilt index")
	}

	start := time.Now()
	matchedIDs := idx.search(query)
	log.Printf("Search found %d documents in %v", len(matchedIDs), time.Since(start))

	for _, id := range matchedIDs {
		doc := docs[id]
		log.Printf("%d\t%s\n", id, doc.Text)
	}
}
