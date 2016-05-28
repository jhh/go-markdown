package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var tests = []struct {
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
	tests[0].meta.Date = d
	tests[2].meta.Date = d
	d, _ = time.Parse(time.RFC3339, "2016-05-28T17:32:43.000Z")
	tests[1].meta.Date = d
}

func checkError(t *testing.T, err error, want bool) {
	got := err != nil
	if got != want {
		t.Errorf("is error got: %v, want: %v\n\t err: %v", got, want, err)
	}
}

func TestCases(t *testing.T) {
	for _, tc := range tests {
		t.Log("test case: " + tc.name)
		f, err := os.Open(filepath.Join("testdata", tc.file))
		if err != nil {
			t.Error(err)
		}
		m, err := NewMeta(f)
		checkError(t, err, tc.isErr)
		if !reflect.DeepEqual(m, tc.meta) {
			t.Errorf("got: %v, want:%v", m, tc.meta)
		}
		f.Close()
	}
}
