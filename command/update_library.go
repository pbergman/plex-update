package command

import (
	"errors"
	"net/http"

	"github.com/pbergman/app"
	"github.com/pbergman/plex-update/application"
	"github.com/spf13/pflag"
)

func NewUpdateLibraryCommand() app.CommandInterface {
	return &updateLibraryCommand{
		Command: app.Command{
			Name:  "update-libraries",
			Usage: "[options] [--] (PLEX-TOKEN) (LIBRARY-ID)",
			Flags: new(pflag.FlagSet),
			Short: "Signal plex to update library",
			Long: `
Arguments:
    PLEX-TOKEN              The subject used for the certificate as a string (see "help plex-token")
    LIBRARY-ID              The library id for updating (see "list-libraries")

Options:
    --quiet                 Disable the application output
    --verbose (-v,-vv,-vvv) Increase the verbosity of application output
	--host                  The plex host (default http://127.0.0.1:32400)
	--force					Refresh the entire library
`,
		},
	}
}

type updateLibraryCommand struct {
	app.Command
}

func (l *updateLibraryCommand) Run(args []string, app *app.App) error {

	if len(args) != 2 {
		return errors.New("missing required arguments")
	}

	client := application.NewClient(args[0], app.Container.(*application.Container).GetLogger())

	uri := "/library/sections/" + args[1] + "/refresh"

	if l.isForce() {
		uri += "?force=1"
	}

	request, err := http.NewRequest("GET", l.getHost()+uri, nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(request)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func (l updateLibraryCommand) getHost() string {
	return l.Flags.(*pflag.FlagSet).Lookup("host").Value.String()
}

func (l updateLibraryCommand) isForce() bool {
	v, _ := l.Flags.(*pflag.FlagSet).GetBool("force")
	return v
}
func (l *updateLibraryCommand) Init(a *app.App) error {
	a.Container.(*application.Container).AddFlags(l.Flags.(*pflag.FlagSet))
	l.Flags.(*pflag.FlagSet).String("host", "http://127.0.0.1:32400", "")
	l.Flags.(*pflag.FlagSet).Bool("force", false, "")
	return nil
}
