package egonest

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Example() {
	var h Host // instantiate the host. this will use the default hostname of 'developer.echonest.com'
	// and the API key from the environment

	args := make(url.Values)
	args.Set("filetype", "wav")

	file := ReaderWrapper{"noise.wav", bytes.NewReader(noise)}                          // noise is defined in a separate file
	resp, err := h.PostCall("track/upload", args, map[string]UploadFile{"track": file}) // see documentation
	// for track/upload: http://developer.echonest.com/docs/v4/track.html#upload

	// check for network error. A non-200 HTTP response is not reflected here.
	if err != nil {
		log.Println(err)
		return
	}

	// This struct mirrors the parts of the resposne to track/upload we need.
	var re struct {
		Response struct {
			Status `json:"status"`
			Track  struct {
				Status string
				Id     string
				Md5    string
			}
		} `json:"response"`
	}

	// CustomUnmarshal is intended to take a pointer to a struct, although it will also
	// accept a pointer to a map[string]interface{}.
	// This will close the response's Body.
	err = CustomUnmarshal(resp, &re)

	// check for non-200 HTTP status
	if err != nil {
		log.Println(err)
		return
	}

	// check for app-level error
	err = re.Response.Status.AsError()
	if err != nil {
		log.Println(err)
		return
	}

	var uploadstatus string
	var m map[string]interface{}
	args = make(url.Values)
	args.Set("id", re.Response.Track.Id)
	args.Set("bucket", "audio_summary")
	for uploadstatus = re.Response.Track.Status; uploadstatus == "pending"; {
		time.Sleep(2 * time.Second)
		resp, err := h.GetCall("track/profile", args)
		if err != nil {
			log.Println(err)
			return
		}

		// This will also close resp's Body.
		m, err = GenericUnmarshal(resp, false)
		if err != nil {
			log.Println(err)
			return
		}

		var ok bool

		// results from Dig still need type assertions
		uploadstatus, ok = Dig(m, "response", "track", "status").(string)
		if !ok {
			log.Printf("status string not found! %v", m)
			return
		}
	}

	analysis_url, ok := Dig(m, "response", "track", "audio_summary", "analysis_url").(string)
	if !ok {
		log.Println("analysis_url not found!")
		log.Println(m)
		return
	}

	resp, err = http.Get(analysis_url)
	if err != nil {
		log.Println(err)
		return
	}

	var full_analysis Analysis
	err = CustomUnmarshal(resp, &full_analysis)
	if err != nil {
		log.Println(analysis_url)
		log.Println(err)
		return
	}

	log.Println("track had", len(full_analysis.Segments), "segments")

	fmt.Println(uploadstatus)
	fmt.Println(Dig(m, "response", "track", "md5")) // Dig is a convenient function for pulling items out of a map[string]interface{} or []interface{}
	// Output:
	// complete
	// 2738fcc4359d716b40965fb406612ba0
}
