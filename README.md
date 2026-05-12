# Gofluxer Demo Bot

Gofluxer Package TL;DR: Gofluxer is an API wrapper for the Go programming language to create fluxer.app bots with.

This repo is the demonstration bot for the Gofluxer package to showcase some of the stuff that you can do with the Gofluxer package. This repo even shows how easy it is to make a simple Fluxer bot with commands with our Gofluxer package.

Please ensure you have Go 1.21 or newer before trying out this project.

## Setup Guide:

1. Fork/Download this source code
2. Create a Fluxer.app bot in the Applications section through your user settings.
3. Fill in everything in the config.json file:

```json
{
	"FLUXERBOTTOKEN": "", #Your Fluxer.app bot token
	"FLUXERBOTPREFIX": "!", #Prefix for the Fluxer.app bot
	"FLUXERBASEURL": "", #The base API url if the Fluxer instance (UNUSED)
	"FLUXERGATEWAYURL": "" #The websocket base url (UNUSED)
}
```

4. Run the following commands in your console (Linux):

```sh
go get github.com/go-fluxer/gofluxer
```

5. Run the bot (Command is: "go run .")