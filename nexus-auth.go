package main

import (
	"flag"
	"fmt"
	"github.com/praveenprem/logging"
	"github.com/praveenprem/nexus-auth/codes"
	"github.com/praveenprem/nexus-auth/configuration"
	"github.com/praveenprem/nexus-auth/notification"
	"github.com/praveenprem/nexus-auth/providers"
	"os"
	"strings"
)

var (
	//SIGNATURE defines a string that hold the application author's signature
	SIGNATURE = `
******************************************
*                                        *
*      ********         ********         *
*      **      **       **      **       *
*      **       **      **       **      *
*      **      **       **      **       *
*      ********         ********         *
*      **               **               *
*      **               **               *
*                                        *
******************************************
         ** Praveen Premaratne **

 * Project name: ProjectNeptune
 * Created by: Praveen Premaratne
 * Created on: 20/03/2020 22:38
`

	// VERSION is a string that hold the build version number passed during the build
	VERSION string
)

//App struct hold application configuration used throughout application
type App struct {
	User       string //KeyChain given for authentication
	Key        string //Key is the public key given by authentication provider
	ConfigInit bool   //ConfigInit defines Boolean parameter to initialise config
	Version    bool
	Debug      bool
}

// Application is an interface of the core application
type Application interface {
	parser()
	verify()
}

const configPath string = "/etc/nexusauth/"

//parser defines a function that creates a flag set from the package flag
func (app *App) parser() {
	flag.StringVar(&app.User, "u", "", "system user authenticating against")
	flag.StringVar(&app.Key, "k", "", "public key produced for authentication by user's agent")
	flag.BoolVar(&app.ConfigInit, "init", false, "initialise default configuration")
	flag.BoolVar(&app.Version, "version", false, "software version")
	flag.BoolVar(&app.Debug, "debug", false, "enable debug mode")
	flag.Parse()
}

//verify defines a validation function for commandline arguments provided at runtime
func (app *App) verify() {
	if app.User == "" {
		logging.Error(codes.CODE1)
	} else if app.Key == "" {
		logging.Error(codes.CODE2)
	}
	logging.Info("input verification successful")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()

	logging.LogFilePath = "/var/log/nexusauth.log"
	logging.Tag = "nexusauth"

	var (
		app    App
		config = configuration.System{
			Path:     configPath,
			Provider: providers.Provider{},
		}
		message notification.Message
	)
	app.parser()

	if app.Version {
		logging.Info("version printed")
		//fmt.Printf(SIGNATURE)
		fmt.Println("Version:", VERSION)
		os.Exit(0)
	}

	if app.ConfigInit {
		logging.Info("config initialised")
		if err := config.ConfigInit(); err != nil {
			logging.Error(err.Error())
		}
		os.Exit(0)
	}

	logging.Info("=============== starting authentication ===============")
	app.verify()

	logging.Info("loading config from system")
	if err := config.ReadConfig(); err != nil {
		logging.Error(err.Error())
	}

	message.Host = config.Host
	message.SystemUser = app.User
	if app.User == config.User || app.User == config.SudoUser {
		if users, err := config.Provider.Authenticate(app.User); err != nil {
			logging.Error(err.Error())
		} else {
			for _, user := range users {
				for _, key := range user.Keys {
					if key == app.Key {
						message.User = user.Username
						message.Provider = strings.Title(config.Provider.Name)
					}
					fmt.Println(key)
				}
			}
		}
	}

	if config.Notification != nil {
		if err := config.Notification.Notify(message); err != nil {
			logging.Warning(err.Error())
		}
	} else {
		logging.Warning(codes.CODE6)
	}

	logging.Info("=============== authentication process completed ===============")

}