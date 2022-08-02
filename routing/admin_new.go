package routing

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.dsrt-int.net/actionmc/actionmc-site-go/auth"
	"git.dsrt-int.net/actionmc/actionmc-site-go/sessions"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template_data"
)

// GET '/admin'
func (RM *RoutingMemory) GET_0_admin(w http.ResponseWriter, r *http.Request) {
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

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if privlevel < 1 {
		// privilege level less than 1, not an admin!
		// for security dont tell them that its an admin endpoint, otherwise
		// we're gonna have people hitting it.
		http.Redirect(w, r, "/profile", 303)
		return
	}

	// set current redirect for errors.
	sess.Data.CurrentRedirect = "/admin"

	// to reduce the amount of necessary branches, we presume level 0 is no privilege.
	err, admin_strings := template_data.Get_AdminPageStrings(privlevel)

	if err != nil {
		// strange, their privilege level is invalid.
		// for safety we want to presume they are not an administrator
		// and throw them out as a regular user.
		http.Redirect(w, r, "/profile", 303)
		return
	}

	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("admin-home", sess)

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%ADMIN%%", admin_strings,
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

// GET '/admin/user'

func (RM *RoutingMemory) GET_0_admin_0_user(w http.ResponseWriter, r *http.Request) {
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

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if privlevel < 2 {
		// privilege level less than 2, not an admin!
		// for security dont tell them that its an admin endpoint, otherwise
		// we're gonna have people hitting it.
		http.Redirect(w, r, "/profile", 303)
		return
	}

	// set current redirect for errors.
	sess.Data.CurrentRedirect = "/admin/user"

	// get base template.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("admin", sess)

	additional_processing_start := time.Now().UnixNano()

	// initialise admin table string.
	var admin_table string

	// initialise admin table string slice for construction.
	var table_strings []string

	// initialise variable for state string as we MUST have this accessible to
	// the whole function.
	var defined_state int

	// check if user wants a specific state.
	if r.FormValue("state") != "" {
		// convert string to int.
		// have to define err here otherwise golang does a weird thing and
		// creates defined_state on a local scope too.
		var err error

		defined_state, err = strconv.Atoi(r.FormValue("state"))

		// if an error is thrown, just pass the root admin dir.
		if err != nil {
			http.Redirect(w, r, "/admin/user", 302)
			return
		}

		// if the defined state isn't within our bounds, pass the root admin
		// directory.
		if (defined_state > 7) || (defined_state < 1) {
			http.Redirect(w, r, "/admin/user", 302)
			return
		}

		// set redirect appropriately.
		sess.Data.CurrentRedirect = "/admin/user?state=" + r.FormValue("state")
	}

	// get all users from the site database.
	allUsers := RM.database.Getall()
	defer allUsers.Close()

	// iterate through all users to generate their table entry.
	for allUsers.Next() {
		var id string
		var userName string
		var state int
		var userChange int

		err := allUsers.Scan(&id, &userName, &state, &userChange)

		if err != nil {
			// this is not good. throw a fatal error.
			RM.logger.Fatal.Fatalln(err)
		}

		// check if we want to add this user to the table.
		if (state == defined_state) || (defined_state == 0) {
			// we want to add them.
			// get the state string data.
			_, selstate, curstate := template_data.Get_StateStrings(state)

			// we don't check for an error since that's really just
			// superficial logging.

			// TODO (fnuer): make this a string replacer function for further
			// optimisation.
			table_strings = append(table_strings, `<tr>
		<td>`+userName+`</td>
	  <td>`+id+`</td>
		<td>`+curstate+`</td>
		<td>
			<form action="/admin/user/edit/change_state" method="POST">
				<input type="text" name="csrf" value="%%CSRF_STRING%%" style="display:none;"\>
				<textarea style="display: none;" name="id">`+id+`</textarea>
				<select name="state" onchange="this.form.submit()">
					`+selstate+`
				</select><br /><br />
				
			</form>
		</td>
	  </tr>`)
		}
	}

	// we now have a populated table string directory.
	// join all the strings together.
	admin_table = strings.Join(table_strings, "")

	// get the current defined state name for the header.
	_, _, head_curstate := template_data.Get_StateStrings(defined_state)

	// again, no error checking here since it's superficial branching.

	// generate CSRF security key.
	sess.Data.CSRFPreventionString = sessions.GenerateCSRFString()

	// replace CSRF stuff in admin table.
	admin_table = strings.Replace(admin_table, "%%CSRF_STRING%%", sess.Data.CSRFPreventionString, -1)

	// we end it here even though we're yet to process the string replacer.
	// this is because another string replaced would have to be called later on,
	// this reducing efficiency.
	end_additional_processing := time.Now().UnixNano()

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%ADMINTABLE%%", admin_table,
		"%%TOTAL%%", strconv.Itoa(len(table_strings)),
		"%%PAGE%%", head_curstate,
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

// GET '/admin/user/edit/change_state'
func (RM *RoutingMemory) GET_0_admin_0_user_0_edit_0_change_state(w http.ResponseWriter, r *http.Request) {
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

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if privlevel < 3 {
		// privilege level less than 3, not allowed to write!
		// throw error.
		http.Redirect(w, r, "/error?type=no_permission_whitelist", 303)
		return
	}

	// check if form values are present.
	if r.FormValue("state") == "" {
		// no state value. return no state error.
		http.Redirect(w, r, "/error?type=change_state_no_state", 303)
		return
	}

	if r.FormValue("id") == "" {
		// no id value. return no id error.
		http.Redirect(w, r, "/error?type=change_state_no_id", 303)
		return
	}

	if r.FormValue("csrf") == "" {
		// csrf exception. security error.
		http.Redirect(w, r, "/error?type=csrf_alert", 303)
		return
	} else {
		// csrf is here. check it.
		if r.FormValue("csrf") != sess.Data.CSRFPreventionString {
			// CSRF error! this is a bad action. alert user.
			http.Redirect(w, r, "/error?type=csrf_alert", 303)
			return
		}
	}

	if sess.Data.CurrentRedirect == "" {
		// no redirection. another security error.
		http.Redirect(w, r, "/error?type=no_redirect_security_risk", 303)
	}

	RM.database.Changestate(r.FormValue("id"), r.FormValue("state"))

	http.Redirect(w, r, sess.Data.CurrentRedirect, 303)
	return
}

// GET '/admin/user/search'
func (RM *RoutingMemory) GET_0_admin_0_user_0_search(w http.ResponseWriter, r *http.Request) {
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

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if privlevel < 3 {
		// privilege level less than 3, not allowed to write!
		// throw error.
		http.Redirect(w, r, "/error?type=no_permission_whitelist", 303)
		return
	}

	// set redirect for legacy reasons.
	sess.Data.CurrentRedirect = "/admin/user/search"

	baseTemp := static_PageRender("admin-search", sess, RM.templateCollection)
	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Write([]byte(baseTemp))
	return
}

// POST '/admin/user/search/id'
func (RM *RoutingMemory) POST_0_admin_0_user_0_search_0_id(w http.ResponseWriter, r *http.Request) {
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

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if privlevel < 2 {
		// privilege level less than 2, not an admin!
		// for security dont tell them that its an admin endpoint, otherwise
		// we're gonna have people hitting it.
		http.Redirect(w, r, "/profile", 303)
		return
	}

	// set current redirect for errors.
	// as much as I dont like using a GET and POST at the same time, here it is
	// the only way to perform the redirect currently.
	sess.Data.CurrentRedirect = "/admin/user/search/id?id=" + r.FormValue("id")

	// get base template.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("admin", sess)

	additional_processing_start := time.Now().UnixNano()

	// initialise admin table string.
	var admin_table string

	// initialise admin table string slice for construction.
	var table_strings []string

	// check if user wants a specific state.
	if r.FormValue("id") == "" {
		// no id passed.
		// redirect user back to root
		http.Redirect(w, r, "/admin/user", 303)
		return
	}

	// get all users from the site database.
	allUsers := RM.database.Getall()
	defer allUsers.Close()

	// iterate through all users to generate their table entry.
	for allUsers.Next() {
		var id string
		var userName string
		var state int
		var userChange int

		err := allUsers.Scan(&id, &userName, &state, &userChange)

		if err != nil {
			// this is not good. throw a fatal error.
			RM.logger.Fatal.Fatalln(err)
		}

		// check if we want to add this user to the table.
		if id == r.FormValue("id") {
			// we want to add them.
			// get the state string data.
			_, selstate, curstate := template_data.Get_StateStrings(state)

			// we don't check for an error since that's really just
			// superficial logging.

			// TODO (fnuer): make this a string replacer function for further
			// optimisation.
			table_strings = append(table_strings, `<tr>
		<td>`+userName+`</td>
	  <td>`+id+`</td>
		<td>`+curstate+`</td>
		<td>
			<form action="/admin/user/edit/change_state" method="POST">
				<input type="text" name="csrf" value="%%CSRF_STRING%%" style="display:none;"\>
				<textarea style="display: none;" name="id">`+id+`</textarea>
				<select name="state" onchange="this.form.submit()">
					`+selstate+`
				</select><br /><br />
				
			</form>
		</td>
	  </tr>`)
		}
	}

	// we now have a populated table string directory.
	// join all the strings together.
	admin_table = strings.Join(table_strings, "")

	// again, no error checking here since it's superficial branching.

	// generate CSRF security key.
	sess.Data.CSRFPreventionString = sessions.GenerateCSRFString()

	// replace CSRF stuff in admin table.
	admin_table = strings.Replace(admin_table, "%%CSRF_STRING%%", sess.Data.CSRFPreventionString, -1)

	// we end it here even though we're yet to process the string replacer.
	// this is because another string replaced would have to be called later on,
	// this reducing efficiency.
	end_additional_processing := time.Now().UnixNano()

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%ADMINTABLE%%", admin_table,
		"%%TOTAL%%", strconv.Itoa(len(table_strings)),
		"%%PAGE%%", r.FormValue("id"),
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

// POST '/admin/user/search/username'
func (RM *RoutingMemory) POST_0_admin_0_user_0_search_0_username(w http.ResponseWriter, r *http.Request) {
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

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if privlevel < 2 {
		// privilege level less than 2, not an admin!
		// for security dont tell them that its an admin endpoint, otherwise
		// we're gonna have people hitting it.
		http.Redirect(w, r, "/profile", 303)
		return
	}

	// set current redirect for errors.
	// as much as I dont like using a GET and POST at the same time, here it is
	// the only way to perform the redirect currently.
	sess.Data.CurrentRedirect = "/admin/user/search/username?username=" + r.FormValue("username")

	// get base template.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("admin", sess)

	additional_processing_start := time.Now().UnixNano()

	// initialise admin table string.
	var admin_table string

	// initialise admin table string slice for construction.
	var table_strings []string

	// check if user wants a specific state.
	if r.FormValue("username") == "" {
		// no username passed.
		// redirect user back to root
		http.Redirect(w, r, "/admin/user", 303)
		return
	}

	// get all users from the site database.
	allUsers := RM.database.Getall()
	defer allUsers.Close()

	// iterate through all users to generate their table entry.
	for allUsers.Next() {
		var id string
		var userName string
		var state int
		var userChange int

		err := allUsers.Scan(&id, &userName, &state, &userChange)

		if err != nil {
			// this is not good. throw a fatal error.
			RM.logger.Fatal.Fatalln(err)
		}

		// check if we want to add this user to the table.
		if userName == r.FormValue("username") {
			// we want to add them.
			// get the state string data.
			_, selstate, curstate := template_data.Get_StateStrings(state)

			// we don't check for an error since that's really just
			// superficial logging.

			// TODO (fnuer): make this a string replacer function for further
			// optimisation.
			table_strings = append(table_strings, `<tr>
		<td>`+userName+`</td>
	  <td>`+id+`</td>
		<td>`+curstate+`</td>
		<td>
			<form action="/admin/user/edit/change_state" method="POST">
				<input type="text" name="csrf" value="%%CSRF_STRING%%" style="display:none;"\>
				<textarea style="display: none;" name="id">`+id+`</textarea>
				<select name="state" onchange="this.form.submit()">
					`+selstate+`
				</select><br /><br />
				
			</form>
		</td>
	  </tr>`)
		}
	}

	// we now have a populated table string directory.
	// join all the strings together.
	admin_table = strings.Join(table_strings, "")

	// again, no error checking here since it's superficial branching.

	// generate CSRF security key.
	sess.Data.CSRFPreventionString = sessions.GenerateCSRFString()

	// replace CSRF stuff in admin table.
	admin_table = strings.Replace(admin_table, "%%CSRF_STRING%%", sess.Data.CSRFPreventionString, -1)

	// we end it here even though we're yet to process the string replacer.
	// this is because another string replaced would have to be called later on,
	// this reducing efficiency.
	end_additional_processing := time.Now().UnixNano()

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%ADMINTABLE%%", admin_table,
		"%%TOTAL%%", strconv.Itoa(len(table_strings)),
		"%%PAGE%%", r.FormValue("username"),
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

// GET '/admin/status'
func (RM *RoutingMemory) GET_0_admin_0_status(w http.ResponseWriter, r *http.Request) {
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

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if privlevel < 1 {
		// privilege level less than 1, not an admin!
		// for security dont tell them that its an admin endpoint, otherwise
		// we're gonna have people hitting it.
		http.Redirect(w, r, "/profile", 303)
		return
	}

	// set current redirect for errors.
	sess.Data.CurrentRedirect = "/admin/status"

	// get base template.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("admin-status", sess)

	additional_processing_start := time.Now().UnixNano()

	// initialise admin table string.
	var admin_table string

	// initialise admin table string slice for construction.
	var table_strings []string

	// get all users from the site database.
	allServices := RM.database.Getallservices()
	defer allServices.Close()

	// iterate through all users to generate their table entry.
	for allServices.Next() {
		var serviceIP string
		var serviceName string
		var serviceState int

		var selState string

		err := allServices.Scan(&serviceIP, &serviceName, &serviceState)

		if err != nil {
			// this is not good. throw a fatal error.
			RM.logger.Fatal.Fatalln(err)
		}

		// we don't check for an error since that's really just
		// superficial logging.

		// TODO (fnuer): make this a string replacer function for further
		// optimisation.

		_, selState, _ = template_data.Get_ServiceStrings(serviceState)

		table_strings = append(table_strings, `<tr>
  <td>`+serviceName+`</td>
  <td>`+serviceIP+`</td>
  <td>
	  <form action="/admin/status/edit/change_state" method="POST" auto>
		  <textarea name="ip" style="display:none;">`+serviceIP+`</textarea>
		  <select name="state" onchange="this.form.submit()">
			  `+selState+`
		  </select>
		<input type="text" name="csrf" value="%%CSRF_STRING%%" style="display:none;"\>
	  </form>
  </td>
  </tr>`)
	}

	// we now have a populated table string directory.
	// join all the strings together.
	admin_table = strings.Join(table_strings, "")

	// again, no error checking here since it's superficial branching.

	// generate CSRF security key.
	sess.Data.CSRFPreventionString = sessions.GenerateCSRFString()

	// replace CSRF stuff in admin table.
	admin_table = strings.Replace(admin_table, "%%CSRF_STRING%%", sess.Data.CSRFPreventionString, -1)

	// we end it here even though we're yet to process the string replacer.
	// this is because another string replaced would have to be called later on,
	// this reducing efficiency.
	end_additional_processing := time.Now().UnixNano()

	// initialise page specific render value things.
	var render_values []string = []string{
		"%%STATUS_TABLE%%", admin_table,
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

// POST '/admin/status/edit/change_state'
func (RM *RoutingMemory) POST_0_admin_0_status_0_edit_0_change_state(w http.ResponseWriter, r *http.Request) {
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

	// check if user is admin
	_, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if privlevel < 3 {
		// privilege level less than 3, not allowed to write!
		// throw error.
		http.Redirect(w, r, "/error?type=no_permission_whitelist", 303)
		return
	}

	// check if form values are present.
	if r.FormValue("state") == "" {
		// no state value. return no state error.
		http.Redirect(w, r, "/error?type=service_state_no_state", 303)
		return
	}

	if r.FormValue("ip") == "" {
		// no id value. return no id error.
		http.Redirect(w, r, "/error?type=service_state_no_ip", 303)
		return
	}

	if r.FormValue("csrf") == "" {
		// csrf exception. security error.
		http.Redirect(w, r, "/error?type=csrf_alert", 303)
		return
	} else {
		// csrf is here. check it.
		if r.FormValue("csrf") != sess.Data.CSRFPreventionString {
			// CSRF error! this is a bad action. alert user.
			http.Redirect(w, r, "/error?type=csrf_alert", 303)
			return
		}
	}

	if sess.Data.CurrentRedirect == "" {
		// no redirection. another security error.
		http.Redirect(w, r, "/error?type=no_redirect_security_risk", 303)
	}

	RM.database.Changeservicestate(r.FormValue("ip"), r.FormValue("state"))

	http.Redirect(w, r, sess.Data.CurrentRedirect, 303)
	return
}
