package main

import (
	"os"

	"github.com/pbergman/app"
	"github.com/pbergman/plex-update/application"
	"github.com/pbergman/plex-update/command"
)

func main() {

	container := new(application.Container)
	container.SetApp(newApp(container))

	// catch no args
	if len(container.GetArgs()) < 1 {
		container.GetFlags().Usage()
	}

	if err := container.GetApp().Run(container.GetArgs()); err != nil {
		container.GetLogger().Error(err)
		os.Exit(int(err.(*app.Error).Code))
	}
}

func newApp(c *application.Container) *app.App {

	application := app.NewApp(
		command.NewPlexTokenCommand(),
		command.NewListLibrariesCommand(),
		command.NewUpdateLibraryCommand(),
	)

	application.Container = c

	application.PreRun = func(cmd app.CommandInterface) error {
		if val, err := c.GetFlags().GetBool("help"); err == nil && val {
			application.PrintHelp(cmd)
			application.Active = false
		}
		c.SetCurrent(cmd)
		return nil
	}

	return application
}
