package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// TODO (fnuer): move configuration to json parsed.

// prolly much more effective than parsing our own shitty format using
// inefficient and buggy tactics.

type Configuration struct {
	IsPreproduction           bool   `json:IsPreproduction`
	IsBeta                    bool   `json:IsBeta`
	BindPort                  string `json:BindPort`
	CallBackURL_Production    string `json:CallBackURL_Production`
	CallBackURL_PreProduction string `json:CallBackURL_PreProduction`
}

// Initialise default configuration.
var Config *Configuration = &Configuration{IsPreproduction: true, BindPort: "6483", IsBeta: false}

func GetDataDir() string {

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	configdir := filepath.Join(user.HomeDir, ".config", "actionmc-site-go")

	if _, err := os.Stat(configdir); os.IsNotExist(err) {
		// config dir not exist
		err = os.Mkdir(filepath.Join(user.HomeDir, ".config", "actionmc-site-go"), 0755)

		if err != nil {
			panic(err)
		}
	}

	return configdir
}

func InitialiseConfig() *Configuration {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Fatal exception when parsing the configuration data. CHECK YOUR CONFIG!", err)
		}
	}()

	config_data_location := filepath.Join(GetDataDir(), "config.json")

	if _, err := os.Stat(config_data_location); os.IsNotExist(err) {
		// config dir not exist

		log.Fatalln("Configuration doesn't exist! Please place your config in ~/.config/actionmc-site-go/config.json in order to start the site.")
	}

	if data, err := os.ReadFile(config_data_location); err == nil {
		// got file.
		var parsed Configuration

		// parse. panic on error.
		if err := json.Unmarshal(data, &parsed); err != nil {
			panic(err)
		}

		if parsed.IsPreproduction && parsed.IsBeta {
			log.Fatalln("Cannot be both in preproduction and beta mode.")
		}

		// return the parsed json.
		return &parsed
	} else if errors.Is(err, os.ErrNotExist) {
		// does not exist.
		log.Fatalln("Configuration doesn't exist! Please place your config in ~/.config/actionmc-site-go/config.json in order to start the site.")
	} else {
		// schrodingers file. presume non-existant.
		log.Fatalln("Configuration doesn't exist! Please place your config in ~/.config/actionmc-site-go/config.json in order to start the site.")
	}

	return Config
}
