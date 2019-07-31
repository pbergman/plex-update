package command

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/pbergman/app"
	"github.com/pbergman/plex-update/application"
	"github.com/spf13/pflag"
)

func NewListLibrariesCommand() app.CommandInterface {
	return &listLibrariesCommand{
		Command: app.Command{
			Name:  "list-libraries",
			Usage: "[options] [--] (PLEX-TOKEN)",
			Flags: new(pflag.FlagSet),
			Short: "List available libraries",
			Long: `
Print an table of all available libraries and from every library the path nd last updated time.

Arguments:
    PLEX-TOKEN              The subject used for the certificate as a string (see "help plex-token")

Options:
    --quiet                 Disable the application output
    --verbose (-v,-vv,-vvv) Increase the verbosity of application output
	--host                  The plex host (default http://127.0.0.1:32400) 
`,
		},
	}
}

type listLibrariesCommand struct {
	app.Command
}

func (l *listLibrariesCommand) Run(args []string, app *app.App) error {

	if len(args) != 1 {
		return errors.New("missing required PLEX-TOKEN argument (see \"help list-libraries\")")
	}

	client := application.NewClient(args[0], app.Container.(*application.Container).GetLogger())

	request, err := http.NewRequest("GET", l.getHost()+"/library/sections", nil)

	if err != nil {
		return err
	}

	resp, err := client.Do(request)

	if err != nil {
		return err
	}

	var mediaContainer struct {
		Directory []struct {
			Type     string `xml:"type,attr"`
			Title    string `xml:"title,attr"`
			Key      int    `xml:"key,attr"`
			UpdateAt int64  `xml:"updatedAt,attr"`
			Location struct {
				Name string `xml:"path,attr"`
			} `xml:"Location"`
		} `xml:"Directory"`
	}

	defer resp.Body.Close()

	if err := xml.NewDecoder(resp.Body).Decode(&mediaContainer); err != nil {
		return err
	}

	writer := tabwriter.NewWriter(os.Stdout, 2, 8, 2, '\t', 0)
	_, _ = writer.Write([]byte("key\ttype\ttitle\tlocation\tlast updated\n"))

	for _, dir := range mediaContainer.Directory {
		_, _ = fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\n", dir.Key, dir.Type, dir.Title, dir.Location.Name, time.Unix(dir.UpdateAt, 0).Format(time.RFC3339))
	}

	_ = writer.Flush()

	return nil
}

func (l listLibrariesCommand) getHost() string {
	return l.Flags.(*pflag.FlagSet).Lookup("host").Value.String()
}

func (l *listLibrariesCommand) Init(a *app.App) error {
	a.Container.(*application.Container).AddFlags(l.Flags.(*pflag.FlagSet))
	l.Flags.(*pflag.FlagSet).String("host", "http://127.0.0.1:32400", "")
	return nil
}
