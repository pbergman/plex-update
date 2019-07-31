package command

import (
	"github.com/pbergman/app"
)

func NewPlexTokenCommand() app.CommandInterface {
	return &app.Command{
		Name:  "plex-token",
		Short: "Information about requiring the plex-token",
		Long: `
This application uses the public plex api for updating the libraries and by default the server
must be authenticated. That means we need the api key for communicating with the api. 

see:

	https://support.plex.tv/articles/201638786-plex-media-server-url-commands

and:

	https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/

why the token is needed.
`,
	}
}
