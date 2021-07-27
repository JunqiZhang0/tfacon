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
	"log"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all constructed information",
	Long:  `list all information constructed from tfa.yml/enviroment variables/cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("--------------------------------------------------")
		fmt.Printf("tfactl  %s\n", rootCmd.Version)
		fmt.Println("Copyright (C) 2021, Red Hat, Inc.")
		fmt.Print("-------------------------------------------------\n\n\n")
		log.Println("Printing the constructed information")
		//printGreen(fmt.Sprintf("Platform URL is: %s\n", platformURL))
		//printGreen(fmt.Sprintf("TFA Classifier URL is: %s\n", tfaURL))
		var info Info
		viper0.Unmarshal(&info)
		fmt.Print(info.PlatformURL, "\n", info.TFAURL)
		//fmt.Print(viper0.GetStringMapString("TFAConfig"))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	initConfig(listCmd)
}
