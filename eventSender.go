package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/calendar/v3"

	"golang.org/x/net/context"
)

//There's got to be a better way to do this. But I can't think of it without
//polymorphism
func eventHandler(events <-chan eventFireInformation) {
	fmt.Println("Starting event Handler...")

	//run and wait for events, when we
	for true {
		ev := <-events
		fmt.Printf("Handling event %+v \n", ev)
		//messy, but just get something through.
		for subNo := range ev.Subscriptions {
			fmt.Printf("Handling event. %v \n", subNo)
			cur := ev.Subscriptions[subNo]
			fmt.Printf("Current : %+v, \n", cur)
			user, err := getNotificationChannel(cur.SubscriberID, cur.NotificationChannelID)

			if err != nil {
				fmt.Printf("There was an error \n")
				//TODO:log error
				panic(err)
			}
			fmt.Printf("User found: %+v\n", user)
			switch user.NotificationChannels[0].NotificationType {
			case "webhook":
				fireWebhook(user.NotificationChannels[0], ev.occurance)
			case "googlecalendar":
				addToGoogleCalendar(user.NotificationChannels[0], ev.occurance)
			}
		}
	}
}

func addToGoogleCalendar(chanToSend notificationChannel, occurance eventOccurance) {
	fmt.Println("Adding event to calendar")

	fmt.Printf("What we get in the event occurance: ")
	config := getConfig()
	token := chanToSend.GoogleCalendarToken
	ctx := context.Background()

	client := config.Client(ctx, &token)

	srv, err := calendar.New(client)

	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}

	event := &calendar.Event{
		Summary:     occurance.EventInformation[0],
		Description: occurance.EventInformation[1],
		Start: &calendar.EventDateTime{
			DateTime: occurance.EventInformation[2],
			TimeZone: "America/Denver",
		},
		End: &calendar.EventDateTime{
			DateTime: occurance.EventInformation[3],
			TimeZone: "America/Denver",
		},
		Location: occurance.EventInformation[4],
	}

	event, err = srv.Events.Insert("primary", event).Do()

	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
}

func fireWebhook(chanToSend notificationChannel, occurance eventOccurance) {
	fmt.Println("Firing Event.")
	addr := chanToSend.Info[0]
	fmt.Printf("Sending hoock to %s\n", addr)
	bits, _ := json.Marshal(occurance.EventInformation)
	//message to send. BASIC
	//TODO: make this more polished. Figure out what needs to be actually sent.
	http.Post(addr, "application/json", bytes.NewBuffer(bits))
}
