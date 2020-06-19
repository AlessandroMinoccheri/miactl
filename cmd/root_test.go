package cmd

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/mia-platform/miactl/fs"
	"github.com/mia-platform/miactl/renderer"
	"github.com/mia-platform/miactl/sdk"

	"github.com/spf13/cobra"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandWithContext(ctx context.Context, root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.ExecuteContext(ctx)

	return buf.String(), err
}

type testOutput struct {
	text string

	factory *Factory
}

func executeRootCommandWithContext(mockError sdk.MockClientError, args ...string) (output testOutput, err error) {
	rootCmd := NewRootCmd()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	f := Factory{
		renderer:         renderer.New(rootCmd.OutOrStderr()),
		miaClientCreator: sdk.WrapperMockMiaClient(mockError),
		fs:               fs.MockFs(),
		homeDir:          "testdata",
	}

	ctx := context.WithValue(context.Background(), FactoryContextKey{}, f)

	err = rootCmd.ExecuteContext(ctx)

	return testOutput{
		text:    buf.String(),
		factory: &f,
	}, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
