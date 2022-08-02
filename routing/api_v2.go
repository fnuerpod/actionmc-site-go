package routing

/*
	route list for implementation:

	/api/state/{username} GET
	/api/status GET
	/api/status/amc GET

	TODO: Implement API tokens. (could be breaking, migration project.)
*/

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"

	"git.dsrt-int.net/actionmc/actionmc-site-go/unsafety"

	"git.dsrt-int.net/actionmc/actionmc-site-go/template_data"
	"github.com/gorilla/mux"
)

// GET '/api/v2/authenticate'
func (RM *RoutingMemory) GET_0_api_0_v2_0_authenticate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("username") == "" {
		w.WriteHeader(400)
		w.Write([]byte("No username specified."))
		return
	}

	if r.FormValue("code") == "" {
		w.WriteHeader(400)
		w.Write([]byte("No TOTP auth code specified."))
		return
	}

	user, valid := RM.database.Getuser_byuname(r.FormValue("username"))

	if !valid {
		// not a valid user
		//w.WriteHeader(401)
		w.Write([]byte(r.FormValue("username") + ";0"))
		return
	}

	key, exists := RM.database.GetTOTPSecret(user.Uid)

	if !exists {
		w.Write([]byte(r.FormValue("username") + ";0"))
		return
	}

	totpTokena, err := totp.GenerateCode(key, time.Now())

	if err != nil {
		RM.logger.Err.Println(err)
		w.Write([]byte(r.FormValue("username") + ";0"))
		return
	}

	if r.FormValue("code") == totpTokena {
		// valid token!!!
		w.Write([]byte(r.FormValue("username") + ";1"))
	} else {
		// invalid token!!!
		w.Write([]byte(r.FormValue("username") + ";0"))
	}
}

// GET '/api/v2/status'
func (RM *RoutingMemory) GET_0_api_0_v2_0_status(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(strings.Join(final, "")))
	return
}

// GET '/api/v2/status/simple'
func (RM *RoutingMemory) GET_0_api_0_v2_0_status_0_simple(w http.ResponseWriter, r *http.Request) {
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

	w.Write([]byte(strings.Join(statTable, "\n")))
	return
}

// GET '/api/v2/maintenance'
func (RM *RoutingMemory) GET_0_api_0_v2_0_maintenance(w http.ResponseWriter, r *http.Request) {
	// use unsafety library to convert bool to uint8, then convert that to int then to string.
	// pls dont kill me.

	w.Write([]byte(strconv.Itoa(int(unsafety.BtoU8(RM.dynamic_config.MaintenanceMode)))))
	return
}

// GET '/api/v2/announcement'
func (RM *RoutingMemory) GET_0_api_0_v2_0_announcement(w http.ResponseWriter, r *http.Request) {
	// initialise text string.
	text := "0"

	if RM.dynamic_config.ShowAnnounceBanner {
		// fill with announcement banner.
		text = "1\n" + RM.dynamic_config.AnnounceBannerText
	}

	// send to caller.
	w.Write([]byte(text))
	return
}

// ActionMGMT endpoints

// GET '/api/v2/actionmgmt/latest.jar'
func (RM *RoutingMemory) GET_0_api_0_v2_0_actionmgmt_0_latest_1_jar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// check if passed DL ticket
	if r.FormValue("dl_ticket") == "" {
		// no dl ticket.
		w.WriteHeader(403)
		w.Write([]byte("Invalid Download Ticket."))
		return
	}

	valid, _, _ := RM.database.GetTicket(r.FormValue("dl_ticket"))

	if valid {
		if defskin, err := os.ReadFile("./ugc_store/latest.jar"); err == nil {
			// jar file exists.

			w.Header().Add("Content-Type", "application/java-archive")
			//w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

			w.Write(defskin)
		} else if errors.Is(err, os.ErrNotExist) {
			// jar file doesn't exist

			w.WriteHeader(404)

			//w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

			w.Write([]byte("JAR file not present on server."))
		} else {
			// ~~ schrodinger: skin may or may not exist. ~~
			// THIS MEANS THERE WAS AN IO ERROR AND WE SHOULD TREAT IT AS SUCH BTW
			// for safety, presume non-existance.
			w.WriteHeader(404)

			//w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

			w.Write([]byte("JAR file not present on server."))

		}
	} else {
		w.WriteHeader(403)
		w.Write([]byte("Invalid Download Ticket."))
	}

	return

}

// GET '/api/v2/actionmgmt/check'
func (RM *RoutingMemory) GET_0_api_0_v2_0_actionmgmt_0_check(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("version") == "" {
		// no version specified
		w.WriteHeader(400)
		w.Write([]byte("No version specified."))
	} else {
		// version specified.
		if r.FormValue("version") == RM.dynamic_config.ActionMGMT_LatestVersion {
			w.Write([]byte("OK"))
		} else {
			w.Write([]byte("FAIL"))
		}
	}

	return
}

// GET '/api/v2/actionmgmt/version'
func (RM *RoutingMemory) GET_0_api_0_v2_0_actionmgmt_0_version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(RM.dynamic_config.ActionMGMT_LatestVersion))

	return
}

// this route has to be handled differently as it uses a path variable.
// i'd change this to work like everything else in the handler but that
// would severely break the handlers for whitelist checking on the
// existing servers.

// GET '/api/v2/state/{username}'

func (RM *RoutingMemory) API_v2_State_Handler(w http.ResponseWriter, r *http.Request) {
	// get params from request.
	params := mux.Vars(r)

	var state_number string = "0"

	var gottenId string = "0"

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
			gottenId = id
			break
		}

	}

	_, privlevel := RM.database.Checkadmin(gottenId)

	w.Write([]byte(state_number + ";" + strconv.Itoa(privlevel)))
	return
}
