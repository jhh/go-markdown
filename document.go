package markdown

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/russross/blackfriday"

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

// NewDocument creates a new document from filename.
func NewDocument(filename string) (doc *Document, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return &Document{}, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// read document by lines
	s := bufio.NewScanner(f)

	// read front-matter metadata
	meta, err := readMeta(s)
	if err != nil {
		return &Document{}, err
	}
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
