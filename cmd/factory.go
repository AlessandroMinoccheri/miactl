package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/mia-platform/miactl/fs"
	"github.com/mia-platform/miactl/renderer"
	"github.com/mia-platform/miactl/sdk"
)

var errFactory = errors.New("factory error")

// FactoryContextKey key of the factory in context
type FactoryContextKey struct{}

type miaClientCreator func(opts sdk.Options) (*sdk.MiaClient, error)

// Factory returns all the clients around the commands
type Factory struct {
	renderer  renderer.IRenderer
	miaClient *sdk.MiaClient
	fs        *fs.Fs

	miaClientCreator miaClientCreator

	homeDir string
}

func (o *Factory) addMiaClientToFactory(opts sdk.Options) error {
	if o.miaClientCreator == nil {
		return fmt.Errorf("%w: newSdk not defined", sdk.ErrCreateClient)
	}
	miaSdk, err := o.miaClientCreator(opts)
	if err != nil {
		return err
	}
	o.miaClient = miaSdk
	return nil
}

// Renderer method to access to renderer field
func (o *Factory) Renderer() renderer.IRenderer {
	if o.renderer == nil {
		panic(fmt.Errorf("%w: renderer not defined", errFactory))
	}
	return o.renderer
}

// MiaClient method to access to miaClient field
func (o *Factory) MiaClient() *sdk.MiaClient {
	if o.miaClient == nil {
		panic(fmt.Errorf("%w: mia client not defined", errFactory))
	}
	return o.miaClient
}

// Fs method to access to fs field
func (o *Factory) Fs() *fs.Fs {
	if o.fs == nil {
		panic(fmt.Errorf("%w: fs not defined", errFactory))
	}
	return o.fs
}

// WithFactoryValue add factory to passed context
func WithFactoryValue(ctx context.Context, writer io.Writer, homeDir string) context.Context {
	return context.WithValue(ctx, FactoryContextKey{}, Factory{
		renderer:         renderer.New(writer),
		miaClientCreator: sdk.New,
		fs:               fs.New(),
		homeDir:          homeDir,
	})
}

// GetFactoryFromContext returns factory from context
func GetFactoryFromContext(ctx context.Context, opts sdk.Options) (*Factory, error) {
	factory, ok := ctx.Value(FactoryContextKey{}).(Factory)
	if !ok {
		return nil, errFactory
	}

	err := factory.addMiaClientToFactory(opts)
	if err != nil {
		return nil, err
	}

	return &factory, nil
}
