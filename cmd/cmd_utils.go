package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdInfoList []map[string]string = []map[string]string{
	{
		"cmdName":        "tfa-url",
		"valName":        "TFA_URL",
		"defaultVal":     "default val for tfa url",
		"cmdDescription": "The url to the TFA Classifier",
	},
	{
		"cmdName":        "platform-url",
		"valName":        "PLATFORM_URL",
		"defaultVal":     "default val for platform url",
		"cmdDescription": "The url to the test platform(example: https://reportportal-ccit.apps.ocp4.prod.psi.redhat.com)",
	},
	{
		"cmdName":        "connector-type",
		"valName":        "CONNECTOR_TYPE",
		"defaultVal":     "RPCon",
		"cmdDescription": "The type of connector you want to use(example: RPCon, PolarionCon, JiraCon)",
	},
	{
		"cmdName":        "launch-id",
		"valName":        "LAUNCH_ID",
		"defaultVal":     "",
		"cmdDescription": "The launch id of report portal",
	},
	{
		"cmdName":        "project-name",
		"valName":        "PROJECT_NAME",
		"defaultVal":     "",
		"cmdDescription": "The project name of report portal",
	},
	{
		"cmdName":        "auth-token",
		"valName":        "AUTH_TOKEN",
		"defaultVal":     "",
		"cmdDescription": "The AUTH_TOKEN of report portal",
	},
	// Client      *http.Client
}

func initConfig(viper *viper.Viper, cmd *cobra.Command, cmdInfoList []map[string]string) {

	viper.AddConfigPath(".")
	viper.SetConfigName("tfacon")
	viper.AutomaticEnv()

	for _, v := range cmdInfoList {

		initViperVal(cmd, viper, v["cmdName"], v["valName"], v["defaultVal"], v["cmdDescription"])
	}
	//initViperVal(cmd, viper0, "tfa-url", "TFAURL", "default val for tfa url", "The url to the TFA Classifier")
	//initViperVal(cmd, viper0, "platform-url", "PLATFORMURL", "default val for platform url", "The url to the test platform")
	viper.ReadInConfig()
}

func initViperVal(cmd *cobra.Command, viper *viper.Viper, cmdName, valName, defaultVal, cmdDescription string) {

	if viper.GetString(valName) == "" {
		viper.SetDefault(valName, defaultVal)
	} else {
		viper.SetDefault(valName, viper.GetString(valName))
	}
	cmd.PersistentFlags().StringP(cmdName, "", viper.GetString(valName), cmdDescription)
	viper.BindPFlag(valName, cmd.PersistentFlags().Lookup(cmdName))

}

func printGreen(str string) {
	color.Green(str)
}
func printHeader() {
	fmt.Println("--------------------------------------------------")
	fmt.Printf("tfacon  %s\n", rootCmd.Version)
	fmt.Println("Copyright (C) 2021, Red Hat, Inc.")
	fmt.Print("-------------------------------------------------\n\n\n")
	log.Println("Printing the constructed information")
}
