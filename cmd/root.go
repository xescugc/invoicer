package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	RootCmd = &cobra.Command{
		Use: "invoicer",
	}
)

func init() {
	cobra.OnInitialize(initViper)

	RootCmd.AddCommand(
		userCmd,
		customersCmd,
		invoicesCmd,
		templatesCmd,
		uiCmd,
	)

	RootCmd.PersistentFlags().String("config", "", "Config flag")
	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
}

func initViper() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
