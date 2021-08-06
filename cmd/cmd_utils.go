package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
}

func initConfig(viper *viper.Viper, cmd *cobra.Command, cmdInfoList []map[string]string) {

	viper.AddConfigPath(".")
	viper.SetConfigName("tfacon")
	viper.AutomaticEnv()

	for _, v := range cmdInfoList {

		initViperVal(cmd, viper, v["cmdName"], v["valName"], v["defaultVal"], v["cmdDescription"])
	}
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

func initTFAConfigFile(viper *viper.Viper) {
	var file []byte
	var err error
	if os.Getenv("TFACON_CONFIG_PATH") != "" {
		file, err = ioutil.ReadFile(os.Getenv("TFACON_CONFIG_PATH"))
	} else {
		file, err = ioutil.ReadFile("./tfacon.cfg")
	}
	defer func() {
		if r := recover(); r != nil {
			os.Create("tfacon.cfg")
		}
	}()
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("ini")
	viper.SetDefault("config.concurrency", true)
	viper.SetDefault("config.retry_times", 1)
	viper.ReadConfig(bytes.NewBuffer(file))
}
