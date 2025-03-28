package main

import (
	"github.com/reiver/barnameh-api/cfg"
	"github.com/reiver/barnameh-api/srv/log"
)

func reveal() {
	log := logsrv.Prefix("reveal").Begin()
	defer log.End()

	var tcpaddr string = cfg.WebServerTCPAddress()
	log.Informf("serving HTTP on TCP address: %q", tcpaddr)
}
