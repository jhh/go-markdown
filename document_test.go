package markdown // import "jhhgo.us/markdown"

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var docTests = []struct {
	name  string
	file  string
	doc   Document
	isErr bool
}{
	{
		"valid complete",
		"ok_complete.md",
		Document{
			Meta: Meta{
				Title:      "Complete",
				Categories: []string{"test"},
				Tags:       []string{"foo", "bar"},
				File:       "ok_complete.md",
			},
			Body: []byte("<h1>Complete</h1>\n"),
		},
		false,
	},
}

func init() {
	d, _ := time.Parse(time.RFC3339, "2016-05-27T23:45:46.000Z")
	docTests[0].doc.Meta.Date = d
}

func TestDocument(t *testing.T) {
	for _, tc := range docTests {
		t.Log("test case: " + tc.name)
		doc, err := NewDocument(filepath.Join("testdata", tc.file))
		checkError(t, err, tc.isErr)
		if !reflect.DeepEqual((*doc).Meta, tc.doc.Meta) {
			t.Errorf("metadata:\ngot: %v\nwant:%v", (*doc).Meta, tc.doc.Meta)
		}
		if !reflect.DeepEqual((*doc).Body, tc.doc.Body) {
			t.Errorf("body:\ngot: %s\nwant:%s", (*doc).Body, tc.doc.Body)
		}
	}
}
