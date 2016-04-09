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
	ID           bson.ObjectId `bson:"_id" json:"-"`
	Name         string
	SubscriberID string //Hex strings of the mongo ID's
	ActionID     string //Hex Strings of the mongo ID's
}

type subscriberAction struct {
	ID          bson.ObjectId `bson:"_id"`
	Name        string
	Description string
	Action      string
	Info        []string //Depending on the action this will be a phone #, e-mail address, or webhook to call, or other information.
}

type user struct {
	ID            bson.ObjectId `bson:"_id"`
	Name          string
	Subscriptions []subscriberAction
	Username      string
	Password      []byte
}

type eventOccurance struct {
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
