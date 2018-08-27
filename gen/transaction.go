package gen

import (
	"os"
	"io/ioutil"
	"goreav/logging"
	"text/template"
	"strings"
	"os/exec"
	"bytes"
)

type IAppTransaction interface {
	Apply() error
	Revert() error
	GetResult() interface{}
}

type AppTransactionStack []IAppTransaction

//Create dir transaction
type AppTransactionCreateDir struct {
	Path string
	Mode os.FileMode
}

func (t *AppTransactionCreateDir) Apply() error {
	logging.Info.Printf("Create Dir: %s", t.Path)
	return os.MkdirAll(t.Path, t.Mode)
}

func (t *AppTransactionCreateDir) Revert() error {
	logging.Info.Printf("Remove Dir: %s", t.Path)
	return os.RemoveAll(t.Path)
}

func (t *AppTransactionCreateDir) GetResult() interface{} {
	return nil
}

//Create file transaction
type AppTransactionCreateFile struct {
	Path string
	file *os.File
}

func (t *AppTransactionCreateFile) Apply() error {
	logging.Info.Printf("Create File: %s", t.Path)
	var err error
	t.file, err = os.Create(t.Path)
	defer t.file.Close()
	return err
}

func (t *AppTransactionCreateFile) Revert() error {
	logging.Info.Printf("Remove File: %s", t.Path)
	return os.RemoveAll(t.Path)
}

func (t *AppTransactionCreateFile) GetResult() interface{} {
	return t.file
}

//Append file transaction
type AppTransactionAppendFile struct {
	Path        string
	Data        []byte
	currentData []byte
}

func (t *AppTransactionAppendFile) Apply() error {
	logging.Info.Printf("Append File: %s", t.Path)
	var err error
	t.currentData, err = ioutil.ReadFile(t.Path)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(t.Path, os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteAt(t.Data, int64(len(t.currentData)))
	return err
}

func (t *AppTransactionAppendFile) Revert() error {
	logging.Info.Printf("Restore File: %s", t.Path)
	file, err := os.OpenFile(t.Path, os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(t.currentData)
	return err
}

func (t *AppTransactionAppendFile) GetResult() interface{} {
	return append(t.currentData, t.Data...)
}

//Create file from text template transaction
type AppTransactionCreateFileFromTemplate struct {
	Path         string
	TemplatePath string
	Data         interface{}
	data         []byte
	file         *os.File
}

func (t *AppTransactionCreateFileFromTemplate) Apply() error {
	logging.Info.Printf("Create File From Template: %s", t.Path)
	var err error
	t.data, err = ioutil.ReadFile(t.TemplatePath)
	if err != nil {
		return err
	}

	t.file, err = os.Create(t.Path)
	if err != nil {
		return err
	}
	defer t.file.Close()

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}

	tml := template.Must(template.New("").Funcs(funcMap).Parse(string(t.data)))

	return tml.Execute(t.file, t.Data)
}

func (t *AppTransactionCreateFileFromTemplate) Revert() error {
	logging.Info.Printf("Remove File From Template: %s", t.Path)
	return os.RemoveAll(t.Path)
}

func (t *AppTransactionCreateFileFromTemplate) GetResult() interface{} {
	return t.file
}

//Format project transaction
type AppTransactionFormatProject struct {
	Path         string
	result       []byte
}

func (t *AppTransactionFormatProject) Apply() error {
	logging.Info.Printf("Format project: %s", t.Path)
	cmd := exec.Command("go", "fmt", t.Path)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	t.result = out.Bytes()
	return nil
}

func (t *AppTransactionFormatProject) Revert() error {
	return nil
}

func (t *AppTransactionFormatProject) GetResult() interface{} {
	return t.result
}

//Execute transaction
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
				logging.Error.Fatal(err)
				return err
			}
		}
	}

	return nil
}
