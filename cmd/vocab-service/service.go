package main

import "context"

type Class struct {
	ID       string
	Vocab    string
	Code     string
	Label    string
	Synonyms []string
}

type Service interface {
	// Get a class by code.
	Get(vocab, code string) (*Class, error)

	// Checks whether a set of codes exist.
	Validate(vocab string, codes []string) ([]bool, error)

	// Match classes by pattern.
	Match(cxt context.Context, vocab, pattern string, res chan<- *Class) error

	// Parents returns all parents of a class.
	Parents(cxt context.Context, vocab, code string, res chan<- *Class) error

	// Children returns all children of a class.
	Children(cxt context.Context, vocab, code string, res chan<- *Class) error

	// Get all ancestors of this class.
	Ancestors(cxt context.Context, vocab, code string, res chan<- *Class) error

	// Get all descendants of this class.
	Descendants(cxt context.Context, vocab, code string, res chan<- *Class) error

	// Flatten takes a set of codes and returns all the codes themselves with
	// all descedents. The use case if for matching
	Flatten(cxt context.Context, vocab string, codes []string, res chan<- *Class) error
}
