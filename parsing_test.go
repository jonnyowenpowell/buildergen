package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestDescribeTaggedStructs(t *testing.T) {
	filename := "testdata/threestructs.go"

	want := []StructDesc{
		{
			Name: "X",
			Fields: []FieldDesc{
				{Name: "a", Type: "int"},
				{Name: "B", Type: "int"},
				{Name: "c", Type: "int", Optional: true},
				{Name: "e", Type: "time.Duration"},
				{Name: "f", Type: "[]struct {\n\tl, n, m []int `json:\"l\"`\n}"},
			},
		},
		{
			Name: "y",
			Fields: []FieldDesc{
				{Name: "a", Type: "string"},
				{Name: "b", Type: "string"},
			},
		},
	}

	got, err := DescribeTaggedStructs(filename, "btest")
	if err != nil {
		t.Fatalf("received unexpected error describing %s", filename)
	}
	sort.Slice(got, func(i, j int) bool {
		return sort.StringsAreSorted([]string{got[i].Name, got[j].Name})
	})

	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nwant: %#v\ngot:  %#v", want, got)
	}
}
