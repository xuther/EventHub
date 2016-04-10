package main

import (
	"bytes"
	"fmt"
	"net/http"
)

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
				fireWebhoock(user.NotificationChannels[0], ev.occurance)
			}
		}
	}
}

//There's got to be a better way to do this. But I can't think of it without
//polymorphism
func fireWebhoock(chanToSend notificationChannel, occurance eventOccurance) {
	fmt.Println("Firing Event.")
	addr := chanToSend.Info[0]
	fmt.Printf("Sending hoock to %s\n", addr)

	//message to send. BASIC
	//TODO: make this more polished. Figure out what needs to be actually sent.
	http.Post(addr, "application/json", bytes.NewBufferString(occurance.EventInformation[0]))
}
