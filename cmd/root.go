package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mia-platform/miactl/sdk"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	opts = sdk.Options{}
)

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "miactl",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.MarkFlagRequired("apiKey")
			cmd.MarkFlagRequired("apiCookie")
			cmd.MarkFlagRequired("apiBaseUrl")
		},
	}
	setRootPersistentFlag(rootCmd)

	// add sub command to root command
	rootCmd.AddCommand(newGetCmd())
	rootCmd.AddCommand(newSetCommand())

	rootCmd.AddCommand(newCompletionCmd(rootCmd))
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	configPath, err := getConfigDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd := NewRootCmd()
	ctx := WithFactoryValue(context.Background(), rootCmd.OutOrStdout(), configPath)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func setRootPersistentFlag(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVar(&opts.APIKey, "apiKey", "", "API Key")
	rootCmd.PersistentFlags().StringVar(&opts.APICookie, "apiCookie", "", "api cookie sid")
	rootCmd.PersistentFlags().StringVar(&opts.APIBaseURL, "apiBaseUrl", "", "api base url")
}

func getConfigDir() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", home, ".miactl"), nil
}

func initConfig() {
	cfgFile := "config"

	configPath, err := getConfigDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(cfgFile)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
