package main

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//TODO: Fix this so it checks if the provider name is in use.
func checkProviderNameInUse(name string) (bool, error) {
	return false, nil
}

func addProvider(prov *provider) error {
	fmt.Printf("Adding provider %s...\n", prov.Name)

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Providers")
	i := bson.NewObjectId()
	prov.ID = i

	c.Insert(prov)

	fmt.Println("Done")
	return nil
}

func removeEvent(eventID string) error {
	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Providers")

	c.Update(bson.M{"events._id": bson.ObjectIdHex(eventID)}, bson.M{"$pull": bson.M{"events": bson.M{"_id": bson.ObjectIdHex(eventID)}}})

	return nil
}

func getEventByID(eventID string, providerID string) (event, error) {
	fmt.Printf("\nGetting the event corresponding to ID: %s:%s \n", providerID, eventID)

	var toReturn = event{}
	var prov = provider{}

	if !bson.IsObjectIdHex(providerID) || !bson.IsObjectIdHex(eventID) {
		return toReturn, errors.New("Invalid ID format\n")
	}
	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return toReturn, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Providers")

	if err != nil {
		return toReturn, err
	}

	//TODO: Figure out how to do this in one step.
	count, err := c.Find(bson.M{"_id": bson.ObjectIdHex(providerID), "events._id": bson.ObjectIdHex(eventID)}).Count()

	if err != nil {
		return toReturn, err
	}
	if count < 1 {
		return toReturn, errors.New("Invalid ID.")
	}
	c.Find(bson.M{"_id": bson.ObjectIdHex(providerID), "events._id": bson.ObjectIdHex(eventID)}).Select(bson.M{"events.$": 1}).One(&prov)

	return prov.Events[0], nil
}

func getProviderByID(providerID string) (provider, error) {
	fmt.Printf("Getting the provider corresponding to ID: %s \n", providerID)

	var toReturn = provider{}
	if !bson.IsObjectIdHex(providerID) {
		return toReturn, errors.New("Invalid ID format")
	}
	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return toReturn, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Providers")

	//TODO: Figure out how to do this in one step.
	count, err := c.FindId(bson.ObjectIdHex(providerID)).Count()

	if err != nil {
		return toReturn, err
	}
	if count < 1 {
		return toReturn, errors.New("Invalid ID.")
	}
	c.FindId(bson.ObjectIdHex(providerID)).One(&toReturn)

	return toReturn, nil
}

func updateProvider(prov *provider) error {
	fmt.Printf("Updating provider %s...\n", prov.Name)

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Providers")

	err = c.UpdateId(prov.ID, bson.M{"$set": &prov})

	if err != nil {
		return err
	}

	fmt.Println("Done")
	return nil
}

func getSubscriptions(eventID string, providerID string) ([]subscription, error) {

	var subs []subscription
	fmt.Printf("Getting all the subscriptions of %s...\n", eventID)

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return subs, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Providers")

	c.Find(bson.M{"events._id": bson.ObjectIdHex(eventID)}).All(&subs)

	return subs, nil
}

func getProviders() ([]provider, error) {
	fmt.Printf("Getting all the users...\n")

	var p []provider

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return p, err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Providers")

	c.Find(bson.M{}).All(&p)

	fmt.Printf("\nProvider:\n%+v", p)

	return p, nil
}

func insertEvent(event string, provider string, sub *subscription) error {
	fmt.Printf("Inserting a subscription provider %s...\n", sub.Name)

	session, err := mgo.Dial(configuration.MongoDBAddress)

	if err != nil {
		return err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(configuration.MongoProviderDBName).C("Providers")

	sub.ID = bson.NewObjectId()

	c.Update(bson.M{"events._id": bson.ObjectIdHex(event)}, bson.M{"$addToSet": bson.M{"events.$.subscribers": sub}})
	return nil
}
