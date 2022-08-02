package template

import (
	"testing"
	"time"

	"git.dsrt-int.net/actionmc/actionmc-site-go/config"
	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	"git.dsrt-int.net/actionmc/actionmc-site-go/sessions"
)

var dummyTemplate = new(TemplateCollection)

/*
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
*/

var session = &sessions.Session{
	CreateTime: int(time.Now().Unix()),
	SessionId:  "test",
	Data:       new(sessions.SessionFields),
}

func TestTemplateCreate(t *testing.T) {
	dummyTemplate = NewTemplateCollection("data/templates", "data/bases", logging.New())

	// reset again.
	dummyTemplate = new(TemplateCollection)

}

func TestConfigInitialisation(t *testing.T) {
	dummyTemplate.logger = logging.New()
	dummyTemplate.baseTemplates.logger = logging.New()

	// site config initialisation
	dummyTemplate.baseTemplates.site_configuration = config.InitialiseConfig()
	dummyTemplate.site_configuration = config.InitialiseConfig()

}

func TestTemplateLoad(t *testing.T) {
	dummyTemplate.LoadTemplates("data/templates", "data/bases")
}

func BenchmarkTemplateCreate(b *testing.B) {

	for i := 0; i < 50; i++ {
		for k := range dummyTemplate.cachedTemplates.gdprTemplates {
			// cannot test non GDPR due to
			dummyTemplate.BuildPopulatedTemplate(k, session)
		}
	}

}

// tests complete.
