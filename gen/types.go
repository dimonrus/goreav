package gen

import (
	"fmt"
	"strings"
	"unicode"
)

func getNestedSpaces(level int) string {
	var spaces string
	for i := 0; i < level; i++ {
		spaces += "    "
	}
	return spaces
}

func getYamlTag(key string) string {
	for i, value := range key {
		if i == 0 {
			continue
		}
		if unicode.IsUpper(value) {
			return fmt.Sprintf(" `yaml:\"%s\"`", key)
		}
	}
	return ""
}

func createNestedStructure(data AppTemplate, nestedLevel int) string {
	var str string
	for key, value := range data {
		switch value.(type) {
		case AppTemplate:
			//fmt spaces filedName fieldType yaml tag
			str += fmt.Sprintf("%s%s %s%s\n", getNestedSpaces(nestedLevel), strings.Title(key.(string)), createNestedStructure(value.(AppTemplate), nestedLevel+1), getYamlTag(key.(string)))
		default:
			str += fmt.Sprintf("%s%s %T%s\n", getNestedSpaces(nestedLevel), strings.Title(key.(string)), value, getYamlTag(key.(string)))
		}
	}
	return fmt.Sprintf("struct {\n%s%s}", str, getNestedSpaces(nestedLevel-1))
}

func CreateTypeStructure(data AppTemplate, name string) (string, error) {
	head := fmt.Sprintf("type %s", name)
	var str string
	for key, value := range data {
		switch value.(type) {
		case AppTemplate:
			//fmt spaces filedName fieldType
			str += fmt.Sprintf("%s%s %s\n", getNestedSpaces(1), strings.Title(key.(string)), createNestedStructure(value.(AppTemplate), 2))
		default:
			str += fmt.Sprintf("%s%s %T%s\n", getNestedSpaces(1), strings.Title(key.(string)), value, getYamlTag(key.(string)))
		}
	}
	return fmt.Sprintf("%s struct {\n%s}\n", head, str), nil
}
