package fs

import "github.com/spf13/afero"

// MockFs is a function to mock fs.
// This function is only for test.
func MockFs() *Fs {
	return &Fs{
		afero.NewMemMapFs(),
	}
}
