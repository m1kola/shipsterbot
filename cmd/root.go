package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "shipster",
	Short: "Shipster is a stupid bot app that helps you maintain your shopping lists",
}

// Execute starts the root command of our CLI
func Execute() {
	rootCmd.Execute()
}
