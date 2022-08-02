package routing

import (
	"log"

	"git.dsrt-int.net/actionmc/actionmc-site-go/gallery_photos"

	"git.dsrt-int.net/actionmc/actionmc-site-go/dynamic_cfg"

	"git.dsrt-int.net/actionmc/actionmc-site-go/auth"
	"git.dsrt-int.net/actionmc/actionmc-site-go/authdatabase"
	"git.dsrt-int.net/actionmc/actionmc-site-go/config"
	"git.dsrt-int.net/actionmc/actionmc-site-go/gdpr"
	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	"git.dsrt-int.net/actionmc/actionmc-site-go/oauth_handler"
	"git.dsrt-int.net/actionmc/actionmc-site-go/sessions"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template"
	"github.com/nu7hatch/gouuid"

	/*"github.com/gofiber/fiber/v2"
	"github.com/gofiber/redirect/v2"*/
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	//"golang.org/x/oauth2"
)

type Route struct {
	Type string
	Exec func(http.ResponseWriter, *http.Request)
}

func static_PageRender(template_name string, session *sessions.Session, templateCollection *template.TemplateCollection) (baseTemp string) {
	baseTemp, renderTime := templateCollection.BuildPopulatedTemplate(template_name, session)
	baseTemp = strings.Replace(baseTemp, "%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0), 1)

	return
}

func gdpr_static_PageRender(template_name string, templateCollection *template.TemplateCollection) (baseTemp string) {
	baseTemp, renderTime := templateCollection.BuildPopulatedGDPRTemplate(template_name)
	baseTemp = strings.Replace(baseTemp, "%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0), 1)

	return
}

func gdpr_nobanner_static_PageRender(template_name string, templateCollection *template.TemplateCollection) (baseTemp string) {
	baseTemp, renderTime := templateCollection.BuildPopulatedGDPR_NoBanner_Template(template_name)
	baseTemp = strings.Replace(baseTemp, "%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0), 1)

	return
}

type RoutingMemory struct {
	database           *authdatabase.MCAuthDB_sqlite3
	gallery_photos     *gallery_photos.GalleryPhotoHandler
	oauth2_struct      *oauth_handler.OAuthHandler
	store              *sessions.Store
	state              string
	templateCollection *template.TemplateCollection
	logger             *logging.Logger
	site_configuration *config.Configuration
	consent_handler    *gdpr.ConsentHandler
	dynamic_config     *dynamic_cfg.DynamicConfig
}

func GenState() string {
	u, err := uuid.NewV4()

	if err != nil {
		log.Panic(err)
	}

	u_two, err := uuid.NewV4()

	if err != nil {
		log.Panic(err)
	}

	return u.String() + "-" + u_two.String()
}

func Routing_Init(logger *logging.Logger) (routes map[string]Route, muxxer *mux.Router) {
	// basic config initialisation
	config.InitialiseConfig()

	// initialise everything needed for routing here, shaves down lines in main func.
	database := authdatabase.InitSQLite3DB(logger)
	gallery_photos := gallery_photos.NewGalleryPhotos(database)
	oauth2_struct := oauth_handler.New()
	store := sessions.New()
	muxxer = mux.NewRouter()
	state := GenState()
	templateCollection := template.NewTemplateCollection("data/templates", "data/bases", logger)
	site_configuration := config.InitialiseConfig()
	consent_handler := gdpr.NewConsentHandler(database)
	dynamic_config := dynamic_cfg.InitialiseDynConfig()

	procmem := &RoutingMemory{
		database:           database,
		gallery_photos:     gallery_photos,
		oauth2_struct:      oauth2_struct,
		store:              store,
		state:              state,
		templateCollection: templateCollection,
		logger:             logger,
		site_configuration: site_configuration,
		consent_handler:    consent_handler,
		dynamic_config:     dynamic_config,
	}

	// Runtime generate hashmap of routes based on function names
	// _0_ defines a new directory
	// _1_ defines a period
	// Prefix with GET to define a GET route, POST for a POST route
	// GET_0_ = GET '/'
	// GET_0_data_1_jpg = GET '/data.jpg'
	// POST_0_user = POST '/user'

	routes = map[string]Route{}

	typ := reflect.TypeOf(procmem)
	ref := reflect.ValueOf(procmem)

	methods := ref.NumMethod()
	syntaxreplacer := strings.NewReplacer("_0_", "/", "_1_", ".")

	// Valid routes are defined here as a hashmap
	// Add more as needed
	validroutes := map[string]struct{}{
		"GET":  {},
		"POST": {},
	}

	for i := 0; i < methods; i++ {
		// Get method name
		f := typ.Method(i).Name

		potroute := strings.Split(f, "_0_")

		// might panic if len(f) == 0, but a function name should never be len 0 so whatever
		if _, ok := validroutes[potroute[0]]; ok {
			routeTyp := potroute[0]

			// Replace syntax of _n_ values
			routedir := syntaxreplacer.Replace(strings.TrimPrefix(f, routeTyp))

			// Test if function is valid args
			fp, ok := ref.MethodByName(f).Interface().(func(http.ResponseWriter, *http.Request))

			if ok {
				if strings.Contains(routedir, "admin") || strings.Contains(routedir, "login") || strings.Contains(routedir, "maintenance") || strings.Contains(routedir, "api") {
					// ADMIN DIRECTORY so this bypasses maintenance mode stuff.
					if !strings.Contains(routedir, "login_required") {
						routes[routedir] = Route{routeTyp, fp}
					} else {
						routes[routedir] = Route{routeTyp, func(a http.ResponseWriter, b *http.Request) {
							// all the stuff here runs before the site loads
							if procmem.dynamic_config.MaintenanceMode {
								http.Redirect(a, b, "/maintenance", 303)
								return
							}
							fp(a, b)
						}}
					}

				} else {
					routes[routedir] = Route{routeTyp, func(a http.ResponseWriter, b *http.Request) {
						// all the stuff here runs before the site loads
						if procmem.dynamic_config.MaintenanceMode {
							http.Redirect(a, b, "/maintenance", 303)
							return
						}
						fp(a, b)
					}}
				}

			} else {
				panic("RoutingMemory." + f + " Not valid function")
			}
		}
	}

	for key, element := range routes {
		muxxer.HandleFunc(key, element.Exec)
	}

	muxxer.Path("/api/state/{username}").Handler(http.HandlerFunc(procmem.API_State_Handler))
	muxxer.Path("/api/v2/state/{username}").Handler(http.HandlerFunc(procmem.API_v2_State_Handler))
	muxxer.Path("/MinecraftSkins/{username}").Handler(http.HandlerFunc(procmem.Skin_GET_Func))
	muxxer.Path("/MinecraftCloaks/{username}").Handler(http.HandlerFunc(procmem.Cape_GET_Func))

	/*fs := http.FileServer(http.Dir("../static/"))
	  muxxer.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))*/

	//muxxer.NotFoundHandler = muxxer.NewRoute().HandlerFunc(procmem.NotFoundHandler).GetHandler()

	/*
		// redirection middleware
		app.Use(redirect.New(redirect.Config{
			Rules: map[string]string{
				"/login":       conf.AuthCodeURL(state),
				"/login/error": "/error?type=oauth2_denied",
				"/favicon.ico": "https://cdn.actionmc.ml/site/favicons/favicon.ico",
			},
			StatusCode: 303,
		}))

		// preproduction logging middleware
		app.Use(func(c *fiber.Ctx) error {
			// Log.

			logger.Debug.Println(string(c.Request().Header.Method()) + " " + c.OriginalURL())

			return c.Next()
		})
	*/

	return
}

func check_authentication(w http.ResponseWriter, r *http.Request, sess *sessions.Session) bool {
	login_check, account_check := auth.CheckAuthentication(sess)

	if !login_check {
		//RM.logger.Debug.Println("Poop")
		http.Redirect(w, r, "/login_required", 302)
		return false
	}

	if !account_check {
		//RM.logger.Debug.Println("Cum	")
		http.Redirect(w, r, "/register", 302)
		return false
	}

	return true
}
