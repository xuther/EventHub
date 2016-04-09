package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
	"gopkg.in/mgo.v2/bson"
)

func registerProvider(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Registering a provider.\n")

	var prov = provider{}
	bits, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read request body.\n")
	}

	err = json.Unmarshal(bits, &prov)

	if err != nil || strings.EqualFold(prov.Name, "") ||
		strings.EqualFold(prov.Type, "") || strings.EqualFold(prov.Description, "") {
		example := provider{Name: "Unique Name. Error will return if not unique.",
			Type:        "Type of events provided. e.g. Food/Sports",
			Description: "Description of updates. e.g. Greek and Go Food Truck Updates"}

		w.WriteHeader(http.StatusBadRequest)
		stuff, _ := json.Marshal(example)
		fmt.Fprintf(w, "Invalid request body. Body must be json in form of: %s", stuff)
		return
	}

	hasName, err := checkProviderNameInUse(prov.Name)

	checkSendError(w, err)

	if hasName {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Name already in use. get /api/providers for a list of occupied names.")
		return
	}

	prov.Secret, err = getSecret()

	checkSendError(w, err)

	var events []event
	prov.Events = events

	err = addProvider(&prov)

	checkSendError(w, err)

	fmt.Fprintf(w, "{\"Message\":\"Success\",\"Secret\":\"%s\",\"ID\":\"%s\"}", prov.Secret, prov.ID.Hex())
}

func registerEvent(c web.C, w http.ResponseWriter, r *http.Request) {
	providerID := c.URLParams["providerID"]

	fmt.Printf("Registering an event for %s...\n", providerID)
	//check for a valid providerID
	prov, err := getProviderByID(providerID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error occured: %s\n", err.Error())
		return
	}

	prov.ID = bson.ObjectIdHex(providerID)

	//TODO: check for duplicate events.
	var eventWrapper = registerEventInfo{}
	bits, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read request body.\n")
		return
	}

	err = json.Unmarshal(bits, &eventWrapper)

	eventToAdd := eventWrapper.Event

	fmt.Printf("Comparing string %s, with %s \n", eventWrapper.Secret, prov.Secret)

	if strings.Compare(eventWrapper.Secret, prov.Secret) != 0 {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error. Not Authorized.\n")
		return
	}

	if err != nil || strings.EqualFold(eventToAdd.Name, "") ||
		strings.EqualFold(eventToAdd.Description, "") {
		example := event{Name: "Unique Name. Error will return if not unique.",
			Description: "Description of the event. e.g. Everytime the menu updates. General Updates, etc."}
		exampleWrapper := registerEventInfo{Secret: "Secret received when created provider.", Event: example}
		w.WriteHeader(http.StatusBadRequest)
		stuff, _ := json.Marshal(exampleWrapper)
		fmt.Fprintf(w, "Invalid request body. Body must be json in form of: %s", stuff)
		return
	}

	eventToAdd.ID = bson.NewObjectId()
	prov.Events = append(prov.Events, eventToAdd)

	err = updateProvider(&prov)

	checkSendError(w, err)

	ev, _ := json.Marshal(prov.Events[len(prov.Events)-1])
	fmt.Fprintf(w, "%s", ev)

	return
}

//TODO: Find a better way to generate a secret that can be uniqe.
//At least investigate
func getSecret() (string, error) {
	b := make([]byte, 100)
	_, err := rand.Read(b)

	if err != nil {
		return "", errors.New("Issue with generating the secret.")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
