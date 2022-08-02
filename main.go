package main

// imports
import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	_ "net/http/pprof"

	"git.dsrt-int.net/actionmc/actionmc-site-go/config"
	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	"git.dsrt-int.net/actionmc/actionmc-site-go/routing"
	"github.com/gabriel-vasile/mimetype"
)

//go:embed static/*

var content embed.FS

func main() {
	// instantise functions for site.
	// initialises database, oauth2 config, otp handler, session storage and the fiber app.
	// routing initialisation will also handle the addition of middlewares to the base app.

	site_configuration := config.InitialiseConfig()

	logger := logging.New()

	_, routing_mux := routing.Routing_Init(logger)

	// start listening
	logger.Log.Println("ActionMC website starting on port " + site_configuration.BindPort + "...")
	//app.Listen(":" + constants.BindPort)

	routing_mux.PathPrefix("/static/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := strings.TrimPrefix(r.RequestURI, "/static/")

		content, err := fs.ReadFile(content, "static/"+uri)

		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("File not found."))
		}

		if strings.HasSuffix(uri, "css") {
			w.Header().Add("Content-Type", "text/css")
		} else {
			w.Header().Add("Content-Type", mimetype.Detect(content).String())
		}

		w.Write(content)
	}))

	routing_mux.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

	err := http.ListenAndServe(":"+site_configuration.BindPort, routing_mux)
	logger.Fatal.Fatalln(err)

}

//EOF
