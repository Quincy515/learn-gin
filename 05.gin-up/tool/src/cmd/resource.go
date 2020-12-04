package cmd

import (
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	Helper "tool/src/helper"
)

// resourceCmd represents the resource command
var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "自动生成模板变量",
	Long:  `自动生成模板变量`,
	Run: func(cmd *cobra.Command, args []string) {
		//遍历文件
		res := Helper.LoadResource("src/templates") //加载 资源文件
		tpl, err := ioutil.ReadFile(Helper.GetWorkDir() + "/src/resource/resource.tpl")
		if err != nil {
			log.Fatal("resource.tpl error", err)
		}
		if res != nil {
			tmpl, err := template.New("resource").Funcs(Helper.NewTplFunction()).Parse(string(tpl))
			if err != nil {
				log.Fatal("resource parse error:", err)
			}
			file, err := os.OpenFile(Helper.GetWorkDir()+"/src/resource/static.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				log.Fatal("load resource  error:", err)
			}
			err = tmpl.Execute(file, res)
			if err != nil {
				log.Fatal("create resource  error:", err)
			}
			log.Println("资源文件刷新成功")
		}
	},
}

func init() {
	rootCmd.AddCommand(resourceCmd)
}
