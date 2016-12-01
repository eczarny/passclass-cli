package command

import (
	"github.com/eczarny/passclass-api/classpass"
	"github.com/wsxiaoys/terminal"
)

func Venues(cmd string, args []string, provider classpass.AuthTokenProvider) {
	stdout := terminal.Stdout
	for _, venue := range classpass.GetVenues(provider) {
		stdout.Color(".").Print(venue.Id).Reset()
		stdout.Print("	").Print(venue.Name)
		if len(venue.Subtitle) > 0 {
			stdout.Print(", ").Print(venue.Subtitle)
		}
		stdout.Nl()
	}
}
