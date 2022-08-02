package auth

import (
	"git.dsrt-int.net/actionmc/actionmc-site-go/authdatabase"
	"git.dsrt-int.net/actionmc/actionmc-site-go/gdpr"
	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	"git.dsrt-int.net/actionmc/actionmc-site-go/oauth_handler"

	//"github.com/gofiber/fiber/v2"
	//	"github.com/gofiber/fiber/v2/middleware/session"
	"net/http"

	"git.dsrt-int.net/actionmc/actionmc-site-go/sessions"
)

var logger *logging.Logger = logging.New()

func InitialiseSession(sess *sessions.Session, db *authdatabase.MCAuthDB_sqlite3, user *oauth_handler.OAuth_User) bool {

	sess.Data.Username = user.Username
	sess.Data.Id = user.DiscordId

	dbInfo, ok := db.Getuser(user.DiscordId)

	if ok {
		logger.Debug.Println("User " + dbInfo.Name + " with Discord ID " + dbInfo.Uid + " and state " + dbInfo.State + " HAS AN ACCOUNT.")

		sess.Data.HasAccount = true
		sess.Data.MinecraftUsername = dbInfo.Name
		sess.Data.StateOnLogin = dbInfo.State
		sess.Data.LoggedIn = true

		return true
	} else {
		logger.Debug.Println("User with Discord ID " + user.DiscordId + " does NOT HAVE AN ACCOUNT.")

		sess.Data.HasAccount = false
		sess.Data.LoggedIn = true

		return false

	}

}

// Grabs a session from the session store and locks its mutex before returning
// Calling functions must call sess.Unlock or sess.Close, or a deadlock will occur on the next call to GetSession
func GetSession(store *sessions.Store, w http.ResponseWriter, r *http.Request, consent *gdpr.CookieConsent) (sess *sessions.Session) {
	// get session from storage
	sess, err := store.Get(w, r, consent)

	// if there is any form of error, panic. should be caught by middleware.
	if err != nil {
		panic(err)
	}

	// lock session before returning it
	// Recievers must call unlock after they are done with the session
	sess.Lock()

	return
}

func EndSession(sess *sessions.Session) bool {
	// destroy session, redirect user to home.
	/*if err := sess.Destroy(); err != nil {
		panic(err)
	}

	return true*/
	sess.Destroy()
	return true
}

func CheckAuthentication(sess *sessions.Session) (is_loggedIn bool, has_account bool) {
	if sess == nil {
		is_loggedIn = false
		has_account = false
		return
	}

	if !sess.Data.LoggedIn {
		// has account.
		is_loggedIn = false
		has_account = false
	} else {
		is_loggedIn = true

		has_account = sess.Data.HasAccount
	}

	return
}
