package globals

import (
	"api/models"
	"github.com/kardianos/service"
)

//Defining all required global variables

// Logger is an application wide logger for windows...
var Logger service.Logger

// AppConfig is an Application config object
var AppConfig models.AppConfig
