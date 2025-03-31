package verboten

import (
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/reiver/go-fediverseid"
	"github.com/reiver/go-http201"
	"github.com/reiver/go-http400"
	"github.com/reiver/go-http415"
	"github.com/reiver/go-http500"
	"github.com/reiver/go-json"

	"github.com/reiver/barnameh-api/srv/log"
)

func servePOST(responsewriter http.ResponseWriter, request *http.Request, account string) {
	log := logsrv.Prefix("www("+path+").servepost").Begin()
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

	var contentType string = request.Header.Get("Content-Type")
	{
		const contentTypeActivity string = `application/activity+json`
		const contentTypeJSONLD   string = `application/ld+json; profile="https://www.w3.org/ns/activitystreams"`

		switch contentType {
		case contentTypeActivity: // application/activity+json
			// good
		case contentTypeJSONLD:   // application/ld+json; profile="https://www.w3.org/ns/activitystreams"
			// good
			
//@TODO: handle variation in those media-types better. ex: extra spacing
			
		default:
			http415.Serve(responsewriter, contentTypeActivity, contentTypeJSONLD)
			log.Debugf("unsupported media-type: %s", contentType)
			return
		}
	}
	log.Debugf("HTTP POST Content-Type: %s", contentType)

	var content []byte
	{
		var body io.ReadCloser = request.Body
		if nil != body {
			var err error
			content, err = io.ReadAll(body)
			if nil != err {
				http500.Serve(responsewriter)
				log.Errorf("problem reading HTTP POST body: %s", err)
				return
			}
		}

		if len(content) <= 0 {
			http400.Serve(responsewriter)
			log.Debug("empty content")
			return
		}
	}

	{
		var data map[string]any = map[string]any{}

		err := json.Unmarshal(content, &data)
		if nil != err {
			http400.Serve(responsewriter)
			log.Debugf("problem parsing content as JSON: %s", err)
			return
		}

		{
			const name string = "type"

			value, found := data[name]

			if !found {
				http400.Serve(responsewriter)
				log.Debugf("missing JSON-LD ActivityStreams %q field", name)
				return
			}
			if reflect.ValueOf(value).IsZero() {
				http400.Serve(responsewriter)
				log.Debugf("empty JSON-LD ActivityStreams %q field", name)
				return
			}
		}

		{
			var names []string = []string{"to", "bto", "audience", "cc", "bcc"}

			var isAddressed bool
			for _, name := range names {
				value, found := data[name]

				if !found {
					continue
				}
				if reflect.ValueOf(value).IsZero() {
					continue
				}

				isAddressed = true
				break
			}

			if !isAddressed {
				http400.Serve(responsewriter)
				log.Debugf("missing a non-empty addressing field — %s", strings.Join(names, ", "))
				return
			}
		}

	}

			
//@TODO: set a real location, so the user can edit, delete, etc.
			
	var location string = "/todo"

	http201.ServeLocation(responsewriter, location)
}
