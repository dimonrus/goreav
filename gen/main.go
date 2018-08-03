package gen

import (
	"os"
	"fmt"
	"errors"
)

type AppTemplate map[string]interface{}

var path string //full path to project

func CreateProject(template AppTemplate) error {
	//Project section is required
	if project, ok := template[KeyWordProject]; ok == true {
		data := project.(map[interface{}]interface{})
		path = fmt.Sprintf("%s/%s", data[KeyWordPath], data[KeyWordName])
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	} else {
		return errors.New("template has no project section")
	}

	return nil
}

func RenderConfig(template AppTemplate) error {
	//Env section is not required
	if environment, ok := template[KeyWordEnvironment]; ok == true {
		configPath := path + "/config"
		err := os.Mkdir(configPath, 0755)
		if err != nil {
			return err
		}

		//Create config files
		for key, _ := range environment.(map[interface{}]interface{}) {
			filePath := configPath + "/" + key.(string) + ".yaml"
			_, err := os.Create(filePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ParseTemplate(template AppTemplate) error {
	//Create project
	if err := CreateProject(template); err != nil {
		return err
	}

	//Render config
	if err := RenderConfig(template); err != nil {
		return err
	}
	
	return nil
}
