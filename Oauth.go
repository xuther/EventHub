package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zenazn/goji/web"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(user string, config *oauth2.Config) string {
	authURL := config.AuthCodeURL(user, oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	return authURL
}

func getOauthToken(c web.C, w http.ResponseWriter, r *http.Request) {

	user := c.URLParams["userID"]

	config := getConfig()

	toReturn := getTokenFromWeb(user, config)

	fmt.Fprintf(w, "%v", toReturn)
}

func getConfig() *oauth2.Config {
	//user := c.URLParams["userID"]

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)

	check(err)

	return config
}

func saveOauthToken(c web.C, w http.ResponseWriter, r *http.Request) {
	//userID is the state
	userID := r.URL.Query().Get("state")
	//token is the code
	tokenCode := r.URL.Query().Get("code")

	config := getConfig()

	ctx := context.Background()

	token, err := config.Exchange(ctx, tokenCode)

	check(err)

	fmt.Printf("Token:\n %+v\n", token)

	err = saveCalendarToken(token, userID)

	//fmt.Printf("%s", err.Error())

	fmt.Fprintf(w, "UserID: %s, Token Code: %s", userID, tokenCode)
}
