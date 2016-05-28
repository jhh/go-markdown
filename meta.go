package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v2"
)

const delim = "---"

// Meta is the metadata for a post
type Meta struct {
	Title      string
	Date       time.Time
	Categories []string
	Tags       []string
}

// MetaError is an error occurred during metadata parsing.
type MetaError struct {
	error
}

func (e MetaError) Error() string {
	return fmt.Sprintf("metadata: parsing error: %v", e.error)
}

// NewMeta initializes metadata from YAML header.
func NewMeta(r io.Reader) (Meta, error) {
	s := bufio.NewScanner(r)

	// scan starting delimeter
	if ok := s.Scan(); !ok {
		return Meta{}, MetaError{s.Err()}
	}
	if s.Text() != delim {
		return Meta{}, MetaError{fmt.Errorf("%q not valid delimeter", s.Text())}
	}

	// scan metadata block
	var b []byte
	for s.Scan() {
		if s.Text() == delim {
			m := Meta{}
			err := yaml.Unmarshal(b, &m)
			if err != nil {
				return m, MetaError{err}
			}
			return m, nil
		}
		b = append(append(b, s.Bytes()...), '\n')
	}
	if err := s.Err(); err != nil {
		return Meta{}, MetaError{err}
	}

	// end of file
	return Meta{}, MetaError{errors.New("no closing delimeter")}
}
