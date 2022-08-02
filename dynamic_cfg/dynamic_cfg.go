package dynamic_cfg

import (
	"git.dsrt-int.net/actionmc/actionmc-site-go/config"
	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"

	//"bytes"
	"encoding/gob"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var logger *logging.Logger = logging.New()

type DynamicConfig struct {
	MaintenanceMode         bool
	AccountCreationDisabled bool

	ShowAnnounceBanner bool
	AnnounceBannerText string

	ActionMGMT_LatestVersion string
}

func InitialiseDynConfig() *DynamicConfig {
	fileHandle, _ := os.Open(filepath.Join(config.GetDataDir(), "dynconf.gob"))
	defer fileHandle.Close()

	dec := gob.NewDecoder(fileHandle)

	var dynconfig DynamicConfig

	err := dec.Decode(&dynconfig)

	if err != nil {
		logger.Err.Println("Error loading Dynamic Config... Creating new one and trying again...")
		dynconfig.Save()

	}

	return &dynconfig
}

func (DC *DynamicConfig) Save() {
	b := new(strings.Builder)

	enc := gob.NewEncoder(hex.NewEncoder(b))
	err := enc.Encode(DC)

	f, err := os.Create(filepath.Join(config.GetDataDir(), "dynconf.gob"))

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	t, err := hex.DecodeString(b.String())

	if err != nil {
		log.Fatal(err)
	}

	_, err2 := f.Write(t)

	if err2 != nil {
		log.Fatal(err2)
	}
}
