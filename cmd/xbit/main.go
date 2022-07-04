package main

import (
	"embed"
	"github.com/xbitgo/core/di"
	"log"

	"github.com/spf13/cobra"

	"xbit/conf"
	"xbit/handler"
)

var rootCmd = &cobra.Command{
	Use:     "xbit",
	Short:   "xbit: An toolkit for xbitgo.",
	Long:    "xbit: An toolkit for xbitgo.",
	Version: "0.0.1",
}

func init() {
	conf.Init()
	rootCmd.AddCommand(handler.CmdList()...)
}

//go:embed template
var projectTpl embed.FS

func main() {
	di.Register("project_tpl", projectTpl)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
