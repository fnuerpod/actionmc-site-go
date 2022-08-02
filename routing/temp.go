package routing

import (
	"git.dsrt-int.net/actionmc/actionmc-site-go/auth"

	//"git.dsrt-int.net/actionmc/actionmc-site-go/oauth_handler"

	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// GET '/retrospective2021'
func (RM *RoutingMemory) GET_0_retrospective2021(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO NOCOOKIE RETROSPECTIVE.
		http.Redirect(w, r, "/nocookie/retrospective2021", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("retrospective2021", sess, RM.templateCollection)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/nocookie/retrospective2021'
func (RM *RoutingMemory) GET_0_nocookie_0_retrospective2021(w http.ResponseWriter, r *http.Request) {
	baseTemp := gdpr_static_PageRender("retrospective2021", RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}
