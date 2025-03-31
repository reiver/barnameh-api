package verboten

import (
	"net/http"

	"github.com/reiver/go-erorr"
	"github.com/reiver/go-http500"

	"github.com/reiver/barnameh-api/srv/http"
	"github.com/reiver/barnameh-api/srv/log"
)

const path string = "/accounts/{account}/outbox"

func init() {
	log := logsrv.Prefix("www("+path+").init").Begin()
	defer log.End()

	var handler httpsrv.PatternHandler = httpsrv.PatternHandlerFunc(serveHTTP)

	err := httpsrv.Mux.HandlePattern(handler, path)
	if nil != err {
		e := erorr.Errorf("problem registering http-handler with path-mux for path %q: %w", path, err)
		log.Error(e)
		panic(e)
	}
}

func serveHTTP(responsewriter http.ResponseWriter, request *httpsrv.ParameterizedRequest) {
	log := logsrv.Prefix("www("+path+").servehttp").Begin()
	defer log.End()

	if nil == responsewriter {
		log.Error("nil response-writer")
		return
	}
	if nil == request {
		http500.Serve(responsewriter)
		log.Error("nil parameterized-request")
		return
	}

	var account string
	{
		var found bool
		account, found = request.ParameterByIndex(0)
		if !found {
			http500.Serve(responsewriter)
			log.Error("could not get account")
			return
		}
	}

	var httpRequest *http.Request
	{
		httpRequest = request.HTTPRequest()
		if nil == httpRequest {
			http500.Serve(responsewriter)
			log.Error("nil http-request")
			return
		}
	}

	log.Debugf("http-method: %q", httpRequest.Method)
	switch httpRequest.Method {
	case http.MethodGet:
		serveGET(responsewriter, httpRequest, account)
		return
	case http.MethodPost:
		servePOST(responsewriter, httpRequest, account)
		return
	default:
		const code int = http.StatusMethodNotAllowed
		http.Error(responsewriter, http.StatusText(code), code)
		log.Errorf("method not supported: %q", httpRequest.Method)
		return
	}
}
