package cmd

import (
	"github.com/spf13/cobra"

	"github.com/m1kola/shipsterbot/internal/pkg/shipsterbot/cmd/telegram"
)

var rootCmd = &cobra.Command{
	Use:   "shipsterbot",
	Short: "Shipster is a stupid bot app that helps you maintain your shopping lists",
}

// Execute starts the root command of our CLI
func Execute() {
	var startBotCmd = &cobra.Command{
		Use:   "startbot",
		Short: "Start a bot",
	}

	rootCmd.AddCommand(startBotCmd)
	startBotCmd.AddCommand(telegram.NewStartTelegramBotCmd())

	rootCmd.Execute()
}
