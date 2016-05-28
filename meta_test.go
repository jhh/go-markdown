package main

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var metaTests = []struct {
	name  string
	file  string
	meta  Meta
	isErr bool
}{
	{
		"valid complete",
		"ok_complete.md",
		Meta{
			Title:      "Complete",
			Categories: []string{"test"},
			Tags:       []string{"foo", "bar"},
		},
		false,
	},
	{
		"valid no tags",
		"ok_no_tags.md",
		Meta{
			Title:      "No Tags",
			Categories: []string{"test"},
		},
		false,
	},
	{
		"valid extra metadata",
		"ok_extra.md",
		Meta{
			Title:      "Has Extra Metadata",
			Categories: []string{"test"},
			Tags:       []string{"foo", "bar"},
		},
		false,
	},
	{"open delim invalid", "err_open_del.md", Meta{}, true},
	{"no closing delimeter", "err_no_close_del.md", Meta{}, true},
	{"no metadata", "err_no_metadata.md", Meta{}, true},
}

func init() {
	d, _ := time.Parse(time.RFC3339, "2016-05-27T23:45:46.000Z")
	metaTests[0].meta.Date = d
	metaTests[2].meta.Date = d
	d, _ = time.Parse(time.RFC3339, "2016-05-28T17:32:43.000Z")
	metaTests[1].meta.Date = d
}

func checkError(t *testing.T, err error, want bool) {
	got := err != nil
	if got != want {
		t.Errorf("is error got: %v, want: %v\n\t err: %v", got, want, err)
	}
}

func TestMeta(t *testing.T) {
	for _, tc := range metaTests {
		t.Log("test case: " + tc.name)
		doc, err := NewDocument(filepath.Join("testdata", tc.file))
		checkError(t, err, tc.isErr)
		if !reflect.DeepEqual(doc.Meta, tc.meta) {
			t.Errorf("got: %v, want:%v", doc.Meta, tc.meta)
		}
	}
}
