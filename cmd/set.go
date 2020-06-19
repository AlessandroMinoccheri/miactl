package cmd

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var errCreateContext = errors.New("fails to create mia context")

func newSetCommand() *cobra.Command {
	setCmd := &cobra.Command{
		Use: "set",
	}

	miaContext := &MiaContext{}

	setCmd.AddCommand(miaContext.newSetContextCmd())

	return setCmd
}

// MiaContext define the context of a console
type MiaContext struct {
	APIBaseURL string `yaml:"apiBaseUrl"`
	APIKey     string `yaml:"apiKey"`
	Name       string `yaml:"name"`
}

func (m *MiaContext) newSetContextCmd() *cobra.Command {
	contextCmd := &cobra.Command{
		Use: "context",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := GetFactoryFromContext(cmd.Context(), opts)
			if err != nil {
				return err
			}

			m.APIBaseURL = opts.APIBaseURL
			m.APIKey = opts.APIKey

			if qs := m.getPromptQuestion(); qs != nil {
				err := survey.Ask(qs, m)
				if err == nil {
					return err
				}
			}
			err = m.createContextFile(f)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStderr(), "Context created")
			return nil
		},
	}

	contextCmd.Flags().StringVar(&m.Name, "name", "", "Set the context name")

	return contextCmd
}

func (m *MiaContext) getPromptQuestion() []*survey.Question {
	var qs = []*survey.Question{}
	if m.Name == "" {
		qs = append(qs, &survey.Question{
			Name:     "name",
			Prompt:   &survey.Input{Message: "Insert context name"},
			Validate: survey.Required,
		})
	}
	if m.APIBaseURL == "" {
		qs = append(qs, &survey.Question{
			Name:     "apiBaseURL",
			Prompt:   &survey.Input{Message: "Insert api base url"},
			Validate: survey.Required,
		})
	}
	if m.APIKey == "" {
		qs = append(qs, &survey.Question{
			Name:     "apiKey",
			Prompt:   &survey.Input{Message: "Insert api key"},
			Validate: survey.Required,
		})
	}
	if len(qs) == 0 {
		return nil
	}
	return qs
}

func (m *MiaContext) createContextFile(f *Factory) error {
	if m.Name == "" {
		return fmt.Errorf("%w: empty name", errCreateContext)
	}
	filePath := fmt.Sprintf("%s/contexts/%s", f.homeDir, m.Name)
	return f.Fs().WriteYAMLFile(filePath, m)
}
