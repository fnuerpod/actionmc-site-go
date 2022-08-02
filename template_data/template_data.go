package template_data

import "errors"

var Whitelist_states_html_unknown string = `
<option value="1">Awaiting Whitelist</option>
<option value="2">Whitelisted</option>
<option value="3">Temporarily Banned</option>
<option value="4">Banned</option>
<option value="5">Not Enough Activity</option>
<option value="6">Not On The Discord</option>
<option value="7">Contact for more Information</option>
<option value="unknown" selected="selected">Unknown - Rectify</option>`

var Service_states_html_unknown string = `
<option value="1">Online</option>
<option value="2">Experiencing Issues</option>
<option value="3">Undergoing Maintenance</option>
<option value="4" selected="selected">Offline</option>
<option value="unknown" selected="selected">Unknown - Rectify</option>`

var Skin_change_states_html_unknown string = "Contact the site administrator."
var Cape_change_states_html_unknown string = "Contact the site administrator."

var Whitelist_states_english_unknown string = "Unknown"

var Service_states_english_unknown string = "Unknown"

var Whitelist_states_html map[int]string = map[int]string{
	1: `
    <option value="1" selected="selected">Awaiting Whitelist</option>
    <option value="2">Whitelisted</option>
    <option value="3">Temporarily Banned</option>
    <option value="4">Banned</option>
	<option value="5">Not Enough Activity</option>
	<option value="6">Not On The Discord</option>
	<option value="7">Contact for more Information</option>`,
	2: `
    <option value="1">Awaiting Whitelist</option>
    <option value="2" selected="selected">Whitelisted</option>
    <option value="3">Temporarily Banned</option>
    <option value="4">Banned</option>
	<option value="5">Not Enough Activity</option>
	<option value="6">Not On The Discord</option>
	<option value="7">Contact for more Information</option>`,
	3: `
    <option value="1">Awaiting Whitelist</option>
    <option value="2">Whitelisted</option>
    <option value="3" selected="selected">Temporarily Banned</option>
    <option value="4">Banned</option>
	<option value="5">Not Enough Activity</option>
	<option value="6">Not On The Discord</option>
	<option value="7">Contact for more Information</option>`,
	4: `
    <option value="1">Awaiting Whitelist</option>
    <option value="2">Whitelisted</option>
    <option value="3">Temporarily Banned</option>
    <option value="4" selected="selected">Banned</option>
	<option value="5">Not Enough Activity</option>
	<option value="6">Not On The Discord</option>
	<option value="7">Contact for more Information</option>`,
	5: `
    <option value="1">Awaiting Whitelist</option>
    <option value="2">Whitelisted</option>
    <option value="3">Temporarily Banned</option>
    <option value="4">Banned</option>
	<option value="5" selected="selected">Not Enough Activity</option>
	<option value="6">Not On The Discord</option>
	<option value="7">Contact for more Information</option>`,
	6: `
    <option value="1">Awaiting Whitelist</option>
    <option value="2">Whitelisted</option>
    <option value="3">Temporarily Banned</option>
    <option value="4">Banned</option>
	<option value="5">Not Enough Activity</option>
	<option value="6" selected="selected">Not On The Discord</option>
	<option value="7">Contact for more Information</option>`,
	7: `
    <option value="1">Awaiting Whitelist</option>
    <option value="2">Whitelisted</option>
    <option value="3">Temporarily Banned</option>
    <option value="4">Banned</option>
	<option value="5">Not Enough Activity</option>
	<option value="6">Not On The Discord</option>
	<option value="7" selected="selected">Contact for more Information</option>`,
}

var Service_states_html map[int]string = map[int]string{
	1: `
	<option value="1" selected="selected">Online</option>
	<option value="2">Experiencing Issues</option>
	<option value="3">Undergoing Maintenance</option>
	<option value="4">Offline</option>`,
	2: `
	<option value="1">Online</option>
	<option value="2" selected="selected">Experiencing Issues</option>
	<option value="3">Undergoing Maintenance</option>
	<option value="4">Offline</option>`,
	3: `
	<option value="1">Online</option>
	<option value="2">Experiencing Issues</option>
	<option value="3" selected="selected">Undergoing Maintenance</option>
	<option value="4">Offline</option>`,
	4: `
	<option value="1">Online</option>
	<option value="2">Experiencing Issues</option>
	<option value="3">Undergoing Maintenance</option>
	<option value="4" selected="selected">Offline</option>`,
}

var Skin_change_states_html map[int]string = map[int]string{
	1: "You must be whitelisted in order to change your skin.",
	2: `You can create your own skin, or download one from any of the skins databases created by the community.<br /><br />To create your own skin, first <a href="/static/img/skin_reference.png">download the reference skin</a>. That's your template for creation. Then edit it with your favourite image editing tool.<br /><br />
	The skin can either be 64x32 or 64x64, whichever one you prefer. However, 64x64 skins have been known to not display properly on Minecraft 1.2.5 due to the lack of the two-pixel skin models.
	<br /><br />
	<form action="/profile/skin_submit" method="POST" enctype="multipart/form-data">
		<input type="file" id="upload" name="upload" accept="image/png">
		<br />
		<br />
		<button type="submit">Upload</button>
	</form>`,
	3: "Please wait until your ban expires before changing your skin.",
	4: "Banned users are not allowed to change their skin.",
	5: "Changing your skin is only available to whitelisted users. Please be more active in the <a href=\"https://discord.gg/hXTehk4PeG\">Action Retro Discord Guild</a>.",
	6: "Changing your skin is only available to whitelisted users. Please join the <a href=\"https://discord.gg/hXTehk4PeG\">Action Retro Discord Guild</a>.",
	7: "Changing your skin is only available to whitelisted users. Please contact us for further information (you'll need to be a member of the <a href=\"https://discord.gg/hXTehk4PeG\">Action Retro Discord Guild</a>).",
}

var Cape_change_states_html map[int]string = map[int]string{
	1: "You must be whitelisted in order to change your cape.",
	2: `You can create your own cape, or download one from any of the cape databases created by the community.<br /><br />To create your own cape, first <a href="/static/img/cape_reference.png">download the reference cape</a>. That's your template for creation. Then edit it with your favourite image editing tool.<br /><br />
	The cape must be 64x32, otherwise it will be rejected.
	<br /><br />
	<form action="/profile/cape_submit" method="POST" enctype="multipart/form-data">
		<input type="file" id="upload" name="upload" accept="image/png">
		<br />
		<br />
		<button type="submit">Upload</button>
	</form>`,
	3: "Please wait until your ban expires before changing your cape.",
	4: "Banned users are not allowed to change their cape.",
	5: "Changing your cape is only available to whitelisted users. Please be more active in the <a href=\"https://discord.gg/hXTehk4PeG\">Action Retro Discord Guild</a>.",
	6: "Changing your cape is only available to whitelisted users. Please join the <a href=\"https://discord.gg/hXTehk4PeG\">Action Retro Discord Guild</a>.",
	7: "Changing your cape is only available to whitelisted users. Please contact us for further information (you'll need to be a member of the <a href=\"https://discord.gg/hXTehk4PeG\">Action Retro Discord Guild</a>).",
}

var Whitelist_states_english map[int]string = map[int]string{
	0: "All",
	1: "Awaiting Whitelist",
	2: "Whitelisted",
	3: "Temp-Banned",
	4: "Banned",
	5: "Not Enough Activity",
	6: "Not On The Discord",
	7: "Please contact us for more information.",
}

var Service_states_english map[int]string = map[int]string{
	1: "Online",
	2: "Experiencing Issues",
	3: "Undergoing Maintenance",
	4: "Offline",
}

var Errors_english map[string]string = map[string]string{
	"change_state_no_state":           "No state provided on CHANGE_STATE.",
	"change_state_no_id":              "No ID provided on CHANGE_STATE.",
	"no_permission_whitelist":         "You don't have permission to modify user whitelist levels.",
	"no_permission_status":            "You don't have permission to modify service status.",
	"no_redirect_security_risk":       "A redirector was not specified in the request. Allowing any page to infer a state change is a security risk, and is disabled.",
	"service_state_no_state":          "No service state provided on CHANGE_STATE.",
	"service_state_no_ip":             "No service IP provided on CHANGE_STATE",
	"unspecified_state_error_service": "An unspecified error occurred while processing your service state change request. Please report this to the site administrator.",
	"unspecified_state_error_user":    "An unspecified error occurred while processing your user state change request. Please report this to the site administrator.",
	"user_creation_error":             "An unexpected error occurred while creating your user. Please contact site administration.",
	"no_permission_status_view":       "You don't have permission to view service status in this way. Please go <a href=\"/status\">here</a> if you wish to see service status.",
	"oauth2_denied":                   "An error was encountered while starting the OAuth2 flow. Chances are you clicked \"Cancel\" on the authorisation prompt. If not, please contact the site administrator.",
	"user_trans_prep_error":           "An error occurred while preparing the database transaction. Your changes may not have been applied. Please contact site administration.",
	"user_trans_exec_error":           "An error occurred while executing the database transaction. Your changes may not have been applied. Please contact site administration.",
	"user_exists_update":              "You tried to change your account username, but the username you are trying to use already exists on ActionMC. Please use another username.",
	"user_exists_create":              "You tried to register an account, but the username you used already exists on ActionMC. Please use another username.",
	"user_trans_init_error":           "An error occurred while initialising the database transaction. Your changes may not have been saved. Please contact site administration.",
	"del_no_confirm":                  "You tried to delete your account, but provided no confirmation code. Please go back and provide a valid confirmation code.",
	"del_invalid_confirm":             "You tried to delete your account, but provided an invalid confirmation code. Please go back, refresh, and re-enter the new confirmation code.",
	"del_trans_init_error":            "An error occurred while initialising the database transaction. Your changes may not have been saved, and your account may not be deleted. Please contact site administration.",
	"del_trans_prep_error":            "An error occurred while preparing the database transaction. Your changes may not have been saved, and your account may not be deleted. Please contact site administration.",
	"del_trans_exec_error":            "An error occurred while executing the database transaction. Your changes may not have been saved, and your account may not be deleted. Please contact site administration.",
	"bandel_trans_init_error":         "An error occurred while initialising the database transaction. Your changes may not have been saved, and your account may not be deleted. Please contact site administration.",
	"bandel_trans_prep_error":         "An error occurred while preparing the database transaction. Your changes may not have been saved, and your account may not be deleted. Please contact site administration.",
	"bandel_trans_exec_error":         "An error occurred while executing the database transaction. Your changes may not have been saved, and your account may not be deleted. Please contact site administration.",
	"banned":                          "You deleted your account while you were banned from ActionMC. Due to this, you are no longer allowed to register. The reason for this is to protect the server and our users from malicious users. If you wish to appeal your ban, please contact site administration.",
	"user_time_trans_init_error":      "An error occurred while initialising the database transaction. Your changes may not have been applied. Please contact site administration.",
	"user_time_trans_prep_error":      "An error occurred while preparing the database transaction. Your changes may not have been applied. Please contact site administration.",
	"user_time_trans_exec_error":      "An error occurred while executing the database transaction. Your changes may not have been applied. Please contact site administration.",
	"user_change_ratelimit":           "You can only change your username once every 30 days for safety reasons. If you require a name change <strong>now</strong>, please contact the site administration.",
	"no_permission_superuser":         "You do not have permission to access this superuser page.",
	"2fa_no_confirm":                  "You tried to enable 2FA, but didn't pass a token in your confirmation request. Please insert a token to confirm that you wish to enable 2FA.",
	"2fa_invalid_token":               "You entered an invalid token while trying to enable 2FA. Please insert a valid token generated by your authenticator.<br/><br/>If your token <strong>still</strong> doesn't work, check the date/time settings on your authenticator device.",
	"2fa_corrupt_token":               "Your token got corrupted during the confirmation process. Remove the entry for ActionMC from your authenticator and try again. If this persists, contact site administration.",
	"mfa_trans_init_error":            "An error occurred while initialising the database transaction. Your changes may not have been applied. Please contact site administration.",
	"mfa_trans_prep_error":            "An error occurred while preparing the database transaction. Your changes may not have been applied. Please contact site administration.",
	"mfa_trans_exec_error":            "An error occurred while executing the database transaction. Your changes may not have been applied. Please contact site administration.",
	"del_2fa_invalid_token":           "An invalid 2FA token was provided while attempting to delete your account. Please insert a valid token generated by your authenticator.<br/><br/>If your token <strong>still</strong> doesn't work, check the date/time settings on your authenticator device.",
	"del_no_2fa":                      "You tried to delete your 2FA-enabled account, but didn't pass a 2FA token in your deletion request. Please insert a valid 2FA token to confirm that you wish to delete your account.",
	"mfa_no_token":                    "You tried to verify your 2FA token, but didn't pass a token in your confirmation request. Please insert a token to confirm your account-impacting action.",
	"mfa_corrupt_key":                 "Your token got corrupted during the confirmation process. Please try again. If this persists, contact site administration.",
	"mfa_invalid_token":               "You entered an invalid token while trying to confirm an account-impacting action. Please insert a valid token generated by your authenticator.<br/><br/>If your token <strong>still</strong> doesn't work, check the date/time settings on your authenticator device.",
	"2fa_blocked":                     "Multi-factor authentication is currently disabled. We do not have an estimated timeframe for 2FA working.",
	"oauth2_error":                    "An unexpected error occurred during the OAuth2 flow. Please try logging in again - if it still doesn't work, contact the site administration.",
	"error_skin_upload":               "An error occurred while uploading your skin. Please check that it is a valid PNG file with a valid resolution and try again. If that doesn't work, please contact site adminstration.",
	"upload_too_large":                "The file you tried to upload is too large! We only support uploading PNG files with a maximum file size of 8 kilobytes.",
	"upload_bad_file":                 "You uploaded a bad file! Please verify that you uploaded a valid PNG file and try again. If it still doesn't work, please contact site administration.",
	"skin_bad_resolution":             "The skin you tried to upload is not a valid resolution. Skins must be either 64x64 or 64x32. Please adjust your skin to fit these dimensions and try again.",
	"error_cape_upload":               "An error occurred while uploading your cape. Please check that it is a valid PNG file with a valid resolution and try again. If that doesn't work, please contact site support.",
	"cape_bad_resolution":             "The cape you tried to upload is not a valid resolution. Capes must be 64x32. PLease adjust your cape to fit these dimensions and try again.",
	"consent_withdraw_db_error":       "There was an error removing your consent from the consent database. Please contact site administration.",
	"cookie_consent_no_preference":    "You have not allowed Preference cookies when you consented to cookies. This means you cannot use this part of the service. Please amend your cookie consent to allow Preference cookies and try again.",
	"reg_disabled":                    "Registration is currently disabled. Please try again later.",
	"image_gallery_not_found":         "Could not find image in the gallery - please contact site administration and try again later.",
}

var Accessibility_state_map = [2]string{
	`<option value="1" autocomplete="off">Enabled</option>
	<option value="0" selected="selected" autocomplete="off">Disabled</option>`,
	`<option value="1" selected="selected" autocomplete="off">Enabled</option>
	<option value="0">Disabled</option>`,
}

var Dynconf_state_map = [2]string{
	`<option value="1" autocomplete="off">Enabled</option>
	<option value="0" selected="selected" autocomplete="off">Disabled</option>`,
	`<option value="1" selected="selected" autocomplete="off">Enabled</option>
	<option value="0">Disabled</option>`,
}

var Admin_Strings map[int]string = map[int]string{
	0: "",
	1: "<h2>User Administrator Access</h2><a href=\"/admin/user\">Click here to access the user administration panel.</a><br /><br />",
	2: "<h2>User Administrator Access</h2><a href=\"/admin/user\">Click here to access the user administration panel.</a><br /><br /><h2>Status Administrator Access</h2><a href=\"/admin/status\">Click here to access the status administration panel.</a><br /><br />",
	3: "<h2>User Administrator Access</h2><a href=\"/admin/user\">Click here to access the user administration panel.</a><br /><br /><h2>Status Administrator Access</h2><a href=\"/admin/status\">Click here to access the status administration panel.</a><br /><br />",
	4: "<h2>User Administrator Access</h2><a href=\"/admin/user\">Click here to access the user administration panel.</a><br /><br /><h2>Status Administrator Access</h2><a href=\"/admin/status\">Click here to access the status administration panel.</a><br /><br />",
	5: "<h2>User Administrator Access</h2><a href=\"/admin/user\">Click here to access the user administration panel.</a><br /><br /><h2>Status Administrator Access</h2><a href=\"/admin/status\">Click here to access the status administration panel.</a><br /><br /><h2>Superuser Site Control Access</h2><a href=\"/admin/superuser\">Click here to access the superuser site control panel.</a><br /><br />",
}

var Admin_profile_strings map[int]string = map[int]string{
	0: "",
	1: "<h2>Administrator Panel</h2><a href=\"/admin\">Click here to access the administration panels.</a><br /><br />",
	2: "<h2>Administrator Panel</h2><a href=\"/admin\">Click here to access the administration panels.</a><br /><br />",
	3: "<h2>Administrator Panel</h2><a href=\"/admin\">Click here to access the administration panels.</a><br /><br />",
	4: "<h2>Administrator Panel</h2><a href=\"/admin\">Click here to access the administration panels.</a><br /><br />",
	5: "<h2>Administrator Panel</h2><a href=\"/admin\">Click here to access the administration panels.</a><br /><br />",
}

func Get_ErrorString(errType string) (error, string) {
	error_english, ok := Errors_english[errType]

	if !ok {
		// error doesn't exist
		return errors.New("Error not found in the Errors hashmap."), ""
	}

	return nil, error_english
}

func Get_StateStrings(state int) (error, string, string) {
	state_html, ok := Whitelist_states_html[state]

	if !ok {
		// something happened
		return errors.New("State not found in HTML hashmap."), Whitelist_states_html_unknown, Whitelist_states_english_unknown
	}

	state_english, eng_ok := Whitelist_states_english[state]

	if !eng_ok {
		// something happened
		return errors.New("State not found in English hashmap."), Whitelist_states_html_unknown, Whitelist_states_english_unknown
	}

	return nil, state_html, state_english

}

func Get_ServiceStrings(state int) (error, string, string) {
	state_html, ok := Service_states_html[state]

	if !ok {
		return errors.New("State not found in HTML hashmap"), Service_states_html_unknown, Service_states_english_unknown
	}

	state_english, eng_ok := Service_states_english[state]

	if !eng_ok {
		return errors.New("State not found in English hashmap"), Service_states_html_unknown, Service_states_english_unknown
	}

	return nil, state_html, state_english
}

func Get_ProfileStrings(state int) (error, string, string, string) {
	skin_html, ok := Skin_change_states_html[state]

	if !ok {
		return errors.New("State not found in skin HTML hashmap."), Skin_change_states_html_unknown, Cape_change_states_html_unknown, Whitelist_states_english_unknown
	}

	cape_html, cape_ok := Cape_change_states_html[state]

	if !cape_ok {
		return errors.New("State not found in skin HTML hashmap."), Skin_change_states_html_unknown, Cape_change_states_html_unknown, Whitelist_states_english_unknown
	}

	state_english, eng_ok := Whitelist_states_english[state]

	if !eng_ok {
		// something happened
		return errors.New("State not found in English hashmap."), Whitelist_states_html_unknown, Whitelist_states_english_unknown, Whitelist_states_english_unknown
	}

	return nil, skin_html, cape_html, state_english
}

func Get_AdminStrings(priv_level int) (error, string) {
	admin_strings, ok := Admin_profile_strings[priv_level]

	if !ok {
		return errors.New("Invalid privilege level given."), ""
	}

	return nil, admin_strings
}

func Get_AdminPageStrings(priv_level int) (error, string) {
	admin_strings, ok := Admin_Strings[priv_level]

	if !ok {
		return errors.New("Invalid privilege level given."), ""
	}

	return nil, admin_strings
}
