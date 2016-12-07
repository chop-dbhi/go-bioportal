# Go-Bioportal

Go client for the [Bioontology Bioportal](http://bioportal.bioontology.org/) API.

## Tools

### Bioportal to Neo4j

Takes a standard Bioportal-based CSV file of a vocabulary and generates CSV files and a script to import using the `neo4j-import` tool. Child-parent relationships are related using `subClassOf` and class-vocabulary relationship is modeled as `classOf`.


#### Install

Requires Go to install, however if you want pre-built binaries [please request it](https://github.com/chop-dbhi/go-bioportal/issues/)!

```
go get github.com/chop-dbhi/go-bioportal/cmd/bioportal-to-neo4j
```

#### Usage

The command takes a vocab ID, label, and the file. It will output a set of files and print a script to standard out. The paths in the script will likely need to be adjusted, but it should serve as a template.

```
bioportal-to-neo4j <vocab-id> <vocab-label> <vocab-file>
```

#### Example

- Download the ICD10-CM vocabulary
- Run the tool
- Show the output

```
$ curl -L http://data.bioontology.org/ontologies/ICD10CM/download?apikey=8b5b7825-538d-40e0-9e9e-5ab9274a9aeb&download_format=csv > icd10cm.csv

$ bioportal-to-neo4j \
  icd10cm \
  "International Classification of Diseases, Version 10 - Clinical Modification" \
  icd10cm.csv > icd10cm_load.sh

$ cat icd10cm_load.sh
./neo4j/bin/neo4j-import \
  --into ./neo4j/data/databases/graph.db \
  --delimiter ',' \
  --array-delimiter '|' \
  --quote '"' \
  --nodes:Vocabulary "icd10cm_nodes_vocabulary.csv" \
  --nodes:Class "icd10cm_nodes_class.csv" \
  --relationships "icd10cm_rels_subclassof.csv" \
  --relationships "icd10cm_rels_classof.csv"

$ ls
icd10cm_nodes_class.csv      icd10cm_rels_classof.csv     icd10cm_load.sh
icd10cm_nodes_vocabulary.csv icd10cm_rels_subclassof.csv

$ bash icd10cm_load.sh  # assuming the paths are correct
```
