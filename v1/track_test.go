package egonest

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"time"
)

func ExampleTrackUploadLocalFile() {
	var h Host

	args := make(url.Values)
	args.Set("filetype", "wav")
	file := ReaderWrapper{"noise.wav", bytes.NewReader(noise)}
	resp, err := h.PostCall("track/upload", args, map[string]UploadFile{"track": file})

	// check network error
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	var re struct {
		TrackUploadResult `json:"response"`
	}

	// check for json error or HTTP error
	err = CustomUnmarshal(resp, &re)
	if err != nil {
		log.Println(err)
		return
	}

	// check for app-level error
	err = re.Status.AsError()
	if err != nil {
		log.Println(err)
		return
	}

	var uploadstatus string
	var m map[string]interface{}
	args = make(url.Values)
	args.Set("id", re.Track.Id)
	for uploadstatus = re.Track.Status; uploadstatus == "pending"; {
		time.Sleep(2 * time.Second)
		resp, err := h.GetCall("track/profile", args)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
		m, err = GenericUnmarshal(resp, false)
		if err != nil {
			log.Println(err)
			return
		}
		var ok bool
		uploadstatus, ok = Dig(m, "response", "track", "status").(string)
		if !ok {
			log.Printf("status string not found! %v", m)
			return
		}
	}
	fmt.Println(uploadstatus)
	fmt.Println(Dig(m, "response", "track", "md5"))
	// Output:
	// complete
	// 2738fcc4359d716b40965fb406612ba0
}

type TrackUploadResult struct {
	Status `json:"status"`
	Track  struct {
		Status string
		Id     string
		Md5    string
	}
}
