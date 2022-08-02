package routing

import (
	"net/http"
	"strings"

	"git.dsrt-int.net/actionmc/actionmc-site-go/auth"
	"git.dsrt-int.net/actionmc/actionmc-site-go/sessions"
	"git.dsrt-int.net/actionmc/actionmc-site-go/template_data"
	"git.dsrt-int.net/actionmc/actionmc-site-go/unsafety"
)

// GET '/accessibility_options'
func (RM *RoutingMemory) GET_0_accessibility_options(w http.ResponseWriter, r *http.Request) {
	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	if !consent.AllowPreferences {
		// no consent to preference cookies.
		http.Redirect(w, r, "/error?type=cookie_consent_no_preference", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	baseTemp := static_PageRender("accessibility_options", sess, RM.templateCollection)

	// check enablement and disablement status

	accessibility := sess.Data.AccessibilityData
	states := template_data.Accessibility_state_map

	joiner_info := []string{
		"%%DYSLEXIA%%", states[unsafety.BtoU8(accessibility.EnableDyslexia)],
		"%%CONTRAST%%", states[unsafety.BtoU8(accessibility.EnableHiContrast)],
		"%%NOIMAGE%%", states[unsafety.BtoU8(accessibility.EnableNoImage)],
	}

	baseTemp = strings.NewReplacer(joiner_info...).Replace(baseTemp)

	w.Header().Add("Content-Type", "text/html")
	w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte(baseTemp))
	return
}

// GET '/accessibility_options/change'
func (RM *RoutingMemory) GET_0_accessibility_options_0_change(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	// do GDPR check first
	consent := RM.consent_handler.GetConsentObject(r)

	if !consent.Exists {
		// NO GDPR CONSENT. DO INTERSTITAL.
		http.Redirect(w, r, "/cookie_consent", 303)
		return
	}

	if !consent.AllowPreferences {
		// no consent to preference cookies.
		http.Redirect(w, r, "/error?type=cookie_consent_no_preference", 303)
		return
	}

	sess := auth.GetSession(RM.store, w, r, consent)
	defer sess.Close()

	// now we're in the session, we just want to reinitialise the bitset cookie
	// TODO: move this to its own handler so we dont have copy paste mumbo jumbo.

	acc_obj := &sessions.AccessibilityData{
		EnableDyslexia:   (r.FormValue("dyslexia_mode") != "0"),
		EnableHiContrast: (r.FormValue("high_contrast") != "0"),
		EnableNoImage:    (r.FormValue("no_image") != "0"),
	}

	// set accessibility data object on server side properly.
	sess.Data.AccessibilityData = *acc_obj

	// Encode to hex
	var bset sessions.AccessibilityBitset

	bset.Set(sessions.EnableDyslexia, acc_obj.EnableDyslexia)
	bset.Set(sessions.EnableHiContrast, acc_obj.EnableHiContrast)
	bset.Set(sessions.EnableNoImage, acc_obj.EnableNoImage)

	hexval := sessions.ByteToHex(byte(bset))

	sessions.SetAccessibilityCookie(w, string(hexval[:]))

	// now with new cookie set, redirect to homepage
	http.Redirect(w, r, "/", 303)
	return
}
