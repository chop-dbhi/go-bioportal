package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var vocabHeader = []string{
	"id:ID(Vocabulary)",
	"label",
}

var classHeader = []string{
	"id:ID(Class)",
	"label",
	"code",
	"synonyms:string[]",
	"definitions:string[]",
	"obsolete:boolean",
	"cui:string[]",
	"semanticTypes:string[]",
}

var subClassOfRelHeader = []string{
	":START_ID(Class)",
	":TYPE",
	":END_ID(Class)",
}

var classOfRelHeader = []string{
	":START_ID(Class)",
	":TYPE",
	":END_ID(Vocabulary)",
}

func newFile(name string, args ...interface{}) *os.File {
	f, err := os.Create(fmt.Sprintf(name, args...))
	if err != nil {
		log.Fatal(f)
	}
	return f
}

func main() {
	log.SetFlags(0)

	args := os.Args[1:]
	if len(args) < 3 {
		log.Fatal("vocabulary id, label, and path required")
	}

	vocabID := args[0]
	vocabLabel := args[1]
	vocabPath := args[2]

	f, err := os.Open(vocabPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	cr := csv.NewReader(f)
	_, err := cr.Read()
	if err != nil {
		log.Fatalf("error reading row: %s", err)
	}

	f := newFile("%s_nodes_vocabulary.csv", vocabID)

	cw := csv.NewWriter(f)
	err = cw.WriteAll([][]string{
		vocabHeader,
		[]string{vocabID, vocabLabel},
	})
	if err != nil {
		log.Fatal(err)
	}
	cw.Flush()
	f.Close()

	f = newFile("%s_nodes_class.csv", vocabID)
	defer f.Close()
	cw = csv.NewWriter(f)
	defer cw.Flush()
	if err := cw.Write(classHeader); err != nil {
		log.Fatal(err)
	}

	vrf := newFile("%s_rels_classof.csv", vocabID)
	defer vrf.Close()
	vrcw := csv.NewWriter(vrf)
	defer vrcw.Flush()
	if err := vrcw.Write(classOfRelHeader); err != nil {
		log.Fatal(err)
	}

	rf := newFile("%s_rels_subclassof.csv", vocabID)
	defer rf.Close()
	rcw := csv.NewWriter(rf)
	defer rcw.Flush()
	if err := rcw.Write(subClassOfRelHeader); err != nil {
		log.Fatal(err)
	}

	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error reading row: %s", err)
		}

		id := row[0]
		if id == "" {
			log.Printf("skipping: %v", row)
			continue
		}

		var parents []string
		if row[7] != "" {
			parents = strings.Split(row[7], "|")
		}

		toks := strings.Split(id, "/")
		code := toks[len(toks)-1]

		err = cw.Write([]string{
			id,
			row[1],
			code,
			row[2],
			row[3],
			row[4],
			row[5],
			row[6],
		})

		if err != nil {
			log.Fatal(err)
		}

		err = vrcw.Write([]string{
			id,
			"classOf",
			vocabID,
		})
		if err != nil {
			log.Fatal(err)
		}

		for _, pid := range parents {
			err = rcw.Write([]string{
				id,
				"subClassOf",
				pid,
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Printf(`./neo4j/bin/neo4j-import \
	--into ./neo4j/data/databases/graph.db \
	--delimiter ',' \
	--array-delimiter '|' \
	--quote '"' \
	--nodes:Vocabulary "%s_nodes_vocabulary.csv" \
	--nodes:Class "%s_nodes_class.csv" \
	--relationships "%s_rels_subclassof.csv" \
	--relationships "%s_rels_classof.csv"
	`, vocabID, vocabID, vocabID, vocabID)
}
