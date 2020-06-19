package fs

import (
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// Fs structure to be used in miactl.
type Fs struct {
	afero.Fs
}

// New creates a new fs object.
func New() *Fs {
	return &Fs{
		afero.NewOsFs(),
	}
}

func (fs *Fs) createFile(p string) (afero.File, error) {
	if err := fs.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return fs.Create(p)
}

// WriteYAMLFile write a yaml file.
func (fs *Fs) WriteYAMLFile(filePath string, content interface{}) error {
	file, err := fs.createFile(filePath)
	if err != nil {
		return err
	}

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(content)
	if err != nil {
		return err
	}

	return nil
}

// Exists execute afero exists function.
func (fs *Fs) Exists(path string) (bool, error) {
	return afero.Exists(fs, path)
}

// ReadFile read file content.
func (fs *Fs) ReadFile(path string) ([]byte, error) {
	return afero.ReadFile(fs, path)
}
