package main

import (
	"hash/maphash"
	"reflect"
	"sort"
	"testing"
)

func TestBuildFuncReceivers(t *testing.T) {
	type test struct {
		input []FieldDesc
		want  []string
	}
	tests := []test{
		{
			input: nil,
			want:  []string{""},
		},
		{
			input: []FieldDesc{},
			want:  []string{""},
		},
		{
			input: []FieldDesc{{Name: "a"}, {Name: "b"}, {Name: "c", Optional: true}, {Name: "d", Optional: true}},
			want:  []string{"AB", "ABC", "ABD", "ABCD"},
		},
	}
	for _, tc := range tests {
		got := BuildFuncReceiverSuffixes(StructDesc{Fields: tc.input})
		want := tc.want
		sort.Strings(got)
		sort.Strings(want)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\nwant: %#v\ngot:  %#v", want, got)
		}
	}
}

func TestPartialBuilderTypeSuffixes(t *testing.T) {
	type test struct {
		input []FieldDesc
		want  []string
	}
	tests := []test{
		{
			input: nil,
			want:  nil,
		},
		{
			input: []FieldDesc{},
			want:  nil,
		},
		{
			input: []FieldDesc{{Name: "a"}, {Name: "b"}, {Name: "c", Optional: true}},
			want:  []string{"A", "B", "C", "AB", "AC", "BC", "ABC"},
		},
	}
	for _, tc := range tests {
		got := PartialBuilderTypeSuffixes(StructDesc{Fields: tc.input})
		want := tc.want
		sort.Strings(got)
		sort.Strings(want)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\nwant: %#v\ngot:  %#v", want, got)
		}
	}
}

func TestWithFuncs(t *testing.T) {
	type test struct {
		input []FieldDesc
		want  []WithFuncDesc
	}
	tests := []test{
		{
			input: nil,
			want:  nil,
		},
		{
			input: []FieldDesc{},
			want:  nil,
		},
		{
			input: []FieldDesc{{Name: "a", Type: "aType"}, {Name: "B", Type: "bType"}, {Name: "c", Type: "cType", Optional: true}},
			want: []WithFuncDesc{
				{FieldName: "a", FieldType: "aType", ReceiverSuffix: "", ReturnSuffix: "A"},
				{FieldName: "a", FieldType: "aType", ReceiverSuffix: "B", ReturnSuffix: "AB"},
				{FieldName: "a", FieldType: "aType", ReceiverSuffix: "C", ReturnSuffix: "AC"},
				{FieldName: "a", FieldType: "aType", ReceiverSuffix: "BC", ReturnSuffix: "ABC"},
				{FieldName: "B", FieldType: "bType", ReceiverSuffix: "", ReturnSuffix: "B"},
				{FieldName: "B", FieldType: "bType", ReceiverSuffix: "A", ReturnSuffix: "AB"},
				{FieldName: "B", FieldType: "bType", ReceiverSuffix: "C", ReturnSuffix: "BC"},
				{FieldName: "B", FieldType: "bType", ReceiverSuffix: "AC", ReturnSuffix: "ABC"},
				{FieldName: "c", FieldType: "cType", ReceiverSuffix: "", ReturnSuffix: "C"},
				{FieldName: "c", FieldType: "cType", ReceiverSuffix: "A", ReturnSuffix: "AC"},
				{FieldName: "c", FieldType: "cType", ReceiverSuffix: "B", ReturnSuffix: "BC"},
				{FieldName: "c", FieldType: "cType", ReceiverSuffix: "AB", ReturnSuffix: "ABC"},
			},
		},
	}
	for _, tc := range tests {
		got := WithFuncs(StructDesc{Fields: tc.input})
		want := tc.want
		sortWithFuncDescs(got)
		sortWithFuncDescs(want)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\nwant: %#v\ngot:  %#v", want, got)
		}
	}
}

func TestCombinations(t *testing.T) {
	type test struct {
		input        []FieldDesc
		optionalOnly bool
		want         [][]string
	}
	tests := []test{
		{
			input: nil,
			want:  nil,
		},
		{
			input: []FieldDesc{},
			want:  nil,
		},
		{
			input: []FieldDesc{{Name: "aField"}},
			want:  [][]string{{"AField"}},
		},
		{
			input: []FieldDesc{{Name: "aField"}, {Name: "BField"}},
			want:  [][]string{{"AField"}, {"BField"}, {"AField", "BField"}},
		},
		{
			input: []FieldDesc{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			want:  [][]string{{"A"}, {"B"}, {"C"}, {"A", "B"}, {"A", "C"}, {"B", "C"}, {"A", "B", "C"}},
		},
		{
			input:        []FieldDesc{{Name: "a", Optional: true}, {Name: "b", Optional: true}, {Name: "c"}},
			want:         [][]string{{"A"}, {"B"}, {"A", "B"}},
			optionalOnly: true,
		},
	}
	for _, tc := range tests {
		got := combinations(tc.input, tc.optionalOnly)
		want := tc.want
		sortStringSlices(got)
		sortStringSlices(want)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("\nwant: %#v\ngot:  %#v", want, got)
		}
	}
}

var seed = maphash.MakeSeed()

func sortStringSlices(s [][]string) {
	sort.Slice(s, func(i, j int) bool {
		var ih, jh maphash.Hash
		ih.SetSeed(seed)
		jh.SetSeed(seed)
		for _, x := range s[i] {
			ih.WriteString(x)
		}
		for _, x := range s[j] {
			jh.WriteString(x)
		}
		return ih.Sum64() > jh.Sum64()
	})
}

func sortWithFuncDescs(s []WithFuncDesc) {
	sort.Slice(s, func(i, j int) bool {
		var ih, jh maphash.Hash
		ih.SetSeed(seed)
		jh.SetSeed(seed)
		ih.WriteString(s[i].FieldName + s[i].FieldType + s[i].ReceiverSuffix + s[i].ReturnSuffix)
		jh.WriteString(s[j].FieldName + s[j].FieldType + s[j].ReceiverSuffix + s[j].ReturnSuffix)
		return ih.Sum64() > jh.Sum64()
	})
}
