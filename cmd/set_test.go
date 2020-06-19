package cmd

import (
	"fmt"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/mia-platform/miactl/fs"
	"github.com/mia-platform/miactl/sdk"
	"github.com/stretchr/testify/require"
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
}

func TestMiaContext(t *testing.T) {
	miaContext := MiaContext{
		Name:       "mia-ctx-name",
		APIBaseURL: "base-url",
		APIKey:     "api key",
	}

	t.Run("create context file", func(t *testing.T) {
		t.Run("returns error if name not set", func(t *testing.T) {
			f := &Factory{
				homeDir: "/",
			}
			miaContext := MiaContext{}
			err := miaContext.createContextFile(f)
			require.EqualError(t, err, fmt.Sprintf("%s: empty name", errCreateContext))
		})

		t.Run("panics if fs not in factory", func(t *testing.T) {
			f := &Factory{}
			require.PanicsWithError(t, fmt.Sprintf("%s: fs not defined", errFactory), func() {
				miaContext.createContextFile(f)
			})
		})

		t.Run("correctly create yaml file", func(t *testing.T) {
			memoryFs := fs.MockFs()
			f := &Factory{
				fs:      memoryFs,
				homeDir: "/home",
			}
			err := miaContext.createContextFile(f)
			require.NoError(t, err)

			expectedPath := fmt.Sprintf("/home/contexts/%s", miaContext.Name)
			isFileExistent, err := memoryFs.Exists(expectedPath)
			require.NoError(t, err)
			require.True(t, isFileExistent)

			content, err := memoryFs.ReadFile(expectedPath)
			require.NoError(t, err)
			require.YAMLEq(t, `apiBaseUrl: base-url
apiKey: api key
name: mia-ctx-name
`, string(content))
		})
	})

	t.Run("get prompt question", func(t *testing.T) {
		t.Run("returns empty array if mia context not filled", func(t *testing.T) {
			miaContext := MiaContext{
				Name:       "mia-ctx-name",
				APIBaseURL: "base-url",
				APIKey:     "api key",
			}
			q := miaContext.getPromptQuestion()
			require.Nil(t, q)
		})

		t.Run("returns name question if name not in mia context", func(t *testing.T) {
			miaContext := MiaContext{
				APIBaseURL: "base-url",
				APIKey:     "api key",
			}
			q := miaContext.getPromptQuestion()
			require.Len(t, q, 1)
			require.Equal(t, "name", q[0].Name)
			require.Equal(t, &survey.Input{Message: "Insert context name"}, q[0].Prompt)
			require.NotEmpty(t, q[0].Validate)
		})

		t.Run("returns api base url question if api base url not in mia context", func(t *testing.T) {
			miaContext := MiaContext{
				Name:   "mia-ctx-name",
				APIKey: "api key",
			}
			q := miaContext.getPromptQuestion()
			require.Len(t, q, 1)
			require.Equal(t, "apiBaseURL", q[0].Name)
			require.Equal(t, &survey.Input{Message: "Insert api base url"}, q[0].Prompt)
			require.NotEmpty(t, q[0].Validate)
		})

		t.Run("returns api key question if api key not in mia context", func(t *testing.T) {
			miaContext := MiaContext{
				Name:       "mia-ctx-name",
				APIBaseURL: "base-url",
			}
			q := miaContext.getPromptQuestion()
			require.Len(t, q, 1)
			require.Equal(t, "apiKey", q[0].Name)
			require.Equal(t, &survey.Input{Message: "Insert api key"}, q[0].Prompt)
			require.NotEmpty(t, q[0].Validate)
		})

		t.Run("returns all 3 questions if context is empty", func(t *testing.T) {
			miaContext := MiaContext{}
			q := miaContext.getPromptQuestion()

			require.Len(t, q, 3)
			require.Equal(t, "name", q[0].Name)
			require.Equal(t, "apiBaseURL", q[1].Name)
			require.Equal(t, "apiKey", q[2].Name)
		})
	})
}

func TestSetContextCommand(t *testing.T) {
	const contextName = "test-context-name"
	const apiBaseURL = "http://base-url/api/"
	const apiKey = "apiKey"

	t.Run("creates context file -- with all context data passed by flag", func(t *testing.T) {
		out, err := executeRootCommandWithContext(sdk.MockClientError{},
			"set", "context",
			"--apiBaseUrl", apiBaseURL,
			"--apiCookie", "my-cookie",
			"--apiKey", apiKey,
			"--name", contextName,
		)
		require.Equal(t, "Context created", out.text)
		require.NoError(t, err)

		t.Run("correctly creates file", func(t *testing.T) {
			homeDir := out.factory.homeDir
			filePath := fmt.Sprintf("%s/contexts/%s", homeDir, contextName)

			exists, err := out.factory.Fs().Exists(filePath)
			require.NoError(t, err)
			require.True(t, exists)

			content, err := out.factory.Fs().ReadFile(filePath)
			require.NoError(t, err)
			require.YAMLEq(t, fmt.Sprintf(`"apiBaseUrl": %s
"apiKey": %s
"name": %s
`, apiBaseURL, apiKey, contextName), string(content))
		})
	})
}
