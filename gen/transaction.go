package gen

import "os"

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
