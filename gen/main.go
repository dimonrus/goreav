package gen

import (
	"fmt"
	"errors"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

var path string //full path to project

var transactions AppTransactionStack

func CreateProject(template AppTemplate) error {
	//Project section is required
	if project, ok := template[KeyWordProject]; ok == true {
		data := project.(AppTemplate)
		path = fmt.Sprintf("%s/%s", data[KeyWordPath], data[KeyWordName])
		transactions = append(transactions, &AppTransactionCreateDir{Path: path, Mode: 0755})
	} else {
		return errors.New("template has no project section")
	}

	return nil
}

func RenderConfig(template AppTemplate) error {
	//environment section is not required
	if environment, ok := template[KeyWordEnvironment]; ok == true {
		configPath := path + "/config"
		yamlConfigPath := configPath + "/yaml"
		//Create config dirs
		transactions = append(transactions, &AppTransactionCreateDir{Path: configPath, Mode: 0755})
		transactions = append(transactions, &AppTransactionCreateDir{Path: yamlConfigPath, Mode: 0755})

		//Create project config struct
		var wholeTemplate = make(AppTemplate)

		//Create config files
		for key, conf := range environment.(AppTemplate) {
			wholeTemplate.Merge(conf.(AppTemplate))
			if conf == nil {
				continue
			}
			filePath := yamlConfigPath + "/" + key.(string) + ".yaml"
			transactions = append(transactions, &AppTransactionCreateFile{Path: filePath})
			data, err := yaml.Marshal(conf)
			if err != nil {
				return err
			}
			transactions = append(transactions, &AppTransactionAppendFile{Path: filePath, Data: data})
		}

		//Render map[interface{}]interface{} to string
		str := CreateTypeStructure(wholeTemplate, "Config", 0)

		configFilePath := configPath + "/" + KeyWordSettings + ".go"
		transactions = append(transactions, &AppTransactionCreateFile{Path: configFilePath})

		str = fmt.Sprintf("package %s\n\n%s", KeyWordSettings, str)
		transactions = append(transactions, &AppTransactionAppendFile{Path: configFilePath, Data: []byte(str)})

		envPath := configPath + "/environment.go"
		templatePath, err := filepath.Abs("")
		if err != nil {
			return err
		}
		templatePath = templatePath + "/gen/tml/environment.tml"
		transactions = append(transactions, &AppTransactionCreateEnvironmentFile{Path: envPath, TemplatePath:templatePath})
	}

	return nil
}

//Function performs parse of map[string]interface and populate transaction stack
//After that all transaction executed by order from 0 to n
func ParseTemplate(template AppTemplate) error {
	//Create project
	if err := CreateProject(template); err != nil {
		return err
	}

	//Render config
	if err := RenderConfig(template); err != nil {
		return err
	}

	//Exec all transaction
	return ExecTransactions(transactions)
}
