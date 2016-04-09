package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
)

func subscribeToEvent(c web.C, w http.ResponseWriter, r *http.Request) {
	provider := c.URLParams["providerID"]
	event := c.URLParams["eventID"]

	fmt.Printf("Subscribing to event: %s:%s\n", provider, event)

	//check to see if the event and provider exist.
	_, err := getEventByID(event, provider)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	var subInfo = subscription{}
	bits, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read request body.\n")
	}

	err = json.Unmarshal(bits, &subInfo)

	fmt.Printf("Values sent in: %+v\n", subInfo)

	if err != nil || strings.EqualFold(subInfo.Name, "") ||
		strings.EqualFold(subInfo.SubscriberID, "") ||
		strings.EqualFold(subInfo.ActionID, "") {
		example := subscription{Name: "Name. Friendly name of the subscription",
			SubscriberID: "Your User ID.",
			ActionID:     "Action ID of the registered action to your user. Call POST:/api/users/ to register. POST:/api/users/:userID/Actions to add an action."}

		w.WriteHeader(http.StatusBadRequest)
		stuff, _ := json.Marshal(example)
		fmt.Fprintf(w, "Invalid request body. Body must be json in form of: %s", stuff)
		return
	}

	//we have our event in subInfo. Now we just insert it into the array.
	err = insertEvent(event, provider, &subInfo)
}
