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
	"github.com/JunqiZhang0/tfacon/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all constructed information",
	Long:  `list all information constructed from tfacon.yml/enviroment variables/cli`,
	Run: func(cmd *cobra.Command, args []string) {
		printHeader()
		con := core.GetInfo(viperList)
		printGreen(con.String())
	},
}
var viperList *viper.Viper

func init() {
	rootCmd.AddCommand(listCmd)
	viperList = viper.New()
	initConfig(viperList, listCmd, cmdInfoList)
}
