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

func CreateTypeStructure(data AppTemplate, name string, nestedLevel int) string {
	var str string
	for key, value := range data {
		stringKey := key.(string)
		switch value.(type) {
		case AppTemplate:
			//fmt spaces filedName struct
			str += fmt.Sprintf("%s%s %s\n", getNestedSpaces(nestedLevel+1), strings.Title(stringKey), CreateTypeStructure(value.(AppTemplate), stringKey, nestedLevel+1))
		default:
			//fmt spaces filedName fieldType
			str += fmt.Sprintf("%s%s %T%s\n", getNestedSpaces(nestedLevel+1), strings.Title(stringKey), value, getYamlTag(stringKey))
		}
	}
	if nestedLevel == 0 {
		return fmt.Sprintf("type %s struct {\n%s}\n", name, str)
	} else {
		return fmt.Sprintf("struct {\n%s%s}%s", str, getNestedSpaces(nestedLevel), getYamlTag(name))
	}
}