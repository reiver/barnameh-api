package main

import (
	"github.com/reiver/barnameh-api/srv/log"
)

func main() {
	log := logsrv.Prefix("main").Begin()
	defer log.End()

	log.Inform("barnameh ⚡")
	shout()

	log.Inform("Let me show you something…")
	reveal()

	log.Inform("Here we go…")
	webserve()
}
