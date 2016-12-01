package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/eczarny/passclass-api/classpass"
	"github.com/eczarny/passclass-api/passclass"
	"github.com/wsxiaoys/terminal"
)

func Schedule(cmd string, args []string, provider classpass.AuthTokenProvider) {
	venueId, _ := strconv.Atoi(args[0])
	when, _ := time.Parse(time.RFC1123, args[1])
	stdout := terminal.Stdout
	for _, schedule := range classpass.GetSchedule(passclass.Timestamp{when}, venueId, provider) {
		details := fmt.Sprintf("%s %s, %s", schedule.Starttime, schedule.Venue.Name, schedule.Class.Name)
		stdout.Color(".").Print(schedule.Id).Reset()
		stdout.Print(" ").Print(details)
		if len(schedule.TeacherName) > 0 {
			stdout.Print(" ").Print("(Instructor: ").Print(schedule.TeacherName).Print(")")
		}
		stdout.Nl()
	}
}
