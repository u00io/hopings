package main

import (
	"github.com/u00io/gomisc/logger"
	"github.com/u00io/hopings/forms/mainform"
	"github.com/u00io/hopings/localstorage"
)

func main() {
	localstorage.Init("hopings")
	logger.Init(localstorage.Path() + "/logs")
	mainform.Run()
}
