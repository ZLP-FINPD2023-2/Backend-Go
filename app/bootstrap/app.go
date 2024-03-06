package bootstrap

import (
	"github.com/spf13/cobra"

	"finapp/commands"
)

var rootCmd = &cobra.Command{
	Use:              "finapp",
	Short:            "Finapp REST API server",
	TraverseChildren: true,
}

// App root of application
type App struct {
	*cobra.Command
}

func NewApp() App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(commands.GetSubCommands(CommonModules)...)
	return cmd
}

var RootApp = NewApp()
