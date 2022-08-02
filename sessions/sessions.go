package sessions

import (
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"git.dsrt-int.net/actionmc/actionmc-site-go/gdpr"

	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	"github.com/nu7hatch/gouuid"
)

var logger *logging.Logger = logging.New()

// initialise typedefs

// container is a map of all sessions by string, which is session ID.
type SessionContainer map[string]*Session

// sessiondata stores all session data by string, which is session ID.
type SessionData map[string]interface{}

// stores the session container in a more user friendly manner.
// also locks map to prevent concurrent r/w
type Store struct {
	mu       *sync.Mutex // mutex to lock map access.
	sessions SessionContainer
}

// session uuid generation and verification.

func (store *Store) genuuid() (error, string) {
	u, err := uuid.NewV4()

	if err != nil {
		return err, ""
	}

	u_two, err := uuid.NewV4()

	if err != nil {
		return err, ""
	}

	result := u.String() + "-" + u_two.String()

	store.mu.Lock()

	_, ok := store.sessions[result]

	store.mu.Unlock()

	if ok {
		// uh... what.
		logger.Debug.Println("Session ID already exists. trying again...")

		// this could be dangerous and lead to like...
		// so MUCH recursion but i dont care.
		_, result = store.genuuid()
	}

	return nil, result
}

// csrf prevention uuid generation
func GenerateCSRFString() string {
	u, err := uuid.NewV4()

	if err != nil {
		return ""
	}

	return u.String()
}

// initialise session.
func (store *Store) InitSession(w http.ResponseWriter, r *http.Request, consent *gdpr.CookieConsent) string {
	// generate new and unique session id.
	err, session_id := store.genuuid()

	// this shouldn't error out but like edgecases and such.
	if err != nil {
		log.Panic(err)
	}

	logger.Debug.Println("New session initialising...")

	session_data := new(SessionFields)

	// we now have session data. now we check for cookies that say
	// "accessibility CSS required"

	cookie, err := r.Cookie("accessibility")

	if err != nil || len(cookie.String()) == 0 {
		// cookie doesn't exist so we need to build.

		// check gdpr

		if consent.AllowPreferences {
			SetAccessibilityCookie(w, "00")
		}

		// if they didn't allow preferences said cookie wont be set.

	}

	// now parse cookie.
	var parsed_acc_obj AccessibilityData

	// Using a reader over the cookie string with a hex decoder and through bitset
	if s := cookie.String(); len(s) == 2 {

		var bitset []byte

		bitset, err = hex.DecodeString(s)

		if err != nil {
			goto errorState
		}

		bits := AccessibilityBitset(bitset[0])

		parsed_acc_obj.EnableDyslexia = bits.Get(EnableDyslexia)
		parsed_acc_obj.EnableHiContrast = bits.Get(EnableHiContrast)
		parsed_acc_obj.EnableNoImage = bits.Get(EnableNoImage)

	} else {
		err = errors.New("Could not decode Accessibility")
	}

errorState:
	if err != nil {
		// check gdpr

		if consent.AllowPreferences {
			SetAccessibilityCookie(w, "00")
		}

		// if they didn't allow preferences said cookie wont be set.
	}

	// Place AccessibilityData in session

	session_data.AccessibilityData = parsed_acc_obj

	new_session := &Session{
		CreateTime: int(time.Now().Unix()),
		SessionId:  session_id,
		Data:       session_data,
	}

	// lock access to map.
	store.mu.Lock()

	// add new session data to the session storage map
	store.sessions[session_id] = new_session

	// unlock map access.
	store.mu.Unlock()

	// create new cookie and set it on the client response.

	c := &http.Cookie{
		Name:  "AMC_SESSION",
		Path:  "/",
		Value: session_id,
	}

	http.SetCookie(w, c)

	// return session id to the calling goroutine.
	return session_id
}

// get session from the store.
func (store *Store) Get(w http.ResponseWriter, r *http.Request, consent *gdpr.CookieConsent) (*Session, error) {
	// check if a cookie exists.
	cookie, err := r.Cookie("AMC_SESSION")

	// string_sess wil contain the cookie value...
	// OR the generated session ID by the session initialiser.
	var string_sess string

	if err != nil {
		// no cookie exists so session needs initialised.
		logger.Debug.Println("No session cookie, initialising...")
		string_sess = store.InitSession(w, r, consent)
	} else {
		// session exists.
		string_sess = cookie.Value
	}

	// lock the store's mutex.
	store.mu.Lock()

	// get session, if possible.
	session_check, ok := store.sessions[string_sess]

	// unlock mutex as we're done with it for now.
	store.mu.Unlock()

	if !ok {
		// session doesn't exist so we can initialise with this id.
		logger.Debug.Println("Session initialisation called from Get function in Store struct.")

		// when initialising session, it will return the new session ID.
		// we need this for verifying session existance.
		string_sess = store.InitSession(w, r, consent)

		// lock the mutex again.
		store.mu.Lock()

		// do a quick check to see if the session actually exists.
		// if not, it wasn't initialised at all (what??)
		session_check, ok = store.sessions[string_sess]

		// once again unlock mutex.
		store.mu.Unlock()

		if !ok {
			logger.Err.Println("Initialisation of session failed - session initialisation function failed to commit new session.")
		}

		// return the session_check variable as that will contain the new session.
		return session_check, nil
	}

	// the session already exists. just return it.

	return session_check, nil

}

type SessionFields struct {
	// Session fields
	HasAccount        bool
	Username          string
	Id                string
	StateOnLogin      string
	MinecraftUsername string
	LoggedIn          bool
	DeletionString    string

	// Redirect fields
	CurrentRedirect string

	// Accessibility fields
	AccessibilityData AccessibilityData

	// CSRF prevention fields
	CSRFPreventionString string
}

// actual session object.
type Session struct {
	CreateTime int
	SessionId  string
	Data       *SessionFields
	sync.Mutex
}

// Session.Close satisfies the io.Closer interface and unlocks the session Lock, it always returns a nil error
func (session *Session) Close() error {
	session.Unlock()
	return nil
}

// resets/"destroys" a session object.
func (session *Session) Destroy() {

	// Clear session fields
	session.Data = new(SessionFields)

}

func New() *Store {
	return &Store{
		mu:       new(sync.Mutex),
		sessions: make(SessionContainer),
	}
}
