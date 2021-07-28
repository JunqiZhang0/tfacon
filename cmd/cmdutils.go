package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var platformURL string
var tfaURL string
var viper0 *viper.Viper

func initConfig(cmd *cobra.Command) {
	viper0 = viper.New()
	viper0.AddConfigPath(".")
	viper0.SetConfigName("tfa")
	//if
	viper0.AutomaticEnv()
	initViperVal(cmd, viper0, "tfa-url", "TFAUL", "default val for tfa url", "The url to the TFA Classifier")
	initViperVal(cmd, viper0, "platform-url", "PLATFORMURL", "default val for platform url", "The url to the test platform")
}

func initViperVal(cmd *cobra.Command, viper *viper.Viper, cmdName, valName, defaultVal, cmdDescription string) {

	if viper.GetString(valName) == "" {
		viper.SetDefault(valName, defaultVal)
	} else {
		viper.SetDefault(valName, viper.GetString(valName))
	}
	cmd.PersistentFlags().StringP(cmdName, "", viper.GetString(valName), cmdDescription)
	viper.BindPFlag(valName, cmd.PersistentFlags().Lookup(cmdName))
	viper.ReadInConfig()
}

func printGreen(str string) {
	color.Green(str)
}
