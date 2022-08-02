package routing

import (
	"embed"
	"errors"
	"net/http"
	"os"
	"strings"

	"path/filepath"

	"git.dsrt-int.net/actionmc/actionmc-site-go/config"

	"github.com/gorilla/mux"
)

//go:embed skin_defaults/*

var content embed.FS

func (RM *RoutingMemory) Skin_GET_Func(w http.ResponseWriter, r *http.Request) {
	// get params from request.
	params := mux.Vars(r)

	name := strings.TrimSuffix(params["username"], ".png")

	userinfo, exists := RM.database.Getuser_byuname(name)

	if !exists {
		// username doesn't exist on this skin server, so we'll want to return the default skin.
		defskin, err := content.ReadFile("skin_defaults/skin.png")

		if err != nil {
			RM.logger.Fatal.Fatalln(err)
		}

		w.Header().Add("Content-Type", "image/png")
		//w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

		w.Write(defskin)
		return
	}

	if defskin, err := os.ReadFile(filepath.Join(config.GetDataDir(), "ugc_store", "skin", userinfo.Uid+".png")); err == nil {
		// user has a skin on the server.

		w.Header().Add("Content-Type", "image/png")
		//w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

		w.Write(defskin)
	} else if errors.Is(err, os.ErrNotExist) {
		// user DOESN'T have a skin on the server.
		defskin, err := content.ReadFile("skin_defaults/skin.png")

		if err != nil {
			RM.logger.Fatal.Fatalln(err)
		}

		w.Header().Add("Content-Type", "image/png")
		//w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

		w.Write(defskin)
	} else {
		// ~~ schrodinger: skin may or may not exist. ~~
		// THIS MEANS THERE WAS AN IO ERROR AND WE SHOULD TREAT IT AS SUCH BTW
		// for safety, presume non-existance.
		defskin, err := content.ReadFile("skin_defaults/skin.png")

		if err != nil {
			RM.logger.Fatal.Fatalln(err)
		}

		w.Header().Add("Content-Type", "image/png")
		//w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

		w.Write(defskin)

	}

	return

}

func (RM *RoutingMemory) Cape_GET_Func(w http.ResponseWriter, r *http.Request) {
	// get params from request.
	params := mux.Vars(r)

	name := strings.TrimSuffix(params["username"], ".png")

	userinfo, exists := RM.database.Getuser_byuname(name)

	if !exists {
		// username doesn't exist on this skin server, so we'll want to return the default cape.
		RM.logger.Debug.Println("Username not known, serve default.")
		http.Error(w, "User doesn't have a cape.", http.StatusNotFound)
	}

	if defskin, err := os.ReadFile(filepath.Join(config.GetDataDir(), "ugc_store", "cape", userinfo.Uid+".png")); err == nil {
		// user has a cape on the server.

		w.Header().Add("Content-Type", "image/png")
		//w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

		w.Write(defskin)
	} else if errors.Is(err, os.ErrNotExist) {
		// user DOESN'T have a skin on the server.
		http.Error(w, "User doesn't have a cape.", http.StatusNotFound)

	} else {
		// schrodinger: skin may or may not exist.
		// THIS MEANS THERE WAS AN IO ERROR AND WE SHOULD TREAT IT AS SUCH BTW, LIKE A CRAZY EDGE CASE REALLY IT SHOULD PANIC
		// for safety, presume non-existance.
		http.Error(w, "User doesn't have a cape.", http.StatusNotFound)

	}

	return

}
