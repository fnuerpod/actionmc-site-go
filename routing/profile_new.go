package routing

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"net/http"

	"git.dsrt-int.net/actionmc/actionmc-site-go/auth"
	"git.dsrt-int.net/actionmc/actionmc-site-go/sessions"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template_data"

	"math/rand"
	"strconv"
	"strings"
	"time"

	"image"
	_ "image/png"
	"io"

	"log"
	"os"
)

func (RM *RoutingMemory) GET_0_profile(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	// get base template.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("profile", sess)

	additional_processing_start := time.Now().UnixNano()

	// get appropriate state HTML data from template data.
	db_state, ok := RM.database.Getstate(sess.Data.Id)

	var rawState string

	if ok {
		// database has a hit for this user's state.
		rawState = db_state
	} else {
		// probably a new user or something so we'll give them the default login
		// state for their session.
		rawState = sess.Data.StateOnLogin
	}

	// convert to int in order to pass to template data.
	intRaw, _ := strconv.Atoi(rawState)

	// get state name and HTML strings.
	_, skinHTML, capeHTML, state := template_data.Get_ProfileStrings(intRaw)

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	// to reduce the amount of necessary branches, we presume level 0 is no privilege.
	err, admin_strings := template_data.Get_AdminStrings(privlevel)

	if err != nil {
		// invalid privilege level...
		RM.logger.Debug.Println("Administrator/user with unknown privilege level on profile page.")
	}

	// we end it here even though we're yet to process the string replacer.
	// this is because another string replaced would have to be called later on,
	// this reducing efficiency.
	end_additional_processing := time.Now().UnixNano()

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%MC_USERNAME%%", sess.Data.MinecraftUsername,
		"%%D_ID%%", sess.Data.Id,
		"%%SKIN_CHANGE%%", skinHTML,
		"%%CAPE_CHANGE%%", capeHTML,
		"%%STATE%%", state,
		"%%ADMIN_ACCESS%%", admin_strings,
		"%%TIME%%", template.GetTemplateTimestamp(renderTime, additional_processing_start, end_additional_processing),
	}

	baseTemp = strings.NewReplacer(render_values...).Replace(baseTemp)

	// set content type and send final template to user.
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return

}

func (RM *RoutingMemory) GET_0_profile_0_edit_0_username(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	// get base template.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("change_name", sess)

	// as this really calls some variables and does string replacement it isn't
	// necessary for us to do additional processing timestamps.

	sess.Data.CSRFPreventionString = sessions.GenerateCSRFString()

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%CURRENT_USERNAME%%", sess.Data.MinecraftUsername,
		"%%D_ID%%", sess.Data.Id,
		"%%CSRF_STRING%%", sess.Data.CSRFPreventionString,
		"%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0),
	}

	baseTemp = strings.NewReplacer(render_values...).Replace(baseTemp)

	// set content type and send final template to user.
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return

}

func (RM *RoutingMemory) GET_0_profile_0_sut(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	// get base template.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("sut", sess)

	// as this really calls some variables and does string replacement it isn't
	// necessary for us to do additional processing timestamps.

	key, ok := RM.database.GetTOTPSecret(sess.Data.Id)

	var totpKey *otp.Key

	var totpToken string

	if !ok {
		// user hasn't got a totp secret yet.
		// generate one
		totpKey, _ = totp.Generate(totp.GenerateOpts{
			Issuer:      "ActionMC SULT",
			AccountName: sess.Data.Id,
		})

		err := RM.database.CreateTOTP(sess.Data.Id, totpKey.Secret())

		if err != nil {
			http.Redirect(w, r, "/error?type="+err.Error(), 303)
			return
		}

		totpToken, err = totp.GenerateCode(key, time.Now())

		if err != nil {
			http.Redirect(w, r, "/error?type=otp_gen_error", 303)
			return
		}

	} else {
		totpTokena, err := totp.GenerateCode(key, time.Now())

		if err != nil {
			RM.logger.Err.Println(err)
			http.Redirect(w, r, "/error?type=otp_gen_error", 303)
			return
		}

		totpToken = totpTokena
	}

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%TOKEN%%", totpToken,
		"%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0),
	}

	baseTemp = strings.NewReplacer(render_values...).Replace(baseTemp)

	// set content type and send final template to user.
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return

}

func (RM *RoutingMemory) POST_0_profile_0_edit_0_username_0_submit(w http.ResponseWriter, r *http.Request) {
	// parse form before doing anything
	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	if r.FormValue("csrf") != sess.Data.CSRFPreventionString {
		// CSRF is taking place. redirect user to error page.
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=csrf_alert", 303)
		return
	}

	// check if new username has been passed by client.
	if r.FormValue("username") == "" {
		http.Redirect(w, r, "/profile", 303)
		return
	}

	// get user's previous name change timestamp.
	previous_change_stamp := RM.database.Getuser_changetime(sess.Data.Id)

	// check if user has changed name within last 30 days - if so, do not allow
	// them to change their name.
	if (int(time.Now().Unix()) - previous_change_stamp) < 2592000 {
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=user_change_ratelimit", 303)
		return
	}

	// update their user.
	err := RM.database.Updateuser(r.FormValue("username"), sess.Data.Id)

	if err == nil {
		// success.
		err := RM.database.Updateuser_skintime(sess.Data.Id, int(time.Now().Unix()))

		if err == nil {
			sess.Data.MinecraftUsername = r.FormValue("username")
			//sess.Save()

			http.Redirect(w, r, "/profile", 303)
			return
		} else {
			sess.Data.CurrentRedirect = "/profile/edit/username"
			http.Redirect(w, r, "/error?type="+err.Error(), 303)
			return
		}

	} else {
		// not success.
		sess.Data.CurrentRedirect = "/profile/edit/username"
		http.Redirect(w, r, "/error?type="+err.Error(), 303)
		return
	}

}

func (RM *RoutingMemory) GET_0_profile_0_delete(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	// get base template.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("delete_confirm", sess)

	additional_processing_start := time.Now().UnixNano()

	// get appropriate state HTML data from template data.
	db_state, ok := RM.database.Getstate(sess.Data.Id)

	var rawState string

	if ok {
		// database has a hit for this user's state.
		rawState = db_state
	} else {
		// probably a new user or something so we'll give them the default login
		// state for their session.
		rawState = sess.Data.StateOnLogin
	}

	// convert to int in order to pass to template data.
	intRaw, _ := strconv.Atoi(rawState)

	// get state name and HTML strings.
	_, _, _, state := template_data.Get_ProfileStrings(intRaw)

	// generate a random four digit number to act as a deletion confirmation
	// string, which is a form of CSRF prevention.
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	rnd := r1.Intn(9999)
	random_key := strconv.FormatInt(int64(rnd), 10)

	// store this key in session so we can refer to it if the user confirms
	// to delete their account.
	sess.Data.DeletionString = random_key

	// we end it here even though we're yet to process the string replacer.
	// this is because another string replaced would have to be called later on,
	// this reducing efficiency.
	end_additional_processing := time.Now().UnixNano()

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%MC_USERNAME%%", sess.Data.MinecraftUsername,
		"%%D_ID%%", sess.Data.Id,
		"%%STATE%%", state,
		"%%DELETION_STRING%%", random_key,
		"%%TIME%%", template.GetTemplateTimestamp(renderTime, additional_processing_start, end_additional_processing),
	}

	baseTemp = strings.NewReplacer(render_values...).Replace(baseTemp)

	// set content type and send final template to user.
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return

}

func (RM *RoutingMemory) POST_0_profile_0_delete_0_confirm(w http.ResponseWriter, r *http.Request) {
	// parse form before doing anything
	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	// check if confirmation string has been passed by client.
	if r.FormValue("confirmation") == "" {
		// redirect to error.
		http.Redirect(w, r, "/error?type=del_no_confirm", 303)
		return
	}

	// check if confirmation string is the same.
	// basically a csrf check.
	if sess.Data.DeletionString != r.FormValue("confirmation") {
		sess.Data.CurrentRedirect = "/profile/delete"
		http.Redirect(w, r, "/error?type=del_invalid_confirm", 303)
		return
	}

	// if user has got this far, they have consented to a full account
	// deletion.

	// check if user is banned.

	// get appropriate state from database.
	db_state, ok := RM.database.Getstate(sess.Data.Id)

	var rawState string

	if ok {
		// database has a hit for this user's state.
		rawState = db_state
	} else {
		// probably a new user or something so we'll give them the default login
		// state for their session.
		rawState = sess.Data.StateOnLogin
	}

	if rawState == "4" || rawState == "3" {
		// user is banned and must be added to banned list.
		_, err := RM.database.Adddeleted_banned(sess.Data.Id)

		if err != "" {
			sess.Data.CurrentRedirect = "/profile/delete"
			http.Redirect(w, r, "/error?type="+err, 303)
			return
		}
	}

	// delete user. NOT UNDOABLE.
	err := RM.database.Deleteuser(sess.Data.Id)

	if err != nil {
		sess.Data.CurrentRedirect = "/profile/delete"
		http.Redirect(w, r, "/error?type="+err.Error(), 302)
		return
	}

	// destroy their session too.
	sess.Destroy()

	// redirect user to completion page.
	http.Redirect(w, r, "/profile/delete/done", 302)
	return
}

// GET '/profile/delete/done'
func (RM *RoutingMemory) GET_0_profile_0_delete_0_done(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	baseTemp := static_PageRender("profile_deleted", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// POST '/profile/skin_submit'
func (RM *RoutingMemory) POST_0_profile_0_skin_submit(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	r.ParseMultipartForm(8192)

	file, handler, err := r.FormFile("upload")
	if err != nil {
		//fmt.Println("Error Retrieving the File")
		//fmt.Println(err)
		RM.logger.Err.Println("Error retreiving skin file.")
		RM.logger.Err.Println(err)

		sess.Data.CurrentRedirect = "/profile"

		http.Redirect(w, r, "/error?type=error_skin_upload", 303)
		return
	}

	defer file.Close()

	if handler.Size > 8192 {
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=upload_too_large", 303)
		return
	}

	if handler.Header["Content-Type"][0] != "image/png" {
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=upload_bad_file", 303)
		return
	}

	im, _, err := image.DecodeConfig(file)

	if err != nil {
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=upload_bad_file", 303)
		return
	}

	// this is great since we won't write to a file till it's verified.

	if im.Width != 64 {
		// width is not 64, which doesn't work under any circumstance
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=skin_bad_resolution", 303)
		return
	} else if (im.Height != 64) && (im.Height != 32) {
		// height doesn't match Steve or Alex model.
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=skin_bad_resolution", 303)
		return
	}

	var fileBytes []byte

	_, err = file.Seek(0, io.SeekStart)

	if err != nil {
		goto errorState
	}

	fileBytes, err = io.ReadAll(file)

	if err != nil {
		goto errorState
	}

	err = os.WriteFile("./ugc_store/skin/"+sess.Data.Id+".png", fileBytes, 0666)

errorState:
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/profile/skin_upload_ok", 303)
	return
}

// POST '/profile/cape_submit'
func (RM *RoutingMemory) POST_0_profile_0_cape_submit(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	r.ParseMultipartForm(8192)

	file, handler, err := r.FormFile("upload")
	if err != nil {
		RM.logger.Err.Println("Error retreiving cape file.")
		RM.logger.Err.Println(err)

		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=error_cape_upload", 303)
		return
	}

	defer file.Close()

	if handler.Size > 8192 {
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=upload_too_large", 303)
		return
	}

	if handler.Header["Content-Type"][0] != "image/png" {
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=upload_bad_file", 303)
		return
	}

	im, _, err := image.DecodeConfig(file)

	if err != nil {
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=upload_bad_file", 303)
		return
	}

	// this is great since we won't write to a file till it's verified.
	if (im.Width != 64) || (im.Height != 32) {
		// width is not 64, which doesn't work under any circumstance
		sess.Data.CurrentRedirect = "/profile"
		http.Redirect(w, r, "/error?type=cape_bad_resolution", 303)
		return
	}

	var fileBytes []byte

	_, err = file.Seek(0, io.SeekStart)

	if err != nil {
		goto errorState
	}

	fileBytes, err = io.ReadAll(file)

	if err != nil {
		goto errorState
	}

	// TODO(ultrabear) make this fucking shit not use a single call to writefile i swear to fuck

	err = os.WriteFile("./ugc_store/cape/"+sess.Data.Id+".png", fileBytes, 0666)

errorState:
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/profile/cape_upload_ok", 303)
	return
}

func (RM *RoutingMemory) GET_0_profile_0_skin_upload_ok(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	baseTemp := static_PageRender("skin-upload-ok", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

func (RM *RoutingMemory) GET_0_profile_0_cape_upload_ok(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	// get their session.
	sess := auth.GetSession(RM.store, w, r, consent)

	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		// return them to the login page
		http.Redirect(w, r, "/login_required", 303)
		return
	}

	baseTemp := static_PageRender("cape-upload-ok", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}
