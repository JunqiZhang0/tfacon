/*
Copyright © 2021 Red Hat D&O Tools Development Team

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"tfactl/utils"

	"github.com/spf13/cobra"
)

var platform utils.Platform
var tfa utils.TFA

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all constructed information",
	Long:  `list all information constructed from tfa.yml/enviroment variables/cli`,
	Run: func(cmd *cobra.Command, args []string) {
		printHeader()
		viper0.Unmarshal(&platform)
		viper0.Unmarshal(&tfa)
		printGreen(fmt.Sprintf("TFA Classifier is: %s\n", tfa))
		printGreen(fmt.Sprintf("Test Platform is: %s\n", platform))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	initConfig(listCmd, cmdInfoList)
}
