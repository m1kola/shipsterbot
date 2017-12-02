package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "shipster",
	Short: "Shipster is a stupid bot app that helps you maintain your shopping lists",
}

func Execute() {
	rootCmd.Execute()
}
