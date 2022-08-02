package gdpr

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"strings"

	//"io"
	"log"
	"net/http"
	"time"

	"git.dsrt-int.net/actionmc/actionmc-site-go/authdatabase"
	"github.com/nu7hatch/gouuid"
)

type CookieConsent struct {
	Exists           bool
	Version          int
	AllowNecessary   bool  // must be true otherwise they cant use the site lmao
	AllowPreferences bool  // optional, just means site settings won't persist.
	ConsentTimestamp int64 // unix timestamp of when user consented.
	ConsentID        string

	// as we have no marketing or analytics cookies, we dont need to worry about
	// asking users whether or not they want this data collected.
}

type ConsentHandler struct {
	ConsentVersion int
	sqlite3db      *authdatabase.MCAuthDB_sqlite3
}

// uuids on this end are literally just UUIDs with a sha1 sum.
func GenerateUUID() string {
	u, err := uuid.NewV4()

	if err != nil {
		log.Panic(err)
	}

	h := sha1.New()

	h.Write([]byte(u.String()))

	bs := hex.EncodeToString(h.Sum(nil))

	if err != nil {
		log.Panic(err)
	}

	return string(bs)
}

func (CH *ConsentHandler) GetConsentObject(r *http.Request) *CookieConsent {
	cookie, err := r.Cookie("AMC_CONSENT")

	if err != nil || len(cookie.String()) == 0 {
		// cookie doesnt exist so user HASNT CONSENTED!
		return new(CookieConsent)
	}

	decoded_hex, err := hex.DecodeString(cookie.Value)

	if err != nil {
		// cookie is corrupt so re-roll the cookie.
		//log.Println("Cookie re-roll during hex decode.", err)
		return new(CookieConsent)
	}

	a := bytes.NewReader(decoded_hex)

	dec := gob.NewDecoder(a)

	var consentData CookieConsent

	err = dec.Decode(&consentData)

	if err != nil {
		// cookie is corrupt so re-roll the cookie.
		//log.Println("Cookie re-roll during decoding.")
		return new(CookieConsent)
	}

	if consentData.Version != CH.ConsentVersion {
		// consent is out of date, get them to re-roll.
		return new(CookieConsent)
	}

	return &consentData
}

func (CH *ConsentHandler) NewConsentObject(allow_necessary bool, allow_preferences bool) string {
	//uuid := GenerateUUID()

	consent := CookieConsent{
		Exists:           true,
		Version:          CH.ConsentVersion,
		AllowNecessary:   allow_necessary,
		AllowPreferences: allow_preferences,
		ConsentTimestamp: time.Now().Unix(),
		ConsentID:        GenerateUUID(),
	}

	// still gotta log these in database for reasons (GDPR is annoying)
	CH.sqlite3db.AddGDPR(consent.ConsentID, consent.Version, allow_necessary, allow_preferences, consent.ConsentTimestamp)

	b := new(strings.Builder)

	enc := gob.NewEncoder(hex.NewEncoder(b))
	err := enc.Encode(consent)

	if err != nil {
		// TODO (fnuer): maybe here instead of crashing, maybe actually
		// handle the error properly, as in just force a re-roll of consent?

		// This'd probably lead to users getting trapped in loops if the
		// encode fails, although this really shouldn't happen.

		panic("Site literally cannot run due to failed encode of GDPR.")
	}

	return b.String()
}

func (CC *CookieConsent) String() string {
	nec := "No"
	pre := "No"
	exi := "No"

	if CC.AllowNecessary {
		nec = "Yes"
	}

	if CC.AllowPreferences {
		pre = "Yes"
	}

	if CC.Exists {
		exi = "Yes"
	}

	return nec + " " + pre + " " + exi
}

func NewConsentHandler(sqlite3db *authdatabase.MCAuthDB_sqlite3) *ConsentHandler {
	return &ConsentHandler{sqlite3db: sqlite3db, ConsentVersion: 2}
}
