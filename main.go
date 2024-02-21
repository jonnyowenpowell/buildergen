package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/tools/imports"
)

func main() {
	file, ok := os.LookupEnv("GOFILE")
	if !ok {
		log.Fatal("$GOFILE must be set")
	}

	pkg, ok := os.LookupEnv("GOPACKAGE")
	if !ok {
		log.Fatal("$GOPACKAGE must be set")
	}

	builderTag := flag.String("tag", "",
		"set the struct tag key used to control builder generation")
	flag.Parse()

	Generate(file, pkg, *builderTag)
}

// Generate parses the provided file and generates one builder per tagged
// struct. The documentation for [DescribeTaggedStructs] details which structs
// are considered tagged.
//
// Generated builders are created in the provided package, which must match the
// package of file. Each builder is output to a separate file, named
// {file}_{struct}builder.gen.go, where struct is the lowercased name of the
// struct from the provided file.
func Generate(file, pkg, builderTag string) {
	descs, err := DescribeTaggedStructs(file, builderTag)
	if err != nil {
		log.Fatalf("failed to parse source file: %v\n", err)
	}

	var errs []error
	for _, desc := range descs {
		buf := new(bytes.Buffer)
		err := RenderBuilder(buf, pkg, desc)
		if err != nil {
			errs = append(errs,
				fmt.Errorf("failed to render builder for %s: %v", desc.Name, err))
			continue
		}
		out := buf.Bytes()

		out, err = imports.Process(file, out, nil)
		if err != nil {
			errs = append(errs,
				fmt.Errorf("failed to run goimports on builder source for %s: %v",
					desc.Name, err))
			continue
		}

		outName := fmt.Sprintf("%s_%sbuilder.gen.go", strings.TrimRight(file, ".go"),
			strings.ToLower(desc.Name))
		outFile, err := os.Create(outName)
		if err != nil {
			errs = append(errs,
				fmt.Errorf("failed to open file %s to write builder for %s: %v", outName,
					desc.Name, err))
			continue
		}

		defer outFile.Close()
		outFile.Write(out)
	}

	if len(errs) > 0 {
		log.Fatalf("failed to generate builders: %v\n", errors.Join(errs...))
	}
}
