// +build ignore

package main

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("skipmap.go")
	if err != nil {
		panic(err)
	}
	filedata, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	w := new(bytes.Buffer)
	// Step 1. Add file header
	w.WriteString(`// Code generated by go run types_gen.go; DO NOT EDIT.` + "\r\n")
	// Step 2. Add imports and package statement
	w.WriteString(string(filedata)[strings.Index(string(filedata), "package skipmap") : strings.Index(string(filedata), ")\n")+1])
	// Step 3. Generate code for all types
	ts := []string{"Float32", "Float64", "Int32", "Int16", "Int", "Uint64", "Uint32", "Uint16", "Uint"} // all types need to be converted
	for _, upper := range ts {
		data := string(filedata)
		// Step 4-1. Remove all string before import
		data = data[strings.Index(data, ")\n")+1:]
		// Step 4-2. Replace all cases
		dataDesc := replace(data, upper, true)
		dataAsc := replace(data, upper, false)
		w.WriteString(dataAsc)
		w.WriteString("\r\n")
		w.WriteString(dataDesc)
		w.WriteString("\r\n")
	}

	out, err := format.Source(w.Bytes())
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("types.go", out, 0660); err != nil {
		panic(err)
	}
}

func replace(data string, upper string, desc bool) string {
	lower := strings.ToLower(upper)

	var descstr string
	if desc {
		descstr = "Desc"
	}
	data = strings.Replace(data, "NewInt64", "New"+upper+descstr, -1)
	data = strings.Replace(data, "newInt64Node", "new"+upper+"Node"+descstr, -1)
	data = strings.Replace(data, "unlockInt64", "unlock"+upper+descstr, -1)
	data = strings.Replace(data, "Int64Map", upper+"Map"+descstr, -1)
	data = strings.Replace(data, "int64Node", lower+"Node"+descstr, -1)
	data = strings.Replace(data, "key int64", "key "+lower, -1)
	data = strings.Replace(data, "key  int64", "key  "+lower, -1)
	data = strings.Replace(data, "key   int64", "key   "+lower, -1)
	data = strings.Replace(data, "int64 skipmap", lower+" skipmap", -1) // comment

	if desc {
		// Special cases for DESC.
		data = strings.Replace(data, "ascending", "descending", -1)
		data = strings.Replace(data, "return n.key < key", "return n.key > key", -1)
	}
	return data
}

func lowerSlice(s []string) []string {
	n := make([]string, len(s))
	for i, v := range s {
		n[i] = strings.ToLower(v)
	}
	return n
}

func inSlice(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}
