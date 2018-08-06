package gen

import (
	"os"
	"io/ioutil"
	"goreav/logging"
)

type IAppTransaction interface {
	Apply() error
	Revert() error
}

type AppTransactionStack []IAppTransaction

type AppTransactionCreateDir struct {
	Path string
	Mode os.FileMode
}

func (t *AppTransactionCreateDir) Apply() error {
	return os.MkdirAll(t.Path, t.Mode)
}

func (t *AppTransactionCreateDir) Revert() error {
	return os.RemoveAll(t.Path)
}

type AppTransactionCreateFile struct {
	Path string
	file *os.File
}

func (t *AppTransactionCreateFile) GetFile() *os.File {
	return t.file
}

func (t *AppTransactionCreateFile) Apply() error {
	var err error
	t.file, err = os.Create(t.Path)
	return err
}

func (t *AppTransactionCreateFile) Revert() error {
	return os.RemoveAll(t.Path)
}

type AppTransactionAppendFile struct {
	Path        string
	Data        []byte
	currentData []byte
}

func (t *AppTransactionAppendFile) Apply() error {
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
	file, err := os.OpenFile(t.Path, os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(t.currentData)
	return err
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
				logging.Error.Fatal(err)
				return err
			}
		}
	}

	return nil
}
