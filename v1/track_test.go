package egonest

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Example() {
	// let's make some audio! (warning: turn volume down before listening)
	var noise = []byte{0x52, 0x49, 0x46, 0x46, 0x1c, 0x30, 0x14, 0x0, 0x57, 0x41, 0x56, 0x45, 0x66, 0x6d, 0x74, 0x20, 0x10, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x22, 0x56, 0x0, 0x0, 0x44, 0xac, 0x0, 0x0, 0x2, 0x0, 0x10, 0x0, 0x64, 0x61, 0x74, 0x61, 0xf8, 0x2f, 0x14, 0x0}
	var noiselen uint32
	err := binary.Read(bytes.NewReader(noise[len(noise)-4:]), binary.LittleEndian, &noiselen)
	if err != nil {
		log.Println(err)
		return
	}

	noiselen /= 2
	var noisedata = new(bytes.Buffer)
	rand.Seed(time.Now().Unix() / 86400)
	type sine struct {
		start, end int
		freq, amp  float64
	}
	sines := make([]sine, int(rand.Intn(27)+5))
	for s := range sines {
		sines[s].start = rand.Intn(int(noiselen))
		sines[s].end = sines[s].start + rand.Intn(int(noiselen)-sines[s].start)
		sines[s].freq = rand.Float64()*(11025-20) + 20
		sines[s].amp = rand.Float64()
		log.Println(s, sines[s].start, sines[s].end, sines[s].freq, sines[s].amp)
	}
	for i := 0; i < int(noiselen); i++ {
		var val float64
		for _, s := range sines {
			if i >= s.start && i <= s.end {
				val += 32767 * s.amp * math.Sin(s.freq*float64(i))
			}
		}
		ov := int16(val)
		err = binary.Write(noisedata, binary.LittleEndian, &ov)
	}
	noise = append(noise, noisedata.Bytes()...)

	/* uncomment this to write the generated audio to a file. heed the warning above.
	f, err := os.Create("noise.wav")
	if err != nil {
		log.Print(err)
		return
	}
	_, err = io.Copy(f, bytes.NewReader(noise))
	if err != nil {
		log.Print(err)
		return
	}
	f.Close()
	*/

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
	log.Println(Dig(m, "response", "track", "md5")) // Dig is a convenient function for pulling items out of a map[string]interface{} or []interface{}
	// Output:
	// complete
}
