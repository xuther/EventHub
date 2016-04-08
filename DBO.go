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
