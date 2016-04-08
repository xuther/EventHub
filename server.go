package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/zenazn/goji"
)

var configuration config

func takeAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Taking Action")
	if checkCookie(r) {
		fmt.Fprintf(w, "Success")
	} else {
		fmt.Fprintf(w, "Yo")
	}
}

func checkCookie(r *http.Request) bool {
	fmt.Println("Checking Cookies")

	cookie, err := r.Cookie("Session")
	if err != nil || !checkSession(cookie.Value) {
		return false
	}

	return true
}

func getData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting all Data")

	if !checkCookie(r) {
		fmt.Fprintf(w, "notLoggedIn")
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logging In")

	decoder := json.NewDecoder(r.Body)

	var login loginInformation

	err := decoder.Decode(&login)
	check(err)

	User, value := checkCredentials(login.Username, login.Password)

	if value {
		//We need to add the sessionID to the DB and delete the old sessionKey
		//associated with that key, if there is one.

		//check if there is a current session associated with the user.
		session := checkSessionByUsername(User.ID.Hex())

		sessionID := generateSessionString(64)

		if session.SessionKey == "" {
			createSession(sessionID, User.ID.Hex())
		} else {
			removeSession(string(session.SessionKey))
			createSession(sessionID, User.ID.Hex())
		}
		fmt.Println("Session created, generating cookie...")

		newCookie := http.Cookie{Name: "Session", Value: sessionID}
		http.SetCookie(w, &newCookie)

		fmt.Println("Done.")
		fmt.Fprintf(w, "Success")

		fmt.Println("Login operation succeeded")
		return
	}
	fmt.Fprintf(w, "Failure")
}

//Import the configuration information from a JSON file
func importConfig(configPath string) config {
	fmt.Printf("Importing the configuration information from %v\n", configPath)

	f, err := ioutil.ReadFile(configPath)
	check(err)

	var c config
	json.Unmarshal(f, &c)

	fmt.Printf("\n%s\n", f)

	fmt.Printf("Done. Configuration data: \n %+v \n", c)
	return c
}

func registerHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering")

	decoder := json.NewDecoder(r.Body)

	var login loginInformation

	err := decoder.Decode(&login)
	check(err)

	//we should check to see if the username already exists
	if !checkForUsername(login.Username) {
		fmt.Fprintf(w, "Username Exists")
		return
	}
	register(login.Username, login.Password)

	fmt.Fprintf(w, "Success")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//var port = flag.Int("port", 443, "The port number you want the server running on. Default is 8080")
	var c = flag.String("cfg", "./config.json", "Location of the configuration file.")
	flag.Parse()

	configuration = importConfig(*c)

	//Front end actions
	http.Handle("/site", http.StripPrefix("/", http.FileServer(http.Dir("Static/"))))
	goji.Post("/api/login", login)
	goji.Post("/api/getAllData", getData)
	goji.Post("/api/register", registerHandle)

	//actions for event registry
	goji.Post("/api/provider", registerProvider)
	//goji.Get("/api/provider", getProviders)
	goji.Post("/api/provider/:providerID/events", registerEvent)

	goji.Serve()
	//err := http.ListenAndServeTLS(":"+strconv.Itoa(*port), "server.pem", "server.key", nil)

	//check(err)
}
