package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// NOTE: Instead of using this script we are using bazel's go_embed_data.
// This script may be removed in the future.

var (
	in         = flag.String("in", "", "The path to the directory that contains text files.")
	out        = flag.String("out", "", "The path to put the generated go file.")
	outPackage = flag.String("out-package", "", "The go package for generated file.")

	codeTemplate = `
package {{ .Package }}

var (
{{- range $name, $data := .Data }}
	{{ formatName $name }} = {{ formatString $data }}
{{- end }}
)
`
)

func main() {
	flag.Parse()
	if err := validateFlags(); err != nil {
		log.Fatalf("invalidate flag: %v", err)
	}

	license, err := os.ReadFile("./hack/boilerplate.go.txt")
	if err != nil {
		log.Fatalf("unable to read license file: %v", err)
	}

	data, err := loadData(*in)
	if err != nil {
		log.Fatalf("unable to load data from %s: %v", *in, err)
	}

	var maxNameLength = 0
	for k := range data {
		if maxNameLength < len(k) {
			maxNameLength = len(k)
		}
	}
	format := "%-" + fmt.Sprintf("%d", maxNameLength) + "s"
	formatName := func(s string) string {
		return fmt.Sprintf(format, s)
	}

	fileTemplate := string(license) + codeTemplate
	tmpl, err := template.New("template").
		Funcs(template.FuncMap{
			"formatName":   formatName,
			"formatString": formatString,
		}).
		Parse(fileTemplate)
	if err != nil {
		log.Fatalf("unable to make template: %v", err)
	}

	generatedCode, err := renderTemplate(tmpl, map[string]interface{}{
		"Package":    *outPackage,
		"Data":       data,
		"NameLength": 50,
	})
	if err != nil {
		log.Fatalf("unable to render go file: %v", err)
	}

	if err = os.WriteFile(*out, []byte(generatedCode), os.ModePerm); err != nil {
		log.Fatalf("unable to write go file to %s: %v", *out, err)
	}
	log.Printf("successfully generated file: %s", *out)
}

func validateFlags() error {
	if *in == "" {
		return fmt.Errorf("in is required")
	}
	if *out == "" {
		return fmt.Errorf("out is required")
	}
	if *outPackage == "" {
		return fmt.Errorf("out-package is required")
	}
	return nil
}

func loadData(root string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		data[filepath.Base(path)] = bytes
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func formatString(bytes []byte) string {
	return fmt.Sprintf("%#v", string(bytes))
}

func renderTemplate(tmpl *template.Template, data map[string]interface{}) (string, error) {
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
