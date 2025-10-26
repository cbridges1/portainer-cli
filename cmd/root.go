package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "portainer-cli",
	Short: "A CLI tool for managing Portainer stacks",
	Long: `A command line interface for managing Docker Compose stacks through the Portainer API.
This tool allows you to create, list, update, delete, and inspect stacks in your Portainer environment.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.portainer-cli.yaml)")
	rootCmd.PersistentFlags().String("url", "", "Portainer URL")
	rootCmd.PersistentFlags().String("token", "", "Portainer API token")
	rootCmd.PersistentFlags().Int("endpoint", 0, "Portainer endpoint ID")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".portainer-cli")
	}

	viper.SetEnvPrefix("PORTAINER")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
