package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var platformURL string
var tfaURL string
var viper0 *viper.Viper

type Info struct {
	PlatformURL string `mapstructure:"platformURL"`
	TFAURL      string `mapstructure:"tfaURL"`
}

func initConfig(cmd *cobra.Command) {
	viper0 = viper.New()
	viper0.AddConfigPath(".")
	viper0.SetConfigName("tfa")
	//if
	viper0.AutomaticEnv()
	if viper0.GetString("PLATFORMURL") == "" {
		viper0.SetDefault("platformURL", "https://platform.com/test")
	} else {
		viper0.SetDefault("platformURL", viper0.GetString("PLATFORMURL"))
	}
	if viper0.GetString("TFAURL") == "" {
		viper0.SetDefault("tfaURL", "default tfaURL")
	} else {
		viper0.SetDefault("tfaURL", viper0.GetString("TFAURL"))
	}
	viper0.ReadInConfig()
	cmd.PersistentFlags().StringVarP(&platformURL, "platform-url", "", viper0.GetString("platformURL"), "The url to the platform")
	cmd.PersistentFlags().StringVarP(&tfaURL, "tfa-url", "", viper0.GetString("tfaURL"), "The url to the TFA Classifer")
}

func printGreen(str string) {
	color.Green(str)
}
