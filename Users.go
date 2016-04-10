package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/zenazn/goji/web"
	"golang.org/x/crypto/bcrypt"
)

func registerUser(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering User....")

	var loginInfo = loginInformation{}
	bits, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read request body.\n")
		return
	}

	err = json.Unmarshal(bits, &loginInfo)

	if err != nil || strings.EqualFold(loginInfo.Name, "") ||
		strings.EqualFold(loginInfo.Username, "") ||
		strings.EqualFold(loginInfo.Password, "") {
		example := loginInformation{Name: "Your Name",
			Username: "Your Username",
			Password: "Password"}

		w.WriteHeader(http.StatusBadRequest)
		stuff, _ := json.Marshal(example)
		fmt.Fprintf(w, "Invalid request body. Body must be json in form of: %s", stuff)
		return
	}

	secret, err := getSecret()
	checkSendError(w, err)

	pass, err := bcrypt.GenerateFromPassword([]byte(loginInfo.Password), 10)
	checkSendError(w, err)

	toAdd := user{Secret: secret, Username: loginInfo.Username, Password: pass, Name: loginInfo.Name}

	err = addUser(&toAdd)
	checkSendError(w, err)

	s, err := json.Marshal(toAdd)
	checkSendError(w, err)

	fmt.Fprintf(w, "%s", s)
}

func addUserNotificationChannel(c web.C, w http.ResponseWriter, r *http.Request) {
	userID := c.URLParams["userID"]

	usr, err := getUserByID(userID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error occured: %s\n", err.Error())
		return
	}

	fmt.Printf("Registering a notification channel for %s...\n", usr.Name)

	fmt.Printf("User found %+v\n", usr)

	var eventWrapper = registerNotificationChannel{}
	bits, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read request body.\n")
		return
	}

	err = json.Unmarshal(bits, &eventWrapper)
	checkSendError(w, err)

	channelToAdd := eventWrapper.Channel

	fmt.Printf("Comparing strings %s with %s\n", eventWrapper.Secret, usr.Secret)

	if strings.Compare(eventWrapper.Secret, usr.Secret) != 0 {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error. Not Authorized.\n")
		return
	}

	//TODO: Make this return a better, understandable help document.
	if strings.EqualFold(channelToAdd.Name, "") ||
		strings.EqualFold(channelToAdd.Description, "") {
		example := event{Name: "Unique Name. Error will return if not unique.",
			Description: "Description of the event. e.g. Everytime the menu updates. General Updates, etc."}
		exampleWrapper := registerEventInfo{Secret: "Secret received when created provider.", Event: example}
		w.WriteHeader(http.StatusBadRequest)
		stuff, _ := json.Marshal(exampleWrapper)
		fmt.Fprintf(w, "Invalid request body. Body must be json in form of: \n %s", stuff)
		return
	}

	channelToAdd.ID = bson.NewObjectId()
	usr.NotificationChannels = append(usr.NotificationChannels, channelToAdd)

	err = updateUser(&usr)

	checkSendError(w, err)

	ev, _ := json.Marshal(usr.NotificationChannels[len(usr.NotificationChannels)-1])
	fmt.Fprintf(w, "%s", ev)

	return
}

func getAllUsers(c web.C, w http.ResponseWriter, r *http.Request) {
	users, err := getUsers()

	checkSendError(w, err)

	for indx := range users {
		users[indx].Secret = ""
	}

	s, err := json.Marshal(users)

	checkSendError(w, err)

	w.Write(s)
}

func checkSendError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
}
