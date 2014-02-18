package egonest

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func init() {
	var _ http.RoundTripper = panicker{}
}

type panicker struct{}

func (p panicker) RoundTrip(*http.Request) (*http.Response, error) {
	panic("It's all I do")
}

func TestClientPanic(t *testing.T) {
	var h Host
	h.Client.Transport = panicker{}
	resp, err := h.GetCall("doesnt/matter", url.Values{})
	t.Log(resp, err)
	if err == nil {
		t.Log("This should have returned an error.")
		t.Fail()
	}
}

func TestBadURL(t *testing.T) {
	var h Host
	h.Hostname = "!@#!@#!@#!@#!@#"
	checkerr := func(resp *http.Response, err error) {
		if err == nil {
			t.Log("This should have returned an error.")
			t.Fail()
		} else if _, ok := err.(*url.Error); !ok {
			t.Logf("unexpected error of type %T %v", err, err)
		}
	}
	resp, err := h.GetCall("doesnt matter", url.Values{})
	checkerr(resp, err)
	resp, err = h.PostCall("doesnt matter", url.Values{}, nil)
	checkerr(resp, err)
}

func TestBadResponses(t *testing.T) {
	var h Host
	checkerr := func(resp *http.Response, err error) {
		t.Log(resp, err)
		if err == nil {
			t.Log("didn't return error from nonexistent call!")
			t.Fail()
		}
		if e, ok := err.(ErrorStatus); !ok {
			t.Logf("Wrong error type: %T for %v", e, e)
			t.Fail()
		} else if *(e.HTTPError) != 404 {
			t.Log("Expected 404")
		}
	}
	resp, err := h.GetCall("artist/sprofile", url.Values{})
	checkerr(resp, err)
	resp, err = h.PostCall("artist/sprofile", url.Values{}, nil)
	checkerr(resp, err)
}

func TestGetCall(t *testing.T) {
	var h Host
	resp, err := h.GetCall("artist/profile", url.Values{"name": []string{"Radiohead"}, "bucket": []string{"hotttnesss"}})
	if err != nil {
		t.Fatal(err)
	}
	var r struct {
		Response struct {
			Status `json:"status"`
			Artist struct{ Id string } `json:"artist"`
		} `json:"response"`
	}
	err = CustomUnmarshal(resp, &r)
	if err != nil {
		t.Log("JSON failure", err)
		t.Fail()
	}
	if r.Response.Artist.Id != "ARH6W4X1187B99274F" {
		t.Log("Wrong artist ID", r.Response.Artist.Id)
		t.Fail()
	}
}

func TestPostCall(t *testing.T) {
	code := ReaderWrapper{"", bytes.NewReader([]byte(`{
  "metadata": {
    "artist": "Michael jackson",
    "release": "800 chansons des annes 80",
    "title": "Billie jean",
    "genre": "",
    "bitrate": 192,
    "sample_rate": 44100,
    "duration": 294,
    "filename": "../billie_jean.mp3",
    "samples_decoded": 220598,
    "given_duration": 20,
    "version": 3.13
  },
  "code": "eJxVlIuNwzAMQ1fxCDL133-xo1rnGqNAEcWy_ERa2aKeZmW9ustWVYrXrl5bthn_laFkzguNWpklEmoTB74JKYZSPlbJ0sy9fQrsrbEaO9W3bsbaWOoK7IhkHFaf_ag2d75oOQSZczbz5CKA7XgTIBIXASvFi0A3W8pMUZ7FZTWTVbujCcADlQ_f_WbdRNJ2vDUwSF0EZmFvAku_CVy440fgiIvArWZZWoJ7GWd-CVTYC5FCFI8GQdECdROE20UQfLoIUmhLC7IiByF1gzbAs3tsSKctyC76MPJlHRsZ5qhSQhu_CJFcKtW4EMrHSIrpTGLFqsdItj1H9JYHQYN7W2nkC6GDPjZTAzL9dx0fS4M1FoROHh9YhLHWdRchQSd_CLTpOHkQQP3xQsA2-sLOUD7CzxU0GmHVdIxh46Oide0NrNEmjghG44Ax_k2AoDHsiV6WsiD6OFm8y-0Lyt8haDBBzeMlAnTuuGYIB4WA2lEPAWbdeOabgFN6TQMs6ctLA5fHyKMBB0veGrjPfP00IAlWNm9n7hEh5PiYYBGKQDP-x4F0CL8HkhoQnRWN997JyEpnHFR7EhLPQMZmgXS68hsHktEVErranvSSR2VwfJhQCnkuwhBUcINNY-xu1pmw3PmBqU9-8xu0kiF1ngOa8vwBSSzzNw=="
}`))}

	var h Host
	resp, err := h.PostCall("song/identify", url.Values{}, map[string]UploadFile{"query": code})
	if err != nil {
		t.Fatal(err)
	}
	var r struct {
		Response struct {
			Status `json:"status"`
			Songs  []struct {
				Score       float64 `json:"score"`
				Title       string  `json:"title"`
				Message     string  `json:"message"`
				Artist_id   string  `json:"artist_id"`
				Artist_name string  `json:"artist_name"`
				Id          string  `json:"id"`
			} `json:"songs"`
		} `json:"response"`
	}
	err = CustomUnmarshal(resp, &r)
	if err != nil {
		t.Log("JSON failure", err)
		t.Fail()
	}
	if len(r.Response.Songs) < 1 {
		t.Log("WARNING: No song returned.", r)
		t.FailNow()
	}
	if r.Response.Songs[0].Artist_id != "ARXPPEY1187FB51DF4" {
		t.Log("Wrong artist ID", r.Response.Songs[0].Artist_id)
		t.Fail()
	}
	if r.Response.Songs[0].Id != "SODJXOA1313438FB61" {
		t.Log("Wrong Song ID", r.Response.Songs[0].Id)
		t.Fail()
	}
}

func TestRateLimit(t *testing.T) {
	var h Host
	h.Throttle = true // commenting this line should make the test fail
	h.SetDefaults()
	var headers = make(http.Header)
	headers.Set("X-RateLimit-Limit", "400")
	headers.Set("X-RateLimit-Used", "400")
	headers.Set("X-RateLimit-Remaining", "0")
	bucket := false
	testPart := func() {
		start := time.Now().Local()
		debugLogger.Println(start.Format(time.RFC850))
		headers.Set("Date", start.Format(time.RFC850))
		nextMinute := start
		timeToNextMinute := (60*time.Second - time.Duration(nextMinute.Second())*time.Second - time.Duration(nextMinute.Nanosecond())*time.Nanosecond)
		nextMinute = nextMinute.Add(timeToNextMinute)
		t.Log(start, timeToNextMinute, nextMinute)

		h.storeRateLimit("oogy/boogy", headers)
		if !bucket && h.CallToBucket("oogy/boogy") != "" {
			t.Log("Rate Limit Bucket set for non-bucket request")
			t.Fail()
		}
		cancel := time.AfterFunc(timeToNextMinute*2, func() {
			t.Fatal("Waited too long!")
		})
		h.delayIfNeeded("oogy/boogy")
		cancel.Stop()
		after := time.Now()
		t.Log("finished at", after)
		if after.Before(nextMinute) {
			t.Log("Didn't wait long enough!", nextMinute, after)
			t.Fail()
		}
	}
	t.Log("part 1")
	testPart()
	headers.Set("X-RateLimit-Bucket", "subbucket")
	bucket = true
	t.Log("part 2")
	testPart()
	if _, ok := h.rateLimits[""]; !ok {
		t.Log("Missing rate info for default bucket")
		t.Fail()
	}
	if _, ok := h.rateLimits["subbucket"]; !ok {
		t.Log("Missing rate info for sub bucket")
		t.Fail()
	}

	ratelimits := h.RateLimits()
	defaultbucket, ok := ratelimits[""]
	if !ok {
		t.Log("no default info stored")
		t.Fail()
	}
	subbucket, ok := ratelimits["subbucket"]
	if !ok {
		t.Log("no default info stored")
		t.Fail()
	}
	bucketcheck := RateLimitInfo{Bucket: "", Remaining: 0, Used: 400, Limit: 400, LastCall: defaultbucket.LastCall, Drift: defaultbucket.Drift}
	if defaultbucket != bucketcheck {
		t.Log("default bucket had wrong values", defaultbucket)
		t.Fail()
	}
	bucketcheck.Bucket = "subbucket"
	bucketcheck.LastCall = subbucket.LastCall
	bucketcheck.Drift = subbucket.Drift
	if subbucket != bucketcheck {
		t.Log("default bucket had wrong values", subbucket)
		t.Fail()
	}
}
