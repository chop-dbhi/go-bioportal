package main

import (
	"context"
	"os"
	"reflect"
	"testing"

	neolog "github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
)

func init() {
	neolog.SetLevel("error")
}

func TestService(t *testing.T) {
	url := os.Getenv("NEO4J_URL")
	if url == "" {
		t.Skip("NEO4J_URL required")
	}

	s, err := NewService(url, -1)
	if err != nil {
		t.Fatal(err)
	}

	cxt := context.Background()
	vocab := "icd10cm"

	t.Run("match", func(t *testing.T) {
		// Execute.
		strm := make(chan *Class)
		go func() {
			defer close(strm)
			if err := s.Match(cxt, vocab, ".*Trisomy 21.*", strm); err != nil {
				t.Fatal(err)
			}
		}()

		// Consume results.
		var res []*Class
		for c := range strm {
			res = append(res, c)
		}

		if len(res) != 4 {
			t.Errorf("expected 4 results, got %d", len(res))
		}
	})

	t.Run("get", func(t *testing.T) {
		code := "Q90"
		c, err := s.Get(vocab, code)
		if err != nil {
			t.Fatal(err)
		}

		if c.Code != code {
			t.Error("codes don't match")
		}
	})

	t.Run("validate", func(t *testing.T) {
		res, err := s.Validate(vocab, []string{"Q90", "XXX", "W17.81XD"})
		if err != nil {
			t.Fatal(err)
		}

		exp := []bool{true, false, true}

		if !reflect.DeepEqual(res, exp) {
			t.Error("result did not match")
		}
	})

	t.Run("parents", func(t *testing.T) {
		// Execute.
		strm := make(chan *Class)
		go func() {
			defer close(strm)
			if err := s.Parents(cxt, vocab, "Q90", strm); err != nil {
				t.Fatal(err)
			}
		}()

		// Consume results.
		var res []*Class
		for c := range strm {
			res = append(res, c)
		}

		if len(res) != 1 {
			t.Errorf("expected 1 result, got %d", len(res))
		}
	})

	t.Run("children", func(t *testing.T) {
		// Execute.
		strm := make(chan *Class)
		go func() {
			defer close(strm)
			if err := s.Children(cxt, vocab, "Q90", strm); err != nil {
				t.Fatal(err)
			}
		}()

		// Consume results.
		var res []*Class
		for c := range strm {
			res = append(res, c)
		}

		if len(res) != 4 {
			t.Errorf("expected 4 results, got %d", len(res))
		}
	})

	t.Run("ancestors", func(t *testing.T) {
		// Execute.
		strm := make(chan *Class)
		go func() {
			defer close(strm)
			if err := s.Ancestors(cxt, vocab, "Q90", strm); err != nil {
				t.Fatal(err)
			}
		}()

		// Consume results.
		var res []*Class
		for c := range strm {
			res = append(res, c)
		}

		if len(res) != 2 {
			t.Errorf("expected 2 results, got %d", len(res))
		}
	})

	t.Run("descendants", func(t *testing.T) {
		// Execute.
		strm := make(chan *Class)
		go func() {
			defer close(strm)
			if err := s.Descendants(cxt, vocab, "Q90", strm); err != nil {
				t.Fatal(err)
			}
		}()

		// Consume results.
		var res []*Class
		for c := range strm {
			res = append(res, c)
		}

		if len(res) != 4 {
			t.Errorf("expected 4 results, got %d", len(res))
		}
	})

	t.Run("flatten", func(t *testing.T) {
		strm := make(chan *Class)
		go func() {
			defer close(strm)
			if err := s.Flatten(cxt, vocab, []string{"Q90"}, strm); err != nil {
				t.Fatal(err)
			}
		}()

		// Consume results.
		var res []*Class
		for c := range strm {
			res = append(res, c)
		}

		if len(res) != 5 {
			t.Errorf("expected 5 results, got %d", len(res))
		}
	})

}
