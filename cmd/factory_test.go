package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/mia-platform/miactl/fs"
	"github.com/mia-platform/miactl/renderer"
	"github.com/mia-platform/miactl/sdk"

	"github.com/stretchr/testify/require"
)

func TestWithFactoryValue(t *testing.T) {
	t.Run("save factory to passed context", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithFactoryValue(ctx, &bytes.Buffer{}, "")
		f := ctx.Value(FactoryContextKey{})
		require.NotNil(t, f)
		if _, ok := f.(Factory); ok {
			return
		}
		t.Fail()
	})
}

func TestAddMiaClientToFactory(t *testing.T) {
	t.Run("throws if newSdk not defined", func(t *testing.T) {
		f := Factory{}
		require.NotNil(t, f)
		opts := sdk.Options{}
		err := f.addMiaClientToFactory(opts)
		require.EqualError(t, err, fmt.Sprintf("%s: newSdk not defined", sdk.ErrCreateClient))
	})

	t.Run("throws if options are not passed", func(t *testing.T) {
		f := Factory{
			miaClientCreator: sdk.New,
		}
		require.NotNil(t, f)
		opts := sdk.Options{}
		err := f.addMiaClientToFactory(opts)
		require.EqualError(t, err, fmt.Sprintf("%s: client options are not correct", sdk.ErrCreateClient))
	})

	t.Run("add MiaClient to factory", func(t *testing.T) {
		opts := sdk.Options{
			APIKey:     "my-apiKey",
			APIBaseURL: "http://base-url.com/",
			APICookie:  "cookie",
		}
		miaClient, err := sdk.New(opts)
		require.NoError(t, err)
		miaClientCreator := func(optsArg sdk.Options) (*sdk.MiaClient, error) {
			require.Equal(t, opts, optsArg)
			return miaClient, nil
		}
		f := Factory{
			miaClientCreator: miaClientCreator,
		}
		require.NotNil(t, f)
		err = f.addMiaClientToFactory(opts)
		require.NoError(t, err)

		require.Equal(t, miaClient, f.MiaClient())
		require.Equal(t, reflect.ValueOf(miaClientCreator).Pointer(), reflect.ValueOf(f.miaClientCreator).Pointer())
	})
}

func TestRendererMethod(t *testing.T) {
	t.Run("panic if renderer is nil", func(t *testing.T) {
		f := Factory{}
		require.PanicsWithError(t, fmt.Sprintf("%s: renderer not defined", errFactory), func() { f.Renderer() })
	})

	t.Run("returns renderer correctly", func(t *testing.T) {
		f := Factory{
			renderer: renderer.New(os.Stdout),
		}
		require.Equal(t, renderer.New(os.Stdout), f.Renderer())
	})
}

func TestMiaClientMethod(t *testing.T) {
	t.Run("panic if miaClient is nil", func(t *testing.T) {
		f := Factory{}
		require.PanicsWithError(t, fmt.Sprintf("%s: mia client not defined", errFactory), func() { f.MiaClient() })
	})

	t.Run("returns miaClient correctly", func(t *testing.T) {
		f := Factory{
			miaClient: &sdk.MiaClient{},
		}
		client := &sdk.MiaClient{}
		require.Equal(t, client, f.MiaClient())
	})
}

func TestFsMethod(t *testing.T) {
	t.Run("panic if renderer is nil", func(t *testing.T) {
		f := Factory{}
		require.PanicsWithError(t, fmt.Sprintf("%s: fs not defined", errFactory), func() { f.Fs() })
	})

	t.Run("returns renderer correctly", func(t *testing.T) {
		f := Factory{
			fs: fs.MockFs(),
		}
		require.Equal(t, fs.MockFs(), f.Fs())
	})
}

func TestGetFactoryFromContext(t *testing.T) {
	t.Run("throws if context error", func(t *testing.T) {
		ctx, cancFn := context.WithTimeout(context.Background(), 0)
		defer cancFn()
		f, err := GetFactoryFromContext(ctx, sdk.Options{})

		require.Nil(t, f)
		require.EqualError(t, err, errFactory.Error())
	})

	t.Run("throws if mia client error", func(t *testing.T) {
		ctx := context.Background()
		buf := &bytes.Buffer{}
		ctx = WithFactoryValue(ctx, buf, "")
		f, err := GetFactoryFromContext(ctx, sdk.Options{})

		require.Nil(t, f)
		require.Error(t, err)
		require.EqualError(t, err, fmt.Sprintf("%s: client options are not correct", sdk.ErrCreateClient))
	})

	t.Run("returns factory", func(t *testing.T) {
		configPath := "config/path"
		ctx := context.Background()
		ctx = WithFactoryValue(ctx, &bytes.Buffer{}, configPath)
		opts := sdk.Options{
			APIBaseURL: "http://base-url.com/",
			APICookie:  "cookie",
			APIKey:     "my-APIKey",
		}

		f, err := GetFactoryFromContext(ctx, opts)
		require.NoError(t, err)

		miaClient, err := sdk.New(opts)

		require.NoError(t, err)
		require.Equal(t, renderer.New(&bytes.Buffer{}), f.Renderer())
		require.Equal(t, miaClient, f.MiaClient())
		require.Equal(t, reflect.ValueOf(sdk.New).Pointer(), reflect.ValueOf(f.miaClientCreator).Pointer())

		require.Equal(t, configPath, f.homeDir)
	})
}
