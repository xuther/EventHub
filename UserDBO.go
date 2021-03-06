package main

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func addUser(toAdd *user) error {
	fmt.Printf("Adding User %s...\n", toAdd.Name)

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Users")
	i := bson.NewObjectId()
	toAdd.ID = i

	c.Insert(toAdd)

	fmt.Println("Done")
	return nil
}

func saveCalendarToken(t *oauth2.Token, userID string) error {
	fmt.Printf("Saving user: %s token string.", userID)

	usr, err := getUserByID(userID)

	check(err)

	has := false

	for i := range usr.NotificationChannels {
		if strings.Compare(usr.NotificationChannels[i].NotificationType, "googlecalendar") == 0 {
			usr.NotificationChannels[i].GoogleCalendarToken = *t
			fmt.Printf("\n%s\n", "FOUND!")
			break
		}
	}

	if !has {
		calendarChan := notificationChannel{ID: bson.NewObjectId(), Name: "googleCalendar", Description: "Google Calendar Insertion", NotificationType: "googlecalendar", GoogleCalendarToken: *t}
		usr.NotificationChannels = append(usr.NotificationChannels, calendarChan)
	}

	err = updateUser(&usr)

	return err
}

func getUserByID(userID string) (user, error) {
	fmt.Printf("Getting the user corresponding to ID: %s \n", userID)

	var toReturn = user{}
	if !bson.IsObjectIdHex(userID) {
		return toReturn, errors.New("Invalid ID format")
	}
	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return toReturn, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Users")

	//TODO: Figure out how to do this in one step.
	count, err := c.FindId(bson.ObjectIdHex(userID)).Count()

	if err != nil {
		return toReturn, err
	}
	if count < 1 {
		return toReturn, errors.New("Invalid ID.")
	}
	c.FindId(bson.ObjectIdHex(userID)).One(&toReturn)

	return toReturn, nil
}

func getUsers() ([]user, error) {
	fmt.Printf("Getting all the users...\n")

	var u []user

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return u, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Users")

	c.Find(bson.M{}).All(&u)

	return u, nil
}

func getNotificationChannel(userID string, notificationChannelID string) (user, error) {
	var toReturn user

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return toReturn, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Users")

	c.Find(bson.M{"notificationchannels._id": bson.ObjectIdHex(notificationChannelID)}).Select(bson.M{"notificationchannels.$": 1}).One(&toReturn)

	return toReturn, nil
}

func updateUser(usr *user) error {
	fmt.Printf("Updating user %s...\n", usr.Name)

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Users")

	err = c.UpdateId(usr.ID, bson.M{"$set": &usr})

	if err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}
