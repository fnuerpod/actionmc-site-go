package template

import (
	"embed"
	"fmt"
	"html"

	//"io/fs"
	"log"

	//"os"
	"strings"
	"time"

	"git.dsrt-int.net/actionmc/actionmc-site-go/config"
	"git.dsrt-int.net/actionmc/actionmc-site-go/dynamic_cfg"

	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	"git.dsrt-int.net/actionmc/actionmc-site-go/sessions"
)

//go:embed data/*

var content embed.FS

type BaseTemplateCollection struct {
	templates          map[string]string
	logger             *logging.Logger
	site_configuration *config.Configuration
}

func (BT *BaseTemplateCollection) GetTemplate(name string) string {
	data, ok := BT.templates[name]

	if !ok {
		panic("Template not found in base templates.")
	}

	return data
}

func (BT *BaseTemplateCollection) BuildBase() (string, string) {
	// thanks to the almighty genius of alex,
	// i shall be using a string builder here.
	var blen int
	var blen2 int

	head := BT.GetTemplate("head")
	body_header := BT.GetTemplate("body_header")
	body_footer := BT.GetTemplate("body_footer")

	blen += len("<html><head>") + len(head) + len("</head><body>") + len(body_header) + len("%%ANNOUNCEMENT%%")

	blen2 += len(body_footer) + len("</body></html>")

	var b strings.Builder
	var b2 strings.Builder
	b.Grow(blen)
	b2.Grow(blen2)

	b.WriteString("<html><head>")
	b.WriteString(head)
	b.WriteString("</head><body>")
	b.WriteString(body_header)
	b.WriteString("%%ANNOUNCEMENT%%")

	// after the content, use b2 builder.
	b2.WriteString(body_footer)
	b2.WriteString("</body></html>")

	// Ram wastage check
	// certified alex moment
	if b.Len() != blen {
		panic(b.Len() - blen)
	}

	if b2.Len() != blen2 {
		panic(b2.Len() - blen2)
	}

	// return...
	return b.String(), b2.String()
}

func (BT *BaseTemplateCollection) BuildGDPRBase() (string, string) {
	// thanks to the almighty genius of alex,
	// i shall be using a string builder here.
	var blen int
	var blen2 int

	head := BT.GetTemplate("head_gdpr")
	body_header := BT.GetTemplate("body_header_gdpr")
	body_footer := BT.GetTemplate("body_footer_gdpr")

	blen += len("<html><head>") + len(head) + len("</head><body>") + len(body_header) + len("%%ANNOUNCEMENT%%")

	blen2 += len(body_footer) + len("</body></html>")

	var b strings.Builder
	var b2 strings.Builder
	b.Grow(blen)
	b2.Grow(blen2)

	b.WriteString("<html><head>")
	b.WriteString(head)
	b.WriteString("</head><body>")
	b.WriteString(body_header)
	b.WriteString("%%ANNOUNCEMENT%%")

	// after the content, use b2 builder.
	b2.WriteString(body_footer)
	b2.WriteString("</body></html>")

	// Ram wastage check
	// certified alex moment
	if b.Len() != blen {
		panic(b.Len() - blen)
	}

	if b2.Len() != blen2 {
		panic(b2.Len() - blen2)
	}

	// return...
	return b.String(), b2.String()
}

func (BT *BaseTemplateCollection) BuildGDPR_NoBanner_Base() (string, string) {
	// thanks to the almighty genius of alex,
	// i shall be using a string builder here.
	var blen int
	var blen2 int

	head := BT.GetTemplate("head_gdpr")
	//body_header := BT.GetTemplate("")
	//body_footer := BT.GetTemplate("body_footer_gdpr")

	blen += len("<html><head>") + len(head) + len("</head><body>")

	blen2 += len("</body></html>")

	var b strings.Builder
	var b2 strings.Builder
	b.Grow(blen)
	b2.Grow(blen2)

	b.WriteString("<html><head>")
	b.WriteString(head)
	b.WriteString("</head><body>")

	// after the content, use b2 builder.
	b2.WriteString("</body></html>")

	// Ram wastage check
	// certified alex moment
	if b.Len() != blen {
		panic(b.Len() - blen)
	}

	if b2.Len() != blen2 {
		panic(b2.Len() - blen2)
	}

	// return...
	return b.String(), b2.String()
}

func (BT *BaseTemplateCollection) LoadBaseTemplates(dir string) bool {
	// empty out template hashmap.
	BT.templates = make(map[string]string)

	files, err := content.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	//logging.Debug_nn("\n")

	for _, f := range files {
		//BT.logger.Debug.Println("Loading base template for snippet \"" + f.Name() + "\" into memory...")
		file, err := content.ReadFile(dir + "/" + f.Name() + "")
		if err != nil {
			//logging.Debug_nn(" \033[31mFAIL\033[0m\n")
			BT.logger.Fatal.Fatalln("Failed to open base template file, this is fatal - program will terminate. More information below...")
			//log.Fatal(err)
		}

		//activeTemplates = append(activeTemplates, &Template{temp_name: strings.TrimSuffix(f.Name(), ".html"), temp_content: final_text})
		BT.templates[strings.TrimSuffix(f.Name(), ".b.html")] = string(file)
		BT.logger.Debug.Println("Loaded template for page \"" + f.Name() + "\" into memory.")
	}

	BT.logger.Debug.Println("All base templates loaded into memory OK.")

	return true
}

type CachedTemplateCollection struct {
	templates             map[string]string
	gdprTemplates         map[string]string
	gdprNoBannerTemplates map[string]string
}

func (CTC *CachedTemplateCollection) GenerateCache(TC *TemplateCollection) bool {
	// empty out cache hashmap.
	CTC.templates = make(map[string]string)
	CTC.gdprTemplates = make(map[string]string)
	CTC.gdprNoBannerTemplates = make(map[string]string)

	dynconf := dynamic_cfg.InitialiseDynConfig()

	for template_name, content := range TC.templates {
		// generate the built template for caching.
		base, base_ac := TC.baseTemplates.BuildBase()

		// start building the actual template
		joinstr := []string{
			base,
			content,
			base_ac,
		}

		result := strings.Join(joinstr, "")

		if dynconf.ShowAnnounceBanner {
			result = strings.Replace(result, "%%ANNOUNCEMENT%%", "<p class=\"announcement_banner\">"+dynconf.AnnounceBannerText+"</p>", 1)
		} else {
			result = strings.Replace(result, "%%ANNOUNCEMENT%%", "", 1)
		}

		// take result and put it into the hashmap
		CTC.templates[template_name] = result
	}

	// gdpr templates
	for template_name, content := range TC.templates {
		// generate the built template for caching.
		base, base_ac := TC.baseTemplates.BuildGDPRBase()

		// start building the actual template
		joinstr := []string{
			base,
			content,
			base_ac,
		}

		result := strings.Join(joinstr, "")

		if dynconf.ShowAnnounceBanner {
			result = strings.Replace(result, "%%ANNOUNCEMENT%%", "<p class=\"announcement_banner\">"+dynconf.AnnounceBannerText+"</p>", 1)
		} else {
			result = strings.Replace(result, "%%ANNOUNCEMENT%%", "", 1)
		}

		// take result and put it into the hashmap
		CTC.gdprTemplates[template_name] = result
	}

	// nobanner gdpr templates
	for template_name, content := range TC.templates {
		// generate the built template for caching.
		base, base_ac := TC.baseTemplates.BuildGDPR_NoBanner_Base()

		// start building the actual template
		joinstr := []string{
			base,
			content,
			base_ac,
		}

		result := strings.Join(joinstr, "")

		if dynconf.ShowAnnounceBanner {
			result = strings.Replace(result, "%%ANNOUNCEMENT%%", "<p class=\"announcement_banner\">"+dynconf.AnnounceBannerText+"</p>", 1)
		} else {
			result = strings.Replace(result, "%%ANNOUNCEMENT%%", "", 1)
		}

		// take result and put it into the hashmap
		CTC.gdprNoBannerTemplates[template_name] = result
	}

	return true
}

func (CTC *CachedTemplateCollection) GetCached(name string) string {
	data, ok := CTC.templates[name]

	if !ok {
		panic("Template " + name + " doesn't exist in the cache.")
	}

	return data
}

func (CTC *CachedTemplateCollection) GetGDPR_NoBanner_Cached(name string) string {
	data, ok := CTC.gdprNoBannerTemplates[name]

	if !ok {
		panic("Template " + name + " doesn't exist in the cache.")
	}

	return data
}

func (CTC *CachedTemplateCollection) GetGDPRCached(name string) string {
	data, ok := CTC.gdprTemplates[name]

	if !ok {
		panic("Template " + name + " doesn't exist in the cache.")
	}

	return data
}

type TemplateCollection struct {
	templates          map[string]string
	baseTemplates      BaseTemplateCollection
	cachedTemplates    CachedTemplateCollection
	logger             *logging.Logger
	site_configuration *config.Configuration
}

func (TC *TemplateCollection) LoadTemplates(dir string, dir2 string) bool {
	// empty out template hashmap.
	TC.templates = make(map[string]string)

	files, err := content.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	//logging.Debug_nn("\n")

	for _, f := range files {
		//TC.logger.Debug.Println("Loading template for page \"" + f.Name() + "\" into memory...")
		file, err := content.ReadFile(dir + "/" + f.Name() + "")
		if err != nil {
			//logging.Debug_nn(" \033[31mFAIL\033[0m\n")
			TC.logger.Fatal.Fatalln("Failed to open template file, this is fatal - program will terminate. More information below...")
			//log.Fatal(err)
		}

		//activeTemplates = append(activeTemplates, &Template{temp_name: strings.TrimSuffix(f.Name(), ".html"), temp_content: final_text})
		TC.templates[strings.TrimSuffix(f.Name(), ".html")] = string(file)
		TC.logger.Debug.Println("Loaded template for page \"" + f.Name() + "\" into memory.")
	}

	TC.logger.Debug.Println("All templates loaded into memory OK.")

	// now load the base templates
	TC.baseTemplates.LoadBaseTemplates(dir2)

	// now do the cache.
	TC.cachedTemplates.GenerateCache(TC)
	return true

}

func (TC *TemplateCollection) BuildTemplate(name string) string {
	return TC.cachedTemplates.GetCached(name)
}

func (TC *TemplateCollection) BuildGDPRTemplate(name string) string {
	return TC.cachedTemplates.GetGDPRCached(name)
}

func (TC *TemplateCollection) BuildGDPR_NoBanner_Template(name string) string {
	return TC.cachedTemplates.GetGDPR_NoBanner_Cached(name)
}

func (TC *TemplateCollection) BuildPopulatedGDPR_NoBanner_Template(temp_name string) (result string, timer int64) {
	startTime_template := time.Now().UnixNano()

	//TC.logger.Debug.Println("----------------------------------------")
	TC.logger.Debug.Println("Creating template for page " + temp_name + "...")
	result = TC.BuildGDPR_NoBanner_Template(temp_name)

	//TC.logger.Debug.Println("Populating base variables...")

	//result = strings.Replace(result, "%%CDN_URL%%", cdn, -1)

	//TC.logger.Debug.Println("Doing login checks in session...")

	var replacestr []string

	if TC.site_configuration.IsPreproduction {
		//TC.logger.Debug.Println("Pre-production, so apply preproduction template changes.")
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo_preprod.png",
			//"%%LOGO_SMALL_URL%%", "/static/img/logo_small_preprod.png",
			"%%PLAYER_OR_PREPROD%%", "PREPRODUCTION - DOES NOT REFLECT FINAL BUILD!",
		}
	} else if TC.site_configuration.IsBeta {
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo_beta.png",
			//"%%LOGO_SMALL_URL%%", "/static/img/logo_small_beta.png",
			"%%PLAYER_OR_PREPROD%%", "BETA - MAY CONTAIN BUGS!",
		}
	} else {
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo.png",
			//"%%LOGO_SMALL_URL%%", "/static/img/logo_small.png",
			"%%PLAYER_OR_PREPROD%%", "0/10 playing on ActionMC.",
		}
	}

	result = strings.NewReplacer(replacestr...).Replace(result)

	endTime_Template := time.Now().UnixNano()
	// Timestamp converted to float miliseconds that is rounded at 3 decimal places
	timer = (endTime_Template - startTime_template) / 1000

	TC.logger.Debug.Println("Created template for page " + temp_name + " in " + fmt.Sprint(float64(timer)/1000) + "ms.")
	//TC.logger.Debug.Println("----------------------------------------")
	return
}

func (TC *TemplateCollection) BuildPopulatedGDPRTemplate(temp_name string) (result string, timer int64) {
	startTime_template := time.Now().UnixNano()

	//TC.logger.Debug.Println("----------------------------------------")
	TC.logger.Debug.Println("Creating template for page " + temp_name + "...")
	result = TC.BuildGDPRTemplate(temp_name)

	//TC.logger.Debug.Println("Populating base variables...")

	//result = strings.Replace(result, "%%CDN_URL%%", cdn, -1)

	//TC.logger.Debug.Println("Doing login checks in session...")

	var replacestr []string

	if TC.site_configuration.IsPreproduction {
		//TC.logger.Debug.Println("Pre-production, so apply preproduction template changes.")
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo_preprod.png",
			"%%LOGO_SMALL_URL%%", "/static/img/logo_small_preprod.png",
			"%%PLAYER_OR_PREPROD%%", "PREPRODUCTION - DOES NOT REFLECT FINAL BUILD!",
		}
	} else if TC.site_configuration.IsBeta {
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo_beta.png",
			"%%LOGO_SMALL_URL%%", "/static/img/logo_small_beta.png",
			"%%PLAYER_OR_PREPROD%%", "BETA - MAY CONTAIN BUGS!",
		}
	} else {
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo.png",
			"%%LOGO_SMALL_URL%%", "/static/img/logo_small.png",
			"%%PLAYER_OR_PREPROD%%", "0/10 playing on ActionMC.",
		}
	}

	result = strings.NewReplacer(replacestr...).Replace(result)

	endTime_Template := time.Now().UnixNano()
	// Timestamp converted to float miliseconds that is rounded at 3 decimal places
	timer = (endTime_Template - startTime_template) / 1000

	TC.logger.Debug.Println("Created template for page " + temp_name + " in " + fmt.Sprint(float64(timer)/1000) + "ms.")
	//TC.logger.Debug.Println("----------------------------------------")
	return
}

func (TC *TemplateCollection) BuildPopulatedTemplate(temp_name string, fiber_session *sessions.Session) (result string, timer int64) {
	startTime_template := time.Now().UnixNano()

	//TC.logger.Debug.Println("----------------------------------------")

	//TC.logger.Debug.Println("Creating template for page " + temp_name + "...")
	result = TC.BuildTemplate(temp_name)

	accessibility_directives := ""
	// lets check their session for accessibility modes

	//fmt.Println("penis %v, %v, %v", fiber_session.Data.AccessibilityData.EnableDyslexia, fiber_session.Data.AccessibilityData.EnableHiContrast, fiber_session.Data.AccessibilityData.EnableNoImage)

	if fiber_session.Data.AccessibilityData.EnableDyslexia {
		// dyslexia accessibility mode ENABLED.
		accessibility_directives = accessibility_directives + "<link rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/a_dyslexia.css\" />"
	}

	if fiber_session.Data.AccessibilityData.EnableHiContrast {
		// high contrast accessibility mode ENABLED.
		accessibility_directives = accessibility_directives + "<link rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/a_hicontrast.css\" />"
	}

	if fiber_session.Data.AccessibilityData.EnableNoImage {
		// no image accessibility mode ENABLED.
		accessibility_directives = accessibility_directives + "<link rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/a_noimage.css\" />"
	}

	login_check := fiber_session.Data.Username

	var login_banner string = "Not logged in. [ <a href=\"/login\">Log in via Discord</a> ]"

	if login_check != "" {
		// user is logged in.

		// Escape the username strings before passing them to avoid a script injection

		if fiber_session.Data.HasAccount {
			//TC.logger.Debug.Println("User is logged in and has an account with us, populate base banner with their MINECRAFT username.")
			login_banner = html.EscapeString(fiber_session.Data.MinecraftUsername) + " [ <a href=\"/logout\">Log out</a> ]"
		} else {
			//TC.logger.Debug.Println("User is logged in and has NO account with us, populate base banner with their DISCORD username.")
			login_banner = html.EscapeString(fiber_session.Data.Username) + " [ <a href=\"/logout\">Log out</a> ]"
		}

	}

	var replacestr []string

	if TC.site_configuration.IsPreproduction {
		//TC.logger.Debug.Println("Pre-production, so apply preproduction template changes.")
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo_preprod.png",
			"%%LOGO_SMALL_URL%%", "/static/img/logo_small_preprod.png",
			"%%PLAYER_OR_PREPROD%%", "PREPRODUCTION - DOES NOT REFLECT FINAL BUILD! <a href=\"/preprod_info\">More info</a> ",
			"%%LOGIN_BANNER%%", login_banner,
			"%%ACC_CSS%%", accessibility_directives,
		}
	} else if TC.site_configuration.IsBeta {
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo_beta.png",
			"%%LOGO_SMALL_URL%%", "/static/img/logo_small_beta.png",
			"%%PLAYER_OR_PREPROD%%", "BETA - MAY CONTAIN BUGS! <a href=\"/beta_info\">More info</a> ",
			"%%LOGIN_BANNER%%", login_banner,
			"%%ACC_CSS%%", accessibility_directives,
		}
	} else {
		replacestr = []string{
			"%%LOGO_BIG_URL%%", "/static/img/logo.png",
			"%%LOGO_SMALL_URL%%", "/static/img/logo_small.png",
			"%%PLAYER_OR_PREPROD%%", "0/10 playing on ActionMC. <a href=\"/\">More info</a> ",
			"%%LOGIN_BANNER%%", login_banner,
			"%%ACC_CSS%%", accessibility_directives,
		}
	}

	result = strings.NewReplacer(replacestr...).Replace(result)

	endTime_Template := time.Now().UnixNano()
	// Timestamp converted to float miliseconds that is rounded at 3 decimal places
	timer = (endTime_Template - startTime_template) / 1000

	//TC.logger.Debug.Println("Created template for page " + temp_name + " in " + fmt.Sprint(float64(timer)/1000) + "ms.")
	//TC.logger.Debug.Println("----------------------------------------")
	return
}

func NewTemplateCollection(dir1 string, dir2 string, logger *logging.Logger) *TemplateCollection {
	var result TemplateCollection

	// logger initialisation
	result.logger = logger
	result.baseTemplates.logger = logger

	// site config initialisation
	result.baseTemplates.site_configuration = config.InitialiseConfig()
	result.site_configuration = config.InitialiseConfig()

	result.LoadTemplates(dir1, dir2)

	return &result
}

func GetTemplateTimestamp(templatetime, extrastart, extraend int64) string {

	extratimestamp := (extraend - extrastart) / 1000

	return "Template Creation: " + fmt.Sprint(float64(templatetime)/1000) + "ms | Additional Processing: " + fmt.Sprint(float64(extratimestamp)/1000) + "ms | Total: " + fmt.Sprint(float64(templatetime+extratimestamp)/1000) + "ms"

}
