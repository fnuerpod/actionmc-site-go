package dynamic_cfg

import (
	"strconv"
	"testing"
)

var dynconf *DynamicConfig

func TestConfigInit(t *testing.T) {
	dynconf = InitialiseDynConfig()
}

func TestConfigEdit(t *testing.T) {
	dynconf.AnnounceBannerText = "test note"
	dynconf.ActionMGMT_LatestVersion = "1.0.2T"
	dynconf.MaintenanceMode = true
	dynconf.AccountCreationDisabled = true
	dynconf.ShowAnnounceBanner = true
}

func TestEditedConfigSave(t *testing.T) {
	dynconf.Save()
}

func TestConfigReset(t *testing.T) {
	dynconf.AnnounceBannerText = ""
	dynconf.ActionMGMT_LatestVersion = ""
	dynconf.MaintenanceMode = false
	dynconf.AccountCreationDisabled = false
	dynconf.ShowAnnounceBanner = false
}

func TestResetConfigSave(t *testing.T) {
	dynconf.Save()
}

func BenchmarkConfigEditSave(b *testing.B) {
	for i := 0; i < 500; i++ {
		a := strconv.Itoa(i)

		dynconf.AnnounceBannerText = a
		dynconf.ActionMGMT_LatestVersion = a

		dynconf.MaintenanceMode = (i%2 == 0)
		dynconf.AccountCreationDisabled = (i%2 == 1)
		dynconf.ShowAnnounceBanner = (i%2 == 0)

		dynconf.Save()
	}
}
