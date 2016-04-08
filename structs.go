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
	ID          int
	Name        string
	Description string
	Subscribers []bson.ObjectId  `json:"-"`
	History     []eventOccurance `json:"-"`
}

type eventOccurance struct {
}

type registerEventInfo struct {
	Event  event
	Secret string
}

//TODO Make ID go away in json
type provider struct {
	ID             bson.ObjectId `bson:"_id" json:"omitempty"`
	EventIDCounter int           `json:"-"`
	Name           string
	Type           string
	Description    string
	Secret         string  `json:"omitempty"`
	Events         []event `json:"omitempty"`
}
