package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"reflect"
	"strings"
)

const (
	DefaultTagName   = "builder"
	OptionalTagValue = "opt"
	IgnoreTagValue   = "ignore"
)

// StructDesc describes a struct with all the information necessary to generate
// a builder.
type StructDesc struct {
	Name   string
	Fields []FieldDesc
}

// FieldDesc describes a single struct field.
type FieldDesc struct {
	Name, Type string
	Optional   bool
}

// DescribeTaggedStructs produces descriptions of all structs within a file
// containing at least one field with a marker tag.
//
// A field will be considered tagged with a marker tag if the tag is formatted
// using the convention from [reflect], and has a key matching the provided
// tagName. If tagName is empty, [DefaultTagName] will be matched instead.
//
// Currently implemented tag values are:
//   - "opt": mark the field as optional
//   - "ignore": behave as if field does not exist
//
// All other tag values have no effect on the description of the field, but
// will still be considered a marker tag.
func DescribeTaggedStructs(file, tagName string) ([]StructDesc, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return nil, err
	}

	types := listTaggedStructTypes(f, tagName)

	var descs []StructDesc
	for structName, t := range types {
		fields := describeStructFields(fset, t, tagName)
		descs = append(descs, StructDesc{
			Name:   structName,
			Fields: fields,
		})
	}

	return descs, nil
}

func describeStructFields(
	fset *token.FileSet,
	t *ast.StructType,
	tagName string,
) []FieldDesc {
	var fields []FieldDesc

	for _, f := range t.Fields.List {
		tag, _ := tagValue(f.Tag, tagName)
		if tag == IgnoreTagValue {
			continue
		}

		fType := asStr(fset, f.Type)

		for _, n := range f.Names {
			fields = append(fields, FieldDesc{
				Name:     n.Name,
				Type:     fType,
				Optional: tag == OptionalTagValue,
			})
		}
	}

	return fields
}

func listTaggedStructTypes(
	f *ast.File,
	tagName string,
) map[string]*ast.StructType {
	types := make(map[string]*ast.StructType)

	for _, d := range f.Decls {
		tDecl := typeDecl(d)
		if tDecl == nil {
			continue
		}
		for _, spec := range tDecl.Specs {
			t, name := structType(spec)
			if t == nil {
				continue
			}
			tagged := false
			for _, f := range t.Fields.List {
				_, ok := tagValue(f.Tag, tagName)
				if ok {
					tagged = true
					break
				}
			}
			if !tagged {
				continue
			}
			types[name.String()] = t
		}
	}

	return types
}

func typeDecl(d ast.Decl) *ast.GenDecl {
	gd, ok := d.(*ast.GenDecl)
	if !ok {
		return nil
	}

	if gd.Tok != token.TYPE {
		return nil
	}

	return gd
}

func structType(s ast.Spec) (*ast.StructType, *ast.Ident) {
	ts, ok := s.(*ast.TypeSpec)
	if !ok {
		return nil, nil
	}

	t, ok := ts.Type.(*ast.StructType)
	if !ok {
		return nil, nil
	}

	return t, ts.Name
}

func asStr(fset *token.FileSet, n ast.Node) string {
	b := new(bytes.Buffer)
	printer.Fprint(b, fset, n)

	return b.String()
}

func tagValue(lit *ast.BasicLit, tagName string) (v string, ok bool) {
	if lit == nil {
		return "", false
	}
	tag := reflect.StructTag(strings.Trim(lit.Value, "`"))

	if tagName == "" {
		v, ok = tag.Lookup(DefaultTagName)
	} else {
		v, ok = tag.Lookup(tagName)
	}
	return
}
