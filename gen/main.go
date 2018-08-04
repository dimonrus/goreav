package gen

import (
	"fmt"
	"errors"
	"goreav/logging"
)

type AppTemplate map[string]interface{}

var path string //full path to project

var transactions AppTransactionStack

func CreateProject(template AppTemplate) error {
	//Project section is required
	if project, ok := template[KeyWordProject]; ok == true {
		data := project.(map[interface{}]interface{})
		path = fmt.Sprintf("%s/%s", data[KeyWordPath], data[KeyWordName])
		transactions = append(transactions, &AppTransactionCreateDir{Path: path, Mode: 0755})
	} else {
		return errors.New("template has no project section")
	}

	return nil
}

func RenderConfig(template AppTemplate) error {
	//Env section is not required
	if environment, ok := template[KeyWordEnvironment]; ok == true {
		configPath := path + "/config"
		transactions = append(transactions, &AppTransactionCreateDir{Path: configPath, Mode: 0755})
		//Create config files
		for key, _ := range environment.(map[interface{}]interface{}) {
			filePath := configPath + "/" + key.(string) + ".yaml"
			transactions = append(transactions, &AppTransactionCreateFile{Path: filePath})
		}
	}

	return nil
}

func ExecTransactions(txs []IAppTransaction) error {
	//Apply app transactions
	var stopped *int
	for index, trx := range txs {
		err := trx.Apply()
		if err != nil {
			logging.Error.Print(err)
			stopped = &index
			break
		}
	}

	//Rollback transactions
	if stopped != nil {
		for i := *stopped; i >= 0; i-- {
			err := transactions[i].Revert()
			if err != nil {
				logging.Error.Print(err)
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

	//Exec all transaction
	return ExecTransactions(transactions)
}
