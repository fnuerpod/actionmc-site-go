package routing

import (

	//"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// GET '/nocookie/skin_server'
func (RM *RoutingMemory) GET_0_nocookie_0_skin_server(w http.ResponseWriter, r *http.Request) {
	baseTemp := gdpr_static_PageRender("skin-server-nc", RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/privacy'
func (RM *RoutingMemory) GET_0_nocookie_0_privacy(w http.ResponseWriter, r *http.Request) {
	baseTemp := gdpr_static_PageRender("privacy", RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/nocookie/leaving'
func (RM *RoutingMemory) GET_0_nocookie_0_leaving(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	baseTemp := gdpr_static_PageRender("leaving", RM.templateCollection)

	baseTemp = strings.Replace(baseTemp, "%%URL%%", r.FormValue("url"), -1)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}
