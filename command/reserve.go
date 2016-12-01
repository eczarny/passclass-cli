package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/eczarny/passclass-api/classpass"
	"github.com/eczarny/passclass-api/passclass"
	"github.com/wsxiaoys/terminal"
)

type desiredReservation struct {
	When    passclass.Timestamp
	VenueId int
}

func Reserve(cmd string, args []string, provider classpass.AuthTokenProvider) {
	n, err := strconv.Atoi(args[0])
	if err != nil && n <= 0 {
		terminal.Stderr.Colorf("@{r}First argument must be the number of days (a positive integer) from now to make reservations.").Nl()
		return
	}
	tuples := args[1:]
	if len(tuples)%2 != 0 {
		terminal.Stderr.Colorf("@{r}Expected an even number of desired reservation date/venue ID pairs.").Nl()
		return
	}
	desiredReservations := make([]desiredReservation, 0, len(tuples)/2)
	for idx, _ := range tuples {
		if idx%2 == 0 {
			starttime, _ := time.Parse("15:04", tuples[idx])
			venueId, _ := strconv.Atoi(tuples[idx+1])
			when := convert(starttime, n)
			desiredReservations = append(desiredReservations, desiredReservation{When: passclass.Timestamp{when}, VenueId: venueId})
		}
	}
	reserveClass(desiredReservations, provider)
}

func convert(t time.Time, n int) time.Time {
	year, month, day := time.Now().Date()
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	return time.Date(year, month, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location).AddDate(0, 0, n)
}

func reserveClass(desiredReservations []desiredReservation, provider classpass.AuthTokenProvider) {
	stdout := terminal.Stdout
	stderr := terminal.Stderr
	venues := classpass.GetVenues(provider)
	for priority, r := range desiredReservations {
		venue, ok := venues[r.VenueId]
		if !ok {
			stderr.Color("r").Print("Venue ").Color(".r").Print(r.VenueId).Color("r").Print(" does not exist; unable to make a reservation.").Reset().Nl()
			continue
		}
		schedule, ok := classpass.GetSchedule(r.When, r.VenueId, provider)[r.When.String()]
		if !ok {
			stderr.Colorf(fmt.Sprintf("@{r}Schedule for class on %s could not be found.", r.When)).Nl()
			continue
		}
		details := fmt.Sprintf("Reserving class for %s at %s (priority %d).", r.When, venue.Name, priority)
		stdout.Color(".g").Print(details).Reset().Nl()
		reservationId, err := classpass.PostReservation(schedule, provider)
		if err != nil || reservationId == 0 {
			stderr.Colorf("@{y}Unable to reserve class trying the next class (if available).").Nl()
		} else {
			stdout.Colorf(fmt.Sprintf("@{g}Reservation %d created for %s at %s!", reservationId, venue.Name, r.When)).Nl()
			return
		}
	}
	stderr.Colorf("@{r}Unable to reserve any classes; try again next time!").Nl()
}
