package main

import "gopkg.in/mgo.v2/bson"

type config struct {
	MongoDBAddress      string
	MongoProviderDBName string
}

type publicProvider struct {
	Name        string
	Type        string
	Description string
}

type event struct {
	ID          bson.ObjectId `bson:"_id" json:"omitempty"`
	Name        string
	Description string
	Subscribers []subscription   `json:"-"`
	History     []eventOccurance `json:"-"`
}

type subscription struct {
	ID                    bson.ObjectId `bson:"_id" json:"-"`
	Name                  string
	SubscriberID          string //Hex strings of the mongo ID's
	NotificationChannelID string //Hex Strings of the mongo ID's
}

type loginInformation struct {
	Username string
	Password string
	Name     string
}

type user struct {
	ID                   bson.ObjectId `bson:"_id"`
	Name                 string
	NotificationChannels []notificationChannel
	Username             string
	Password             []byte `json:"-"`
	Secret               string
}

type eventOccurance struct {
}

type notificationChannel struct {
	ID               bson.ObjectId `bson:"_id"`
	Name             string
	Description      string
	NotificationType string
	Info             []string //Depending on the action this will be a phone #, e-mail address, or webhook to call, or other information.
}

type registerNotificationChannel struct {
	Secret  string
	Channel notificationChannel
}

type registerEventInfo struct {
	Event  event
	Secret string
}

//TODO Make ID go away in json
type provider struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"omitempty"`
	Name        string
	Type        string
	Description string
	Secret      string  `json:"omitempty"`
	Events      []event `json:"omitempty"`
}
