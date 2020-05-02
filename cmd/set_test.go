package cmd

import (
	"testing"

	"github.com/mia-platform/miactl/sdk"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestNewSetCommand(t *testing.T) {
	t.Run("returns a command", func(t *testing.T) {
		cmd := newSetCommand()
		require.NotNil(t, cmd)
	})

	t.Run("command name is set", func(t *testing.T) {
		cmd := newSetCommand()
		require.NotNil(t, cmd)
		require.Equal(t, "set", cmd.Use)
	})

	t.Run("accept context as argument", func(t *testing.T) {
		cmd := newSetCommand()
		require.NotNil(t, cmd)
		require.Contains(t, cmd.ValidArgs, "context")
	})

	t.Run("accept only context as argument", func(t *testing.T) {
		cmd := newSetCommand()
		require.NotNil(t, cmd)
		require.Error(t, cmd.Args(cmd, []string{"something-wrong"}))
		require.Error(t, cmd.Args(cmd, []string{}))
		require.Error(t, cmd.Args(cmd, nil))
		require.NoError(t, cmd.Args(cmd, []string{"context"}))
	})
}

func TestSetContextCommand(t *testing.T) {
	t.Run("not returns error", func(t *testing.T) {
		out, err := executeRootCommandWithContext(sdk.MockClientError{}, "set", "context")
		require.Equal(t, "Context created", out)
		require.NoError(t, err)
	})
}

func TestWriteContextFile(t *testing.T) {
	t.Run("write a new file", func(t *testing.T) {
		appFs := afero.NewMemMapFs()
		f := &Factory{
			fs: appFs,
		}
		filePath := "/my/path"
		miaContext := MiaContext{
			APIBaseURL: "https://my-host",
			APIKey:     "api-key",
		}
		err := writeContextFile(f, filePath, &miaContext)
		require.NoError(t, err)
		isFileExistent, err := afero.Exists(appFs, filePath)
		require.NoError(t, err)
		require.True(t, isFileExistent)

		t.Run("with correct content", func(t *testing.T) {
			content, err := afero.ReadFile(appFs, filePath)
			require.NoError(t, err)
			actualSavedContext := MiaContext{}
			err = yaml.Unmarshal(content, &actualSavedContext)
			require.NoError(t, err)
			require.Equal(t, miaContext, actualSavedContext)
		})

		t.Run("with correct string content", func(t *testing.T) {
			content, err := afero.ReadFile(appFs, filePath)
			require.NoError(t, err)
			require.NoError(t, err)
			require.Equal(t, `apiBaseUrl: https://my-host
apiKey: api-key
`, string(content))
		})
	})
}
