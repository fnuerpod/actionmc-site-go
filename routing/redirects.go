package routing

import (
	"net/http"
)

// GET '/login'
func (RM *RoutingMemory) GET_0_login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, RM.oauth2_struct.GetConf().AuthCodeURL(RM.state), 303)
	return
}

// GET '/login/error'
func (RM *RoutingMemory) GET_0_login_0_error(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/error?type=oauth2_denied", 303)
	return
}

// GET '/favicon.ico'
func (RM *RoutingMemory) GET_0_favicon_1_ico(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/favicon.ico", 303)
	return
}
