package gen

import (
	"fmt"
	"strings"
	"unicode"
)

func createNestedStructure(data AppTemplate, nestedLevel int) string {
	var str string
	for key, value := range data {
		switch value.(type) {
		case AppTemplate:
			str += fmt.Sprintf("%s%s %s", getNestedSpaces(nestedLevel), strings.Title(key.(string)), createNestedStructure(value.(AppTemplate), nestedLevel+1))
		default:
			str += fmt.Sprintf("%s%s %T %s\n", getNestedSpaces(nestedLevel), strings.Title(key.(string)), value, getYamlTag(key.(string)))
		}
	}

	return fmt.Sprintf("struct {\n%s%s}\n", str, getNestedSpaces(nestedLevel-1))
}

func getNestedSpaces(level int) string {
	var spaces string
	for i := 0; i < level; i++ {
		spaces += "    "
	}
	return spaces
}

func getYamlTag(key string) string {
	for _, value := range key {
		if unicode.IsUpper(value) {
			return fmt.Sprintf("`yaml:\"%s\"`", key)
		}
	}
	return ""
}

func CreateTypeStructure(data AppTemplate, name string) (string, error) {
	head := fmt.Sprintf("type %s", name)
	var str string
	for key, value := range data {
		switch value.(type) {
		case AppTemplate:
			str += fmt.Sprintf("%s%s %s", getNestedSpaces(1), strings.Title(key.(string)), createNestedStructure(value.(AppTemplate), 2))
		default:
			str += fmt.Sprintf("%s%s %T %s\n", getNestedSpaces(1), strings.Title(key.(string)), value, getYamlTag(key.(string)))
		}

	}
	return fmt.Sprintf("%s struct {\n%s}\n", head, str), nil
}
