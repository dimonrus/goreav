package gen

import (
	"fmt"
	"errors"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

const AppTemplateFileMode = 0755

var (
	ProjectPath  string
	TemplatePath string
)

var transactions AppTransactionStack

func CreateProject(template AppTemplate) error {
	var err error
	//Project section is required
	if project, ok := template[KeyWordProject]; ok == true {
		TemplatePath, err = filepath.Abs("")
		if err != nil {
			return err
		}
		TemplatePath += "/gen/tml"
		data := project.(AppTemplate)
		ProjectPath = fmt.Sprintf("%s/%s", data[KeyWordPath], data[KeyWordName])
		transactions = append(transactions, &AppTransactionCreateDir{Path: ProjectPath, Mode: AppTemplateFileMode})
	} else {
		return errors.New("template has no project section")
	}

	return nil
}

func RenderConfig(template AppTemplate) error {
	//environment section is not required
	if environment, ok := template[KeyWordEnvironment]; ok == true {
		configPath := ProjectPath + "/config"
		yamlConfigPath := configPath + "/yaml"
		//Create config dirs
		transactions = append(transactions, &AppTransactionCreateDir{Path: configPath, Mode: AppTemplateFileMode})
		transactions = append(transactions, &AppTransactionCreateDir{Path: yamlConfigPath, Mode: AppTemplateFileMode})

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
		str := CreateTypeStructure(wholeTemplate, KeyWordSettings, 0)

		configFilePath := configPath + "/" + KeyWordConfig + ".go"
		transactions = append(transactions, &AppTransactionCreateFile{Path: configFilePath})

		str = fmt.Sprintf("package %s\n\n%s", KeyWordConfig, str)
		transactions = append(transactions, &AppTransactionAppendFile{Path: configFilePath, Data: []byte(str)})

		envPath := configPath + "/environment.go"
		envTemplatePath := TemplatePath + "/config/environment.tml"

		templateData := struct {
			Package    string
			ConfigType string
		}{Package: "config", ConfigType: "settings"}

		transactions = append(transactions, &AppTransactionCreateFileFromTemplate{
			Path:         envPath,
			TemplatePath: envTemplatePath,
			Data:         templateData,
		})
	}

	return nil
}

func CreateMainFile(template AppTemplate) error {
	if serve, ok := template[KeyWordServe]; ok == true {
		project := template[KeyWordProject].(AppTemplate)
		mainPath := ProjectPath + "/main.go"
		templatePath := TemplatePath + "/main.tml"
		templateData := struct {
			Project string
			Apps AppTemplate
		}{
			Project: project[KeyWordName].(string),
			Apps: serve.(AppTemplate),
		}

		transactions = append(transactions, &AppTransactionCreateFileFromTemplate{
			Path: mainPath,
			TemplatePath: templatePath,
			Data:         templateData,
		})
	}
	return nil
}

func RenderLogger(template AppTemplate) error {
	helperDir := ProjectPath + "/helper"
	transactions = append(transactions, &AppTransactionCreateDir{Path: helperDir, Mode: AppTemplateFileMode})

	loggingDir := helperDir + "/logging"
	transactions = append(transactions, &AppTransactionCreateDir{Path: loggingDir, Mode: AppTemplateFileMode})

	loggerPath := loggingDir + "/logger.go"
	templatePath := TemplatePath + "/helper/logging/logger.tml"
	templateData := struct {
		Loggers []string
	}{Loggers: []string{"Query"}}

	transactions = append(transactions, &AppTransactionCreateFileFromTemplate{
		Path:         loggerPath,
		TemplatePath: templatePath,
		Data:         templateData,
	})

	return nil
}

func FormatProject() error {
	transactions = append(transactions, &AppTransactionFormatProject{Path: ProjectPath + "/..."})
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

	//Render logger
	if err := RenderLogger(template); err != nil {
		return err
	}

	//Create Main file
	if err := CreateMainFile(template); err != nil {
		return err
	}

	//format project
	if err := FormatProject(); err != nil {
		return err
	}

	//Exec all transaction
	if err := ExecTransactions(transactions); err != nil {
		return err
	}

	return nil
}
