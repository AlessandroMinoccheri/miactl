package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type MiaContext struct {
	APIBaseURL string
	APIKey     string
}

func newSetCommand() *cobra.Command {
	validArgs := []string{"context"}

	return &cobra.Command{
		Use:       "set",
		ValidArgs: validArgs,
		Args:      cobra.ExactValidArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// f, err := GetFactoryFromContext(cmd.Context(), opts)
			// if err != nil {
			// 	return err
			// }

			// writeContextFile(f, fmt.Sprintf("%s/contexts/%s", cfgHome, "prova.yml"), &MiaContext{})

			fmt.Fprint(cmd.OutOrStderr(), "Context created")
			return nil
		},
	}
}

func writeContextFile(f *Factory, filePath string, miaContext *MiaContext) error {
	file, err := f.Fs.Create(filePath)
	if err != nil {
		return err
	}

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	err = encoder.Encode(miaContext)
	if err != nil {
		return err
	}

	return nil
}
