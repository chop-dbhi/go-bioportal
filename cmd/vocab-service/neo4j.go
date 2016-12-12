package main

import (
	"context"
	"errors"
	"io"
	"strings"

	neo "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

var DefaultMaxConn = 20

type service struct {
	url  string
	pool neo.DriverPool
}

func (s *service) Match(cxt context.Context, vocab, pattern string, res chan<- *Class) error {
	if vocab == "" {
		return errors.New("vocab cannot be empty")
	}

	if pattern == "" {
		return errors.New("pattern cannot be empty")
	}

	// Case insensitive.
	pattern = strings.ToLower(pattern)

	conn, err := s.pool.OpenPool()
	if err != nil {
		return err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(`
		MATCH (v:Vocabulary {id: {vocab}})<-[:classOf]-(c:Class)
		WHERE lower(c.label) =~ {pattern}
			OR any(syn in c.synonyms where lower(syn) =~ {pattern})
		RETURN c.id, v.id, c.label, c.code, c.synonyms
	`)
	if err != nil {
		return err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"pattern": pattern,
		"vocab":   vocab,
	})
	if err != nil {
		return err
	}
	defer rows.Close()

	for {
		vals, _, err := rows.NextNeo()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		c := &Class{
			ID:    vals[0].(string),
			Vocab: vals[1].(string),
			Label: vals[2].(string),
			Code:  vals[3].(string),
		}

		if vals[4] != nil {
			for _, v := range vals[4].([]interface{}) {
				c.Synonyms = append(c.Synonyms, v.(string))
			}
		}

		select {
		case <-cxt.Done():
			return nil

		case res <- c:
		}
	}

	return nil
}

// Get a class by code.
func (s *service) Get(vocab, code string) (*Class, error) {
	if vocab == "" {
		return nil, errors.New("vocab cannot be empty")
	}

	if code == "" {
		return nil, errors.New("code cannot be empty")
	}

	conn, err := s.pool.OpenPool()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(`
		MATCH (v:Vocabulary {id: {vocab}})<-[:classOf]-(c:Class {code: {code}})
		RETURN c.id, v.id, c.label, c.code, c.synonyms
		LIMIT 1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"vocab": vocab,
		"code":  code,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vals, _, err := rows.NextNeo()
	if err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	c := &Class{
		ID:    vals[0].(string),
		Vocab: vals[1].(string),
		Label: vals[2].(string),
		Code:  vals[3].(string),
	}

	if vals[4] != nil {
		for _, v := range vals[4].([]interface{}) {
			c.Synonyms = append(c.Synonyms, v.(string))
		}
	}

	return c, nil
}

// Checks whether a set of codes exist.
func (s *service) Validate(vocab string, codes []string) ([]bool, error) {
	if vocab == "" {
		return nil, errors.New("vocab cannot be empty")
	}

	if len(codes) == 0 {
		return nil, nil
	}

	conn, err := s.pool.OpenPool()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(`
		MATCH (:Vocabulary {id: {vocab}})<-[:classOf]-(c:Class {code: {code}})
		RETURN 1
		LIMIT 1
	`)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"code":  "",
		"vocab": vocab,
	}

	bools := make([]bool, len(codes))

	for i, code := range codes {
		params["code"] = code

		rows, err := stmt.QueryNeo(params)
		if err != nil {
			return nil, err
		}

		_, _, err = rows.NextNeo()
		rows.Close()

		if err == io.EOF {
			continue
		} else if err != nil {
			return nil, err
		}

		bools[i] = true
	}

	return bools, nil
}

func (s *service) traversal(cxt context.Context, query, vocab, code string, res chan<- *Class) error {
	if vocab == "" {
		return errors.New("vocab cannot be empty")
	}

	if code == "" {
		return errors.New("code must be specified")
	}

	conn, err := s.pool.OpenPool()
	if err != nil {
		return err
	}
	defer conn.Close()

	stmt, err := conn.PrepareNeo(query)
	if err != nil {
		return err
	}

	rows, err := stmt.QueryNeo(map[string]interface{}{
		"code":  code,
		"vocab": vocab,
	})
	if err != nil {
		return err
	}
	defer rows.Close()

	for {
		vals, _, err := rows.NextNeo()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		c := &Class{
			ID:    vals[0].(string),
			Vocab: vals[1].(string),
			Label: vals[2].(string),
			Code:  vals[3].(string),
		}

		if vals[4] != nil {
			for _, v := range vals[4].([]interface{}) {
				c.Synonyms = append(c.Synonyms, v.(string))
			}
		}

		select {
		case <-cxt.Done():
			return nil

		case res <- c:
		}
	}

	return nil
}

// Parents returns all parents of a class.
func (s *service) Parents(cxt context.Context, vocab, code string, res chan<- *Class) error {
	query := `
		MATCH (v:Vocabulary {id: {vocab}})<-[:classOf]-(:Class {code: {code}})-[:subClassOf]->(c:Class)
		RETURN c.id, v.id, c.label, c.code, c.synonyms
	`
	return s.traversal(cxt, query, vocab, code, res)
}

// Children returns all children of a class.
func (s *service) Children(cxt context.Context, vocab, code string, res chan<- *Class) error {
	query := `
		MATCH (v:Vocabulary {id: {vocab}})<-[:classOf]-(:Class {code: {code}})<-[:subClassOf]-(c:Class)
		RETURN c.id, v.id, c.label, c.code, c.synonyms
	`
	return s.traversal(cxt, query, vocab, code, res)
}

// Get all ancestors of this class.
func (s *service) Ancestors(cxt context.Context, vocab, code string, res chan<- *Class) error {
	query := `
		MATCH (v:Vocabulary {id: {vocab}})<-[:classOf]-(:Class {code: {code}})-[:subClassOf*1..]->(c:Class)
		RETURN c.id, v.id, c.label, c.code, c.synonyms
	`
	return s.traversal(cxt, query, vocab, code, res)
}

// Get all descendants of this class.
func (s *service) Descendants(cxt context.Context, vocab, code string, res chan<- *Class) error {
	query := `
		MATCH (v:Vocabulary {id: {vocab}})<-[:classOf]-(:Class {code: {code}})<-[:subClassOf*1..]-(c:Class)
		RETURN c.id, v.id, c.label, c.code, c.synonyms
	`
	return s.traversal(cxt, query, vocab, code, res)
}

// Flatten takes a set of codes and returns all the codes themselves with
// all descendants. The use case if for matching
func (s *service) Flatten(cxt context.Context, vocab string, codes []string, res chan<- *Class) error {
	query := `
			MATCH (v:Vocabulary {id: {vocab}})<-[:classOf]-(:Class {code: {code}})<-[:subClassOf*0..]-(c:Class)
			RETURN c.id, v.id, c.label, c.code, c.synonyms
		`
	for _, code := range codes {
		if err := s.traversal(cxt, query, vocab, code, res); err != nil {
			return err
		}
	}

	return nil
}

func NewService(url string, maxconn int) (Service, error) {
	if maxconn <= 0 {
		maxconn = DefaultMaxConn
	}

	if !strings.HasPrefix(url, "bolt://") {
		url = "bolt://" + url
	}

	pool, err := neo.NewDriverPool(url, maxconn)
	if err != nil {
		return nil, err
	}

	return &service{
		url:  url,
		pool: pool,
	}, nil
}
