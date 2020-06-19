package fs

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestNew(t *testing.T) {
	t.Run("returns os fs", func(t *testing.T) {
		fs := New()
		require.Equal(t, &Fs{afero.NewOsFs()}, fs)
	})
}

func TestWriteYAMLFile(t *testing.T) {
	type fileContent struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	}

	content := fileContent{
		Name: "my name",
		Type: "my type",
	}
	appFs := New()

	t.Run("write a new file", func(t *testing.T) {
		filePath := "testdata/write-yaml/path/file.yml"
		t.Cleanup(func() {
			appFs.RemoveAll("testdata/write-yaml")
		})

		err := appFs.WriteYAMLFile(filePath, &content)
		require.NoError(t, err)

		isFileExistent, err := afero.Exists(appFs, filePath)
		require.NoError(t, err)
		require.True(t, isFileExistent)

		t.Run("with correct content", func(t *testing.T) {
			file, err := afero.ReadFile(appFs, filePath)
			require.NoError(t, err)
			var actualSavedContext = fileContent{}
			err = yaml.Unmarshal(file, &actualSavedContext)
			require.NoError(t, err)
			require.Equal(t, content, actualSavedContext)
		})

		t.Run("with correct string content", func(t *testing.T) {
			content, err := afero.ReadFile(appFs, filePath)
			require.NoError(t, err)
			require.NoError(t, err)
			require.YAMLEq(t, `name: my name
type: my type
`, string(content))
		})
	})

	t.Run("fails to create file", func(t *testing.T) {
		filePath := "testdata/file"
		appFs.createFile(filePath)
		t.Cleanup(func() {
			appFs.RemoveAll(filePath)
		})

		err := appFs.WriteYAMLFile(fmt.Sprintf("%s/file.yaml", filePath), content)
		require.Error(t, err)
	})
}

func TestCreateFile(t *testing.T) {
	appFs := New()

	t.Run("correctly create new empty file", func(t *testing.T) {
		filePath := "testdata/test-file"
		appFs.Remove(filePath)

		_, err := appFs.createFile(filePath)
		t.Cleanup(func() {
			appFs.Remove(filePath)
		})

		require.NoError(t, err)
		isFileExistent, err := appFs.Exists(filePath)
		require.NoError(t, err)
		require.True(t, isFileExistent)

		content, err := appFs.ReadFile(filePath)
		require.Equal(t, []byte{}, content)
	})

	t.Run("correctly create new empty file under multiple sub directories", func(t *testing.T) {
		filePath := "testdata/nested/new/path/dir/test-file"
		appFs.RemoveAll("testdata/nested")

		_, err := appFs.createFile(filePath)
		t.Cleanup(func() {
			appFs.RemoveAll("testdata/nested")
		})

		require.NoError(t, err)
		isFileExistent, err := appFs.Exists(filePath)
		require.NoError(t, err)
		require.True(t, isFileExistent)

		content, err := appFs.ReadFile(filePath)
		require.Equal(t, []byte{}, content)
	})

	t.Run("returns no error if file already exists", func(t *testing.T) {
		filePath := "testdata/test-file-exists"
		t.Cleanup(func() {
			appFs.Remove(filePath)
		})

		_, err := appFs.createFile(filePath)
		require.NoError(t, err)

		_, err = appFs.createFile(filePath)
		require.NoError(t, err)
	})

	t.Run("returns error if directory to create is already a file", func(t *testing.T) {
		filePath := "testdata/test-collision"
		t.Cleanup(func() {
			appFs.RemoveAll(filePath)
		})
		_, err := appFs.createFile(filePath)
		require.NoError(t, err)

		_, err = appFs.createFile(fmt.Sprintf("%s/file.yml", filePath))
		require.Error(t, err)
	})
}

func TestReadFile(t *testing.T) {
	appFs := New()

	t.Run("read file", func(t *testing.T) {
		content, err := appFs.ReadFile("testdata/file.yaml")
		require.NoError(t, err)
		require.YAMLEq(t, string(content), "file: content")
	})

	t.Run("returns error if file not exists", func(t *testing.T) {
		content, err := appFs.ReadFile("testdata/not-exists")
		require.Empty(t, content)
		require.Error(t, err)
	})
}

func TestExists(t *testing.T) {
	appFs := New()

	t.Run("returns true if file exists", func(t *testing.T) {

		exists, err := appFs.Exists("testdata/file.yaml")
		require.NoError(t, err)
		require.True(t, exists, "file exists")
	})

	t.Run("returns false if file exists", func(t *testing.T) {
		exists, err := appFs.Exists("testdata/not-exists")
		require.NoError(t, err)
		require.False(t, exists, "file exists")
	})
}
