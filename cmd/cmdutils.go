package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var platformURL string
// var tfaURL string
var viper0 *viper.Viper

var cmdInfoList []map[string]string = []map[string]string{
	{
		"cmdName":        "tfa-url",
		"valName":        "TFAURL",
		"defaultVal":     "default val for tfa url",
		"cmdDescription": "The url to the TFA Classifier",
	},
	{
		"cmdName":        "platform-url",
		"valName":        "PLATFORMURL",
		"defaultVal":     "default val for platform url",
		"cmdDescription": "The url to the test platform",
	},
}

func initConfig(cmd *cobra.Command, cmdInfoList []map[string]string) {
	viper0 = viper.New()
	viper0.AddConfigPath(".")
	viper0.SetConfigName("tfa")
	viper0.AutomaticEnv()
	for _, v := range cmdInfoList {

		initViperVal(cmd, viper0, v["cmdName"], v["valName"], v["defaultVal"], v["cmdDescription"])
	}
	//initViperVal(cmd, viper0, "tfa-url", "TFAURL", "default val for tfa url", "The url to the TFA Classifier")
	//initViperVal(cmd, viper0, "platform-url", "PLATFORMURL", "default val for platform url", "The url to the test platform")
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
func printHeader() {
	fmt.Println("--------------------------------------------------")
	fmt.Printf("tfactl  %s\n", rootCmd.Version)
	fmt.Println("Copyright (C) 2021, Red Hat, Inc.")
	fmt.Print("-------------------------------------------------\n\n\n")
	log.Println("Printing the constructed information")
}
