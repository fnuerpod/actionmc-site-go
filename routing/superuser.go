package routing

import (
	//"encoding/hex"
	"git.dsrt-int.net/actionmc/actionmc-site-go/auth"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template"
	"git.dsrt-int.net/actionmc/actionmc-site-go/unsafety"

	"git.dsrt-int.net/actionmc/actionmc-site-go/template_data"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/cpu"

	//"math/rand"
	//"strconv"

	"net/http"
	"strconv"
	"strings"
	"time"
	//"time"
)

/*
	route list for implementation:

	/admin/superuser

	/admin/superuser/template-reload/page
	/admin/superuser/template-reload/base
	/admin/superuser/template-reload/all

	/admin/superuser/sessions/purge-stale-gt5h
	/admin/superuser/sessions/purge-stale-gt1h
	/admin/superuser/sessions/purge-all

*/

// GET '/admin/superuser'
func (RM *RoutingMemory) GET_0_admin_0_superuser(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		return
	}

	// check if user is admin
	isadmin, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if !isadmin {
		// user is not admin.
		http.Redirect(w, r, "/profile", 302)
		return
	}

	if privlevel < 5 {
		http.Redirect(w, r, "/error?type=not_authorised", 302)
		return
	}

	cpuram := ""

	percent, _ := cpu.Percent(time.Second, false)

	if percent[0] > 50 && (((memory.FreeMemory() / memory.TotalMemory()) * 100) > 50) {
		cpuram = "<p class=\"error-text\">WARNING: RAM and CPU usage is above 50&#37;!</p><p class=\"error-text\">Consider killing some useless processes and manually purging sessions that are >5 hours old.</p><br />"
	} else if ((memory.FreeMemory() / memory.TotalMemory()) * 100) > 50 {
		cpuram = "<p class=\"error-text\">WARNING: RAM usage is above 50&#37;!</p><p class=\"error-text\">Consider manually purging sessions that are >5 hours old.</p><br />"
	} else if percent[0] > 50 {
		cpuram = "<p class=\"error-text\">WARNING: CPU usage is above 50&#37;!</p><p class=\"error-text\">Consider killing some useless processes, or upgrade server host machine.</p><br />"
	}

	replacestr := []string{
		"%%MEM_CUR%%", strconv.FormatUint(memory.FreeMemory()/100000000, 10) + "GB",
		"%%MEM_SYS%%", strconv.FormatUint(memory.TotalMemory()/1000000000, 10) + "GB",
		"%%CPU_PER%%", strconv.FormatInt(int64(percent[0]), 10),
		"%%CPU_RAM_ALERT%%", cpuram,
	}

	// user is an admin if they're beyond this point. create base template and send to user.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("admin-superuser", sess)

	sess.Data.CurrentRedirect = "/admin/superuser"

	//sess.Save()

	baseTemp = strings.NewReplacer(replacestr...).Replace(baseTemp)

	baseTemp = strings.Replace(baseTemp, "%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0), -1)

	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(baseTemp))

	return
}

// GET /admin/superuser/template_reload
func (RM *RoutingMemory) GET_0_admin_0_superuser_0_template_reload(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		return
	}

	// check if user is admin
	isadmin, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if !isadmin {
		// user is not admin.
		http.Redirect(w, r, "/profile", 302)
		return
	}

	if privlevel < 5 {
		http.Redirect(w, r, "/error?type=not_authorised", 302)
		return
	}

	// reload templates
	RM.templateCollection = template.NewTemplateCollection("data/templates", "data/bases", RM.logger)

	// redirect to the superuser page.
	http.Redirect(w, r, "/admin/superuser", 303)

}

// GET '/admin/superuser/dyn_cfg`
func (RM *RoutingMemory) GET_0_admin_0_superuser_0_dyn_cfg(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		return
	}

	// check if user is admin
	isadmin, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if !isadmin {
		// user is not admin.
		http.Redirect(w, r, "/profile", 302)
		return
	}

	if privlevel < 5 {
		http.Redirect(w, r, "/error?type=not_authorised", 302)
		return
	}

	// user is an admin if they're beyond this point. create base template and send to user.
	baseTemp, renderTime := RM.templateCollection.BuildPopulatedTemplate("admin-superuser-dyncfg", sess)

	sess.Data.CurrentRedirect = "/admin/superuser/dyn_cfg"

	//sess.Save()

	states := template_data.Dynconf_state_map

	joiner_info := []string{
		"%%MAINTENANCE_MODE%%", states[unsafety.BtoU8(RM.dynamic_config.MaintenanceMode)],
		"%%USER_REGISTRATION%%", states[1-unsafety.BtoU8(RM.dynamic_config.AccountCreationDisabled)],
		"%%ANNOUNCEMENT_BANNER%%", states[unsafety.BtoU8(RM.dynamic_config.ShowAnnounceBanner)],
		"%%ABAN_TEXT%%", RM.dynamic_config.AnnounceBannerText,
		"%%AMGMT_VS%%", RM.dynamic_config.ActionMGMT_LatestVersion,
	}

	baseTemp = strings.NewReplacer(joiner_info...).Replace(baseTemp)

	baseTemp = strings.Replace(baseTemp, "%%TIME%%", template.GetTemplateTimestamp(renderTime, 0, 0), -1)

	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte(baseTemp))

	return
}

// POST '/admin/superuser/dyn_cfg/submit'
func (RM *RoutingMemory) GET_0_admin_0_superuser_0_dyn_cfg_0_submit(w http.ResponseWriter, r *http.Request) {

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

	// check authentication
	auth_check := check_authentication(w, r, sess)

	if !auth_check {
		return
	}

	// check if user is admin
	isadmin, privlevel := RM.database.Checkadmin(sess.Data.Id)

	if !isadmin {
		// user is not admin.
		http.Redirect(w, r, "/profile", 302)
		return
	}

	if privlevel < 5 {
		http.Redirect(w, r, "/error?type=not_authorised", 302)
		return
	}

	if r.FormValue("user_reg") == "" {
		http.Redirect(w, r, "/admin/superuser/dyn_cfg", 303)
		return
	} else {
		RM.dynamic_config.AccountCreationDisabled = (r.FormValue("user_reg") != "1")
	}

	if r.FormValue("maintenance_mode") == "" {
		http.Redirect(w, r, "/admin/superuser/dyn_cfg", 303)
		return
	} else {
		RM.dynamic_config.MaintenanceMode = (r.FormValue("maintenance_mode") == "1")
	}

	if r.FormValue("announcement_banner_enable") == "" {
		http.Redirect(w, r, "/admin/superuser/dyn_cfg", 303)
		return
	} else {
		RM.dynamic_config.ShowAnnounceBanner = (r.FormValue("announcement_banner_enable") == "1")

	}

	RM.dynamic_config.ActionMGMT_LatestVersion = r.FormValue("actionmgmt_version")

	RM.dynamic_config.AnnounceBannerText = r.FormValue("announcement_text")

	// save now.
	RM.dynamic_config.Save()

	// reload templates
	RM.templateCollection = template.NewTemplateCollection("data/templates", "data/bases", RM.logger)

	http.Redirect(w, r, "/admin/superuser", 303)

}
