// Package markdown translates plain text with Markdown formatting and YAML
// front-matter into HTML documents with metadata.
// Markdown processing is done by "github.com/russross/blackfriday".
// Document format is:
//
// 	---
// 	title: Complete
// 	date: 2016-05-27T23:45:46.000Z # RFC3339 format
// 	categories:
//   	  - test
// 	tags:
//   	  - foo
//   	  - bar
// 	---
// 	# Heading
// 	The remainder of the document in *Markdown* format.
package markdown // import "jhhgo.us/markdown"

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/russross/blackfriday"

	"gopkg.in/yaml.v2"
)

const delim = "---"

// Meta is the metadata for a post, taken from the document's YAML front-matter.
// RelPath is set to the path (relative to basepath in NewDocument) of the
// parsed document and has '/' path separators. For example, if
// NewDocument("/a", "/a/b/document.md") is called, Meta.RelPath will be
// "b/document.md"
type Meta struct {
	Title      string
	Date       time.Time
	Categories []string
	Tags       []string
	RelPath    string
}

// MetaError is an error occurred during metadata parsing.
type MetaError struct {
	error
}

func (e MetaError) Error() string {
	return fmt.Sprintf("metadata: parsing error: %v", e.error)
}

func readMeta(s *bufio.Scanner) (Meta, error) {

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

// A Document represents a Markdown document with metadata front-matter.
type Document struct {
	Meta
	Body []byte
}

// DocumentError is an error occurred during document parsing.
type DocumentError struct {
	error
}

func (e DocumentError) Error() string {
	return fmt.Sprintf("document: parsing error: %v", e.error)
}

// NewDocument creates a new document from filename. `Meta.RelPath` is set to the
// file path relative to `basepath`. See Meta documentation for example.
func NewDocument(basepath, filename string) (doc *Document, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return &Document{}, DocumentError{err}
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = DocumentError{cerr}
		}
	}()

	// read document by lines
	s := bufio.NewScanner(f)

	// read front-matter metadata
	meta, err := readMeta(s)
	if err != nil {
		return &Document{}, DocumentError{err}
	}
	rel, err := filepath.Rel(basepath, filename)
	if err != nil {
		return &Document{}, DocumentError{err}
	}
	meta.RelPath = rel
	doc = &Document{}
	doc.Meta = meta

	// read body
	var b []byte
	for s.Scan() {
		b = append(append(b, s.Bytes()...), '\n')
	}
	if err := s.Err(); err != nil {
		return &Document{}, DocumentError{err}
	}
	doc.Body = blackfriday.MarkdownCommon(b)
	return doc, nil
}
