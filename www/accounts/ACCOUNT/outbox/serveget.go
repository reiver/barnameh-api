package verboten

import (
	"io"
	"net/http"

	"github.com/reiver/go-fediverseid"
	"github.com/reiver/go-http400"
	"github.com/reiver/go-http500"

	"github.com/reiver/barnameh-api/srv/log"
)

func serveGET(responsewriter http.ResponseWriter, request *http.Request, account string) {
	log := logsrv.Prefix("www("+path+").serveget").Begin()
	defer log.End()

	if nil == responsewriter {
		log.Error("nil response-writer")
		return
	}
	if nil == request {
		http500.Serve(responsewriter)
		log.Error("nil request")
		return
	}

	{
		fid, err := fediverseid.ParseFediverseIDString(account)
		if nil != err {
			http400.Serve(responsewriter)
			log.Debug("bad fediverse-id %q: %s", account, err)
			return
		}
		log.Debugf("fediverse-id: %q", fid)
	}

			
//@TODO: implement the actual outbox
			

	{

		const content string = `{"@context": "https://www.w3.org/ns/activitystreams"}`
		io.WriteString(responsewriter, content)
	}
}
