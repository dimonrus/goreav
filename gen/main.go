package gen

import (
	"os"
	"fmt"
)

type AppTemplate map[string]interface{}

func ParseTemplate(template AppTemplate) error {
	for key, value := range template {
		if key == KeyWordProject {
			data := value.(map[interface {}]interface {})
			path := fmt.Sprintf("%s/%s", data[KeyWordPath], data[KeyWordName])
			err := os.Mkdir(path, 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
