# Vocab Service

This is a simple service for querying vocabulary classes extracted from the Bioportal. Since many vocabularies are hierarchical, Neo4j is used as the database to store and query the classes.

## Install

Requires Go to install, however if you want pre-built binaries [please request it](https://github.com/chop-dbhi/go-bioportal/issues/)!

```
go get github.com/chop-dbhi/go-bioportal/cmd/vocab-service
```

## Usage

It requires a connection to a Neo4j server through the Bolt interface. The address defaults to the default Bolt address `bolt://127.0.0.1:7687`. If authentication is required, the `<user>:<password>` parts of the URI must be supplied.

```
vocab-service [-http] [-neo4j] [-tlscert] [-tlskey]
```

## Methods

Below are list of the methods the service supports.

- `Get(vocab, class)` - Get the details of a class in the target vocabulary.
- `Validate(vocab, classes)`- Validate the set of classes are in the vocabulary.
- `Match(vocab, pattern)` - Perform a pattern match on the label of the class.
- `Parents(vocab, class)` - Get the parents of the class.
- `Children(vocab, class)` - Get the children of the class.
- `Ancestors(vocab, class)` - Get all ancestors of the class.
- `Descendants(vocab, class)` - Get all descendants of the class.
- `Flatten(vocab, classes)` - Flatten the hierarchy for all classes passed.

## Interfaces

### HTTP

A request is a POST with a JSON-encoded body with the following structure. The keys of the the `params` entry correspond to the parameter names in the method signatures.

```
{
  "method": "match",
  "params": {
    "vocab": "icd10cm",
    "pattern": "down syndrome.*"
  }
}
```
