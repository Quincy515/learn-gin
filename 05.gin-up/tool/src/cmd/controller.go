/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	Helper "tool/src/helper"

	"github.com/spf13/cobra"
)

// controllerCmd represents the controller command
var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "控制器",
	Long:  `生成控制器.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 根据模板生成控制器
		Helper.GenFile(string(Helper.ReadFile("src/templates/controller.tpl")),
			fmt.Sprintf("/src/classes/%sClass.go", Helper.Ucfirst(args[0])),
			map[string]interface{}{
				"ControllerName": args[0],
			},
		)
		fmt.Println("控制器生成成功")
	},
	Args: cobra.MinimumNArgs(1), // 至少有一个参数
}

func init() {
	newCmd.AddCommand(controllerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// controllerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// controllerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
