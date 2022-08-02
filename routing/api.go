package routing

/*
	route list for implementation:

	/api/state/{username} GET
	/api/status GET
	/api/status/amc GET

	TODO: Implement API tokens. (could be breaking, migration project.)
*/

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"git.dsrt-int.net/actionmc/actionmc-site-go/template_data"
	"github.com/gorilla/mux"
)

// GET '/api'
func (RM *RoutingMemory) GET_0_api(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if RM.site_configuration.IsPreproduction || RM.site_configuration.IsBeta {
		w.Write([]byte("https://b.actionmc.ml/api/v2"))
	} else {
		w.Write([]byte("https://actionmc.ml/api/v2"))
	}

	return
}

// GET '/api/status'
func (RM *RoutingMemory) GET_0_api_0_status(w http.ResponseWriter, r *http.Request) {
	var statTable []string

	allServices := RM.database.Getallservices()
	defer allServices.Close()

	for allServices.Next() {
		var serviceIP string
		var serviceName string
		var serviceState int

		err := allServices.Scan(&serviceIP, &serviceName, &serviceState)
		if err != nil {
			log.Fatal(err)
		}

		statTable = append(statTable, "{\"serviceIP\": \""+serviceIP+"\", \"serviceName\": \""+serviceName+"\", \"rawStateID\": \""+strconv.Itoa(serviceState)+"\"}")

	}

	final := []string{
		"[",
		strings.Join(statTable, ","),
		"]",
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, strings.Join(final, ""))
	return
}

// GET '/api/status/amc'
func (RM *RoutingMemory) GET_0_api_0_status_0_amc(w http.ResponseWriter, r *http.Request) {
	var statTable []string

	allServices := RM.database.Getallservices()
	defer allServices.Close()

	for allServices.Next() {
		var serviceIP string
		var serviceName string
		var serviceState int

		var selState string

		err := allServices.Scan(&serviceIP, &serviceName, &serviceState)
		if err != nil {
			log.Fatal(err)
		}

		err, _, selState = template_data.Get_ServiceStrings(serviceState)

		if err != nil {
			RM.logger.Debug.Println("An unknown service state was found.")
		}

		serviceString := []string{
			serviceIP,
			serviceName,
			selState,
			strconv.Itoa(serviceState),
		}

		statTable = append(statTable, strings.Join(serviceString, ":"))
	}

	fmt.Fprintf(w, strings.Join(statTable, "\n"))
	return
}

// this route has to be handled differently as it uses a path variable.
// i'd change this to work like everything else in the handler but that
// would severely break the handlers for whitelist checking on the
// existing servers.

// GET '/api/state/{username}'

func (RM *RoutingMemory) API_State_Handler(w http.ResponseWriter, r *http.Request) {
	// get params from request.
	params := mux.Vars(r)

	var state_number string = "0"

	allUsers := RM.database.Getall()
	defer allUsers.Close()

	for allUsers.Next() {
		var id string
		var userName string
		var state int
		var nameChangeTime int

		/*var selstate string
		var curstate string = "Unknown"*/

		err := allUsers.Scan(&id, &userName, &state, &nameChangeTime)
		if err != nil {
			log.Fatal(err)
		}

		if userName == params["username"] {
			// create their thingy

			state_number = strconv.Itoa(state)
			break
		}

	}

	fmt.Fprintf(w, state_number)
	return
}
