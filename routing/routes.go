package routing

import (
	"math/rand"
	"strconv"

	"git.dsrt-int.net/actionmc/actionmc-site-go/auth"

	//"git.dsrt-int.net/actionmc/actionmc-site-go/oauth_handler"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template_data"

	//"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// 404 Handler
func (RM *RoutingMemory) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("not-found", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// GET '/gallery'
func (RM *RoutingMemory) GET_0_gallery(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	r.ParseForm()

	var imageNumber int = 1

	if r.FormValue("image") != "" {
		// specified.
		var err error

		imageNumber, err = strconv.Atoi(r.FormValue("image"))
		if err != nil {
			// presume imageNumber is 1
			imageNumber = 1
		}

		if imageNumber > RM.gallery_photos.ImageCount {
			// greater than image count, rubberband to last available image.
			imageNumber = RM.gallery_photos.ImageCount
		} else if imageNumber < 1 {
			imageNumber = 1
		}
	}

	// get image from gallery.
	image, err := RM.gallery_photos.Get(imageNumber)

	if err != nil {
		// TODO (fnuer): replace this with actual handler after testing
		http.Redirect(w, r, "/error?type=image_gallery_not_found", 303)
		return
	}

	var replacestr []string

	var next_disabled = ""
	var previous_disabled = ""

	if imageNumber == RM.gallery_photos.ImageCount {
		// disable next.
		next_disabled = "disabled=\"disabled\""
	} else if imageNumber == 1 {
		// disable previous.
		previous_disabled = "disabled=\"disabled\""
	}
	if image.HasCredit {
		// author credited.
		replacestr = []string{
			"%%IMAGE_NAME%%", image.Name,
			"%%IMAGE_URL%%", image.ImageURL,
			"%%IMAGE_DESCRIPTION%%", image.Description,
			"%%IMAGE_CREDIT%%", image.CreditedAuthor,
			"%%IMAGE_CURRENT%%", strconv.Itoa(imageNumber),
			"%%IMAGE_MAX%%", strconv.Itoa(RM.gallery_photos.ImageCount),

			"%%IMAGE_NEXT_URL%%", strconv.Itoa(imageNumber + 1),
			"%%IMAGE_PREVIOUS_URL%%", strconv.Itoa(imageNumber - 1),

			"%%NEXT_DISABLED%%", next_disabled,
			"%%PREVIOUS_DISABLED%%", previous_disabled,
		}
	} else {
		// author not credited.
		replacestr = []string{
			"%%IMAGE_NAME%%", image.Name,
			"%%IMAGE_URL%%", image.ImageURL,
			"%%IMAGE_DESCRIPTION%%", image.Description,
			"%%IMAGE_CREDIT%%", "Unknown",
			"%%IMAGE_CURRENT%%", strconv.Itoa(imageNumber),
			"%%IMAGE_MAX%%", strconv.Itoa(RM.gallery_photos.ImageCount),

			"%%IMAGE_NEXT_URL%%", strconv.Itoa(imageNumber + 1),
			"%%IMAGE_PREVIOUS_URL%%", strconv.Itoa(imageNumber - 1),

			"%%NEXT_DISABLED%%", next_disabled,
			"%%PREVIOUS_DISABLED%%", previous_disabled,
		}
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("photo-gallery", sess, RM.templateCollection)
	baseTemp = strings.NewReplacer(replacestr...).Replace(baseTemp)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// GET '/cookie_consent'
func (RM *RoutingMemory) GET_0_cookie_consent(w http.ResponseWriter, r *http.Request) {
	baseTemp := gdpr_static_PageRender("cookie_consent", RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/maintenance'
func (RM *RoutingMemory) GET_0_maintenance(w http.ResponseWriter, r *http.Request) {
	baseTemp := gdpr_nobanner_static_PageRender("maintenance", RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/cookie_consent/withdraw'
func (RM *RoutingMemory) GET_0_cookie_consent_0_withdraw(w http.ResponseWriter, r *http.Request) {
	// get old consent to remove from db
	cookie, err := r.Cookie("AMC_CONSENT")

	if err != nil {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	err = RM.database.Deleteconsent(cookie.Value)

	if err != nil {
		http.Redirect(w, r, "/error?type=consent_withdraw_db_error", 303)
		return
	}

	// remove all cookies.

	c_consent := &http.Cookie{
		Name:    "AMC_CONSENT",
		Path:    "/",
		Expires: time.Unix(0, 0),
		Value:   "CONSENT WITHDRAWN",
	}

	c_session := &http.Cookie{
		Name:    "AMC_SESSION",
		Path:    "/",
		Expires: time.Unix(0, 0),
		Value:   "CONSENT WITHDRAWN",
	}

	c_accessibility := &http.Cookie{
		Name:    "accessibility",
		Path:    "/",
		Expires: time.Unix(0, 0),
		Value:   "CONSENT WITHDRAWN",
	}

	http.SetCookie(w, c_consent)
	http.SetCookie(w, c_session)
	http.SetCookie(w, c_accessibility)

	baseTemp := gdpr_static_PageRender("cookie_consent_withdrawn", RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// POST '/cookie_consent/submit'
func (RM *RoutingMemory) GET_0_cookie_consent_0_submit(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	uuid := RM.consent_handler.NewConsentObject(true, r.FormValue("preference") == "preference")

	c := &http.Cookie{
		Name:    "AMC_CONSENT",
		Path:    "/",
		Value:   uuid,
		Expires: time.Now().Add(17520 * time.Hour),
	}

	http.SetCookie(w, c)

	http.Redirect(w, r, "/", 303)
}

// GET '/'
func (RM *RoutingMemory) GET_0_(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("index", sess, RM.templateCollection)

	// get from db.
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	random := r1.Intn(RM.gallery_photos.ImageCount)

	imageURL := "/static/img/screenshot-preview.png"

	if random < 1 {
		random = 1
	}

	img, err := RM.gallery_photos.Get(random)

	if err == nil {
		imageURL = img.ImageURL
	}

	baseTemp = strings.Replace(baseTemp, "%%RANDOM_URL%%", imageURL, 1)
	baseTemp = strings.Replace(baseTemp, "%%RANDOM_GOTO%%", "/gallery?image="+strconv.Itoa(random), 1)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/'
func (RM *RoutingMemory) GET_0_skin_server(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("skin-server", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/leaving'
func (RM *RoutingMemory) GET_0_leaving(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("leaving", sess, RM.templateCollection)

	baseTemp = strings.Replace(baseTemp, "%%URL%%", r.FormValue("url"), -1)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/login_required'
func (RM *RoutingMemory) GET_0_login_required(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("login-required", sess, RM.templateCollection)

	baseTemp = strings.Replace(baseTemp, "%%URL%%", r.FormValue("url"), -1)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/web_button'
func (RM *RoutingMemory) GET_0_web_button(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("web_button", sess, RM.templateCollection)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/cookies'
func (RM *RoutingMemory) GET_0_cookies(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("cookies", sess, RM.templateCollection)

	var consentState string = "Allow necessary cookies only"

	if consent.AllowPreferences {
		consentState = "Allow all (Necessary, Preferences)"
	}

	baseTemp = strings.Replace(baseTemp, "%%CONSENT_STATE%%", consentState, 1)

	baseTemp = strings.Replace(baseTemp, "%%CONSENT_ID%%", consent.ConsentID, 1)
	baseTemp = strings.Replace(baseTemp, "%%CONSENT_DATE%%", time.Unix(consent.ConsentTimestamp, 0).Format(time.UnixDate), 1)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// GET '/privacy'
func (RM *RoutingMemory) GET_0_privacy(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("privacy", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// GET '/credits'
func (RM *RoutingMemory) GET_0_credits(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("credits", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// GET '/licenses'
func (RM *RoutingMemory) GET_0_licenses(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()
	baseTemp := static_PageRender("licenses", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// GET '/error'
func (RM *RoutingMemory) GET_0_error(w http.ResponseWriter, r *http.Request) {
	// need to parse form beforehand
	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	rd_check := sess.Data.CurrentRedirect

	if rd_check == "" {
		// no redirect available so go to root.
		rd_check = "/"
	}

	var error_text string

	if r.FormValue("type") == "" {
		// no error passed.
		error_text = "You found the secret error. How you found this, I have no idea.<br /><br />You should go home, though..."
	} else {
		// scan hashmap for error types.
		hash_error, english_text := template_data.Get_ErrorString(r.FormValue("type"))

		if hash_error != nil {
			error_text = "You found yet another secret error! This is because an invalid error type was passed to the error handler!"
		} else {
			error_text = english_text
		}
	}

	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("error", sess)
	baseTemp = strings.Replace(baseTemp, "%%ENCOUNTERED_ERROR_INFO%%", error_text, 1)
	baseTemp = strings.Replace(baseTemp, "%%REDIRECTOR_URL%%", rd_check, 1)
	baseTemp = strings.Replace(baseTemp, "%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0), -1)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// GET '/preprod_info'
func (RM *RoutingMemory) GET_0_preprod_info(w http.ResponseWriter, r *http.Request) {
	if !RM.site_configuration.IsPreproduction {
		// not preprod.
		http.Redirect(w, r, "/", 302)
		return
	}

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	baseTemp := static_PageRender("preprod_info", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// GET '/beta_info'
func (RM *RoutingMemory) GET_0_beta_info(w http.ResponseWriter, r *http.Request) {
	if !RM.site_configuration.IsBeta {
		// not preprod.
		http.Redirect(w, r, "/", 302)
		return
	}

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	baseTemp := static_PageRender("beta_info", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// OAUTH2 FLOW STUFF //

// GET '/login/callback'
func (RM *RoutingMemory) GET_0_login_0_callback(w http.ResponseWriter, r *http.Request) {
	// need to parse form beforehand
	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	if r.FormValue("state") != RM.state {
		// user passed no or invalid state
		//return c.Status(400).SendString("OAuth2 state does not match!")
		http.Redirect(w, r, "/error?type=oauth2_error", 303)
	}

	if r.FormValue("error") != "" {
		// something on discord's end cocked up.
		// user more than likely rejected OAuth2 flow
		http.Redirect(w, r, "/error?type=oauth2_denied", 303)
		return
	}

	// exchange oauth code for an access token
	err, user := RM.oauth2_struct.GetInformation(r.FormValue("code"))

	if err != nil {
		// something major happened - maybe an expired code?
		//log.Println(err)

		http.Redirect(w, r, "/error?type=oauth2_error", 303)
		return
		//return c.Status(500).SendString("Error occurred during OAuth2 flow. Please contact an administrator.")
	}

	user_banned := RM.database.Checkdeleted(user.DiscordId)

	if user_banned {
		// user is banned!
		http.Redirect(w, r, "/error?type=banned", 303)
		return
	}

	// get session, initialise login session

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	session_init := auth.InitialiseSession(sess, RM.database, user)

	if !session_init {
		// user needs to register.
		http.Redirect(w, r, "/register", 303)
		return
	} else {
		// user is logged in.

		// note to alex... before he screams at me again.
		// THIS FUNCTION RETURNS NO FUCKING ERRORS EVER.
		// IT HAS NO ERROR RETURN VALUE AT ALL!!!!!
		// SO STOP SCREAMING AT ME ;-;
		isadmin, _ := RM.database.Checkadmin(user.DiscordId)

		if isadmin && RM.dynamic_config.MaintenanceMode {
			http.Redirect(w, r, "/admin", 303)
		} else {
			http.Redirect(w, r, "/profile", 303)
		}

		return
	}

	// if they manage to get past that if statement, something went horribly wrong.
	// throw them back to the home page.
	http.Redirect(w, r, "/", 303)
	return
}

// GET '/logout'
func (RM *RoutingMemory) GET_0_logout(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	auth.EndSession(sess)
	http.Redirect(w, r, "/", 303)
	return
}

// DYNAMIC ROUTES //

// STATUS //

// GET '/status'
func (RM *RoutingMemory) GET_0_status(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("status", sess)

	start_addProc := time.Now().UnixNano()

	allServices := RM.database.Getallservices()
	defer allServices.Close()

	serviceTable := ""

	for allServices.Next() {
		var serviceName string
		var serviceIP string
		var serviceStatus int

		var statName string = "Unknown"

		err := allServices.Scan(&serviceIP, &serviceName, &serviceStatus)
		if err != nil {
			log.Fatal(err)
		}

		err, _, statName = template_data.Get_ServiceStrings(serviceStatus)

		if err != nil {
			RM.logger.Err.Println("An unknown service state was found.")
		}

		serviceTable = serviceTable + `<tr><td>` + serviceName + `</td><td>` + serviceIP + `</td><td>` + statName + `</td></tr>`

	}

	baseTemp = strings.Replace(baseTemp, "%%STATUS_TABLE%%", serviceTable, -1)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	end_addProc := time.Now().UnixNano()

	baseTemp = strings.Replace(baseTemp, "%%TIME%%", template.GetTemplateTimestamp(renderTime, start_addProc, end_addProc), -1)

	w.Write([]byte(baseTemp))
	return
}

// PROFILE //

// REGISTRATION //

// GET '/register'
func (RM *RoutingMemory) GET_0_register(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	login_check, account_check := auth.CheckAuthentication(sess)

	if !login_check {
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	if account_check {
		http.Redirect(w, r, "/profile", 303)
		return
	}

	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("register", sess)
	baseTemp = strings.Replace(baseTemp, "%%D_ID%%", ""+sess.Data.Id, -1)

	baseTemp = strings.Replace(baseTemp, "%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0), -1)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

	w.Write([]byte(baseTemp))
	return
}

// GET '/register/submit'
func (RM *RoutingMemory) GET_0_register_0_submit(w http.ResponseWriter, r *http.Request) {
	// need to parse form beforehand
	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	login_check, account_check := auth.CheckAuthentication(sess)

	if !login_check {
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	if account_check {
		http.Redirect(w, r, "/profile", 303)
		return
	}

	if r.FormValue("username") == "" {
		http.Redirect(w, r, "/register", 303)
		return
	}

	if RM.dynamic_config.AccountCreationDisabled {
		http.Redirect(w, r, "/error?type=reg_disabled", 303)
		return
	}

	// create their user.
	err := RM.database.Createuser(r.FormValue("username"), sess.Data.Id)

	if err == nil {
		// success.
		sess.Data.MinecraftUsername = r.FormValue("username")
		sess.Data.StateOnLogin = "1"
		sess.Data.HasAccount = true

		//sess.Save()

		http.Redirect(w, r, "/profile", 303)
		return
	} else {
		// not success.
		sess.Data.CurrentRedirect = "/register"
		http.Redirect(w, r, "/error?type="+err.Error(), 303)
		return
	}
}

// USERNAME CHANGE AND SUBMIT //
