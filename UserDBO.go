package main

import (
	"errors"
	"fmt"

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
