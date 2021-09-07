package cmd

import (
	"fmt"
	"os"

	"github.com/danielgatis/go-between/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel string
	logJSON  bool
)

var rootCmd = &cobra.Command{
	Use:   "go-between",
	Short: "A distributed order matchmaking",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-between.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "INFO", "set the log level")
	rootCmd.PersistentFlags().BoolVar(&logJSON, "log-json", false, "set the format as json")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		logger.GetLogrusInstance().Fatal(err)
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".go-between")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
