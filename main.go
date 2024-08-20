package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/iancoleman/strcase"
)

// Field represents a single field in the Go struct.
type Field[T any] struct {
	Name string
	Type string
	Tag  string
}

// Struct represents a Go struct.
type Struct[T any] struct {
	Name   string
	Fields []Field[T]
}

// GenerateStruct generates a Go struct from a JSON object.
func GenerateStruct(structName string, jsonObject map[string]interface{}) Struct[string] {
	fields := make([]Field[string], 0, len(jsonObject))

	// キーでソートしてフィールドの順序を一定に保つ
	keys := make([]string, 0, len(jsonObject))
	for k := range jsonObject {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := jsonObject[key]
		fieldName := strcase.ToCamel(key)
		fieldType := getType(value)
		tag := fmt.Sprintf("`json:\"%s\"`", key)

		fields = append(fields, Field[string]{
			Name: fieldName,
			Type: fieldType,
			Tag:  tag,
		})
	}

	return Struct[string]{
		Name:   structName,
		Fields: fields,
	}
}

func getType(value interface{}) string {
	switch value.(type) {
	case float64:
		return "float64"
	case string:
		return "string"
	case bool:
		return "bool"
	case map[string]interface{}:
		return "struct"
	default:
		return "interface{}"
	}
}

// ToCamelCase converts a snake_case or kebab-case string to CamelCase.
func ToCamelCase(input string) string {
	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(input, "_", " ")), " ", "")
}

// guessType returns a string representing the Go type of the value.
func guessType[T any](value T) string {
	switch any(value).(type) {
	case string:
		return "string"
	case float64:
		return "float64"
	case bool:
		return "bool"
	case map[string]interface{}:
		return "map[string]interface{}"
	case []interface{}:
		return "[]interface{}"
	case int, int8, int16, int32, int64:
		return "int"
	case uint, uint8, uint16, uint32, uint64:
		return "uint"
	default:
		return "interface{}"
	}
}

// PrintStruct prints the struct in a Go source file format.
func PrintStruct[T any](s Struct[T]) {
	fmt.Printf("type %s struct {\n", s.Name)
	for _, field := range s.Fields {
		fmt.Printf("\t%s %s %s\n", field.Name, field.Type, field.Tag)
	}
	fmt.Println("}")
}

func main() {
	// Define command line arguments.
	jsonFile := flag.String("json", "", "Path to the JSON file.")
	structName := flag.String("name", "AutoGeneratedStruct", "Name of the Go struct.")
	outputFile := flag.String("output", "", "Path to the output file (default is stdout).")

	flag.Parse()

	if *jsonFile == "" {
		log.Fatal("Please provide a JSON file using the -json flag.")
	}

	// Read the JSON file.
	data, err := os.ReadFile(*jsonFile)
	if err != nil {
		log.Fatalf("Failed to read the JSON file: %v", err)
	}

	// Parse the JSON data.
	var jsonObject map[string]interface{}
	if err := json.Unmarshal(data, &jsonObject); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Generate the Go struct.
	generatedStruct := GenerateStruct(*structName, jsonObject)

	// Print or write the struct.
	if *outputFile != "" {
		f, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Failed to create the output file: %v", err)
		}
		defer f.Close()
		_, err = f.WriteString(fmt.Sprintf("package main\n\n%s", formatStruct(generatedStruct)))
		if err != nil {
			log.Fatalf("Failed to write to the output file: %v", err)
		}
	} else {
		PrintStruct(generatedStruct)
	}
}

func formatStruct[T any](s Struct[T]) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("type %s struct {\n", s.Name))
	for _, field := range s.Fields {
		builder.WriteString(fmt.Sprintf("\t%s %s %s\n", field.Name, field.Type, field.Tag))
	}
	builder.WriteString("}\n")
	return builder.String()
}
