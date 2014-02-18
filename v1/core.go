// Copyright 2013 The Echo Nest. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
egonest is a wrapper for the Echo Nest API for the Go programming language.
It provides a basic interface for making calls to the API, using Go's built-in json decoding
to make it simpler to deal with responses.



*/

package egonest

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type RateLimitInfo struct {
	Bucket                 string
	Limit, Used, Remaining int
	LastCall               time.Time
}

// A Host contains the most basic information necessary for communicating with The Echo Nest; the API hostname and key.
// The zero value of a Host is usable; the default hostname is "developer.echonest.com" and the default API key is the value of the environment variable ECHO_NEST_API_KEY.
// A Host is safe to use from multiple goroutines except for changes to its exported fields.

// Hostname, if left unset, will be set to "developer.echonest.com" on first use.

// BasePath, if left unset, will be set to "/api/v4/" on first use.

// ApiKey, if left unset, will be set to the value of the environment variable ECHO_NEST_API_KEY on first use.

// Client may be altered or replaced as needed to suit your environment's needs.

// Throttle, if true, will use rate-limiting information from the Echo Nest API to delay requests that would otherwise cause an error due to exceeding the API key's rate limit.

// A method call against a Host will result in at most one call against the API unless otherwise noted, and will not panic unless otherwise noted.
type Host struct {
	Hostname, BasePath, ApiKey string
	Client                     http.Client
	Throttle                   bool
	callToBucket               map[string]string
	rateLimits                 map[string]RateLimitInfo
	rateLimitLock              sync.RWMutex
	set                        sync.Once
}

const DefaultHost = "developer.echonest.com"
const DefaultBasePath = "/api/v4/"

var defaultApiKey string

var debugLogger = log.New(ioutil.Discard, "", 0)

func init() {
	defaultApiKey = os.Getenv("ECHO_NEST_API_KEY")
	if os.Getenv("EGONEST_DEBUG") != "" {
		debugLogger = log.New(os.Stderr, "egonest-", log.Lshortfile|log.LstdFlags)
	}
}

// SetDefaults sets a zero-valued Host structure's internal fields to the default values.
// This is called by GetCall and PostCall if needed during the first request and will only
// take effect once.
func (h *Host) SetDefaults() {
	h.set.Do(func() {
		if h.Hostname == "" {
			h.Hostname = DefaultHost
		}
		if h.BasePath == "" {
			h.BasePath = DefaultBasePath
		}
		if h.ApiKey == "" {
			h.ApiKey = defaultApiKey
		}
		h.rateLimits = make(map[string]RateLimitInfo)
		h.callToBucket = make(map[string]string)
	})
}

func copyValues(in url.Values) (out url.Values) {
	out = make(url.Values)
	for k, v := range in {
		out[k] = v // not copying individual slices as we don't mutate these
	}
	return
}

// GetCall will obtain the raw response for an Echo Nest API call made through the GET method.
// The caller must call resp.Body.Close() (directly or through GenericUnmarshal or CustomUnmarshal) if err is not nil.
// Calling this function will make a single API request.

func (h *Host) GetCall(call string, args url.Values) (resp *http.Response, err error) {
	defer func() {
		if r := recover(); r != nil {
			if resp != nil {
				resp.Body.Close()
			}
			err = convertToErr(r)
		}
	}()

	h.SetDefaults()

	args = copyValues(args)

	args.Set("api_key", h.ApiKey)
	args.Set("format", "json")
	u := &url.URL{Scheme: "http", Host: h.Hostname, Path: path.Join(h.BasePath, call), RawQuery: args.Encode()}
	debugLogger.Println(u)
	req, reqerr := http.NewRequest("GET", u.String(), nil)
	if reqerr != nil {
		err = reqerr
		return
	}
	// if there's a need for another GET and POST header here, refactor this to a new function
	req.Header.Add("User-Agent", userAgent)
	h.delayIfNeeded(call)
	resp, httperr := h.Client.Do(req)
	debugLogger.Println(u.String())
	if httperr != nil {
		err = httperr
		return
	}
	h.storeRateLimit(call, resp.Header)
	if resp.StatusCode != http.StatusOK {
		code := resp.StatusCode
		err = ErrorStatus{Status: nil, HTTPError: &code}
	} else {
		h.storeRateLimit(call, resp.Header)
	}
	return resp, err
}

// PostCall will obtain the raw response for an Echo Nest API call made through the POST method.

// The caller must call resp.Body.Close() (directly or through GenericUnmarshal or CustomUnmarshal) if err is not nil. If the io.Readers supplied need to be closed, the caller is responsible for that too.

// Calling this function will make a single API request.

func (h *Host) PostCall(call string, args url.Values, files map[string]UploadFile) (resp *http.Response, err error) {
	defer func() {
		if r := recover(); r != nil {
			if resp != nil {
				resp.Body.Close()
			}
			err = convertToErr(r)
		}
	}()

	h.SetDefaults()

	args = copyValues(args)

	args.Set("api_key", h.ApiKey)
	args.Set("format", "json")

	u := &url.URL{Scheme: "http", Host: h.Hostname, Path: path.Join(h.BasePath, call)}
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	go func() {
		var err error
		defer func() {
			// if no errors have occurred, err is nil here
			pw.CloseWithError(err)
		}()
		if args != nil {
			for k, vals := range args {
				for _, val := range vals {
					err = mw.WriteField(k, val)
					if err != nil {
						mw.Close() // this would return same err as WriteField()
						return
					}
				}
			}
		}
		if files != nil {
			for name, reader := range files {
				var fw io.Writer
				fw, err = mw.CreateFormFile(name, filepath.Base(reader.Name()))
				if err != nil {
					mw.Close()
					return
				}
				_, err = io.Copy(fw, reader)
				if err != nil {
					mw.Close()
					return
				}
			}
		}

		err = mw.Close()
		// pipewriter close is deferred!
	}()

	req, reqerr := http.NewRequest("POST", u.String(), pr)
	if reqerr != nil {
		err = reqerr
		return
	}

	// if there's a need for another header for both GET and POST here, refactor this to a new function
	req.Header.Add("Content-Type", mw.FormDataContentType())
	req.Header.Add("User-Agent", userAgent)
	h.delayIfNeeded(call)
	resp, httperr := h.Client.Do(req)

	if httperr != nil {
		err = httperr
		return
	}
	err = pr.Close()
	if err != nil {
		resp.Body.Close()
		resp = nil
	}
	if resp.StatusCode != http.StatusOK {
		code := resp.StatusCode
		err = ErrorStatus{Status: nil, HTTPError: &code}
	} else {
		h.storeRateLimit(call, resp.Header)
	}
	return resp, err
}

func (h *Host) storeRateLimit(call string, headers http.Header) {
	// check for Bucket header
	var err error
	l, err := strconv.ParseInt(headers.Get("x-ratelimit-limit"), 10, 32)
	if err != nil {
		debugLogger.Println("Error parsing rate limit:", err)
	}
	u, err := strconv.ParseInt(headers.Get("x-ratelimit-used"), 10, 32)
	if err != nil {
		debugLogger.Println("Error parsing rate limit used:", err)
	}
	r, err := strconv.ParseInt(headers.Get("x-ratelimit-remaining"), 10, 32)
	if err != nil {
		debugLogger.Println("Error parsing rate limit remaining:", err)
	}
	lc, err := http.ParseTime(headers.Get("Date"))
	if err != nil {
		debugLogger.Println("Error parsing date header:", err)
	}
	info := RateLimitInfo{Bucket: headers.Get("x-ratelimit-bucket"),
		Used: int(u), Limit: int(l), Remaining: int(r), LastCall: lc}

	h.rateLimitLock.Lock()
	defer h.rateLimitLock.Unlock()

	if h.rateLimits[info.Bucket].LastCall.Before(info.LastCall) {
		h.callToBucket[call] = info.Bucket
		h.rateLimits[info.Bucket] = info
	}

}

// Returns a map of bucket name to last retrieved rate limit information since program start.
func (h *Host) RateLimits() map[string]RateLimitInfo {
	h.SetDefaults()
	h.rateLimitLock.RLock()
	defer h.rateLimitLock.RUnlock()

	result := make(map[string]RateLimitInfo, len(h.rateLimits))
	for k, v := range h.rateLimits {
		result[k] = v
	}
	return result
}

// Returns the rate limit bucket for call or "" if there is none.
// If call has not been used since the instantiation of the Host struct then the result will be "".
func (h *Host) CallToBucket(call string) string {
	h.SetDefaults()
	h.rateLimitLock.RLock()
	defer h.rateLimitLock.RUnlock()
	return h.callToBucket[call]
}

func (h *Host) delayIfNeeded(call string) {
	if !h.Throttle {
		return
	}
	h.rateLimitLock.RLock()
	defer h.rateLimitLock.RUnlock()
	bucket, ok := h.callToBucket[call]
	limit, ok := h.rateLimits[bucket]
	if !ok {
		return
	}
	if limit.Remaining > 0 {
		return
	}

	// seconds that were left until the next rate limit reset at the time
	// of the last call
	secondsleft := time.Duration(60-limit.LastCall.Second()) * time.Second

	if time.Now().Sub(limit.LastCall) > secondsleft {
		// sleep until top of the minute
		sleept := time.Duration(60-time.Now().Second()) * time.Second
		debugLogger.Println("sleeping for", sleept)
		time.Sleep(sleept)
	}
}

func convertToErr(r interface{}) (err error) {
	switch e := r.(type) {
	case error:
		err = e
	case string:
		err = errors.New(e)
	case fmt.Stringer:
		err = fmt.Errorf("%s", e.String())
	case fmt.GoStringer:
		err = fmt.Errorf("%#v", e)
	default:
		err = fmt.Errorf("%v", e)
	}
	return
}

type Status struct {
	Version string
	Code    int
	Message string
}

func (s *Status) AsError() error {
	if s.Code != 0 {
		return ErrorStatus{Status: s}
	}
	return nil
}

// ErrorStatus will be returned by function calls when the error is above the HTTP transport layer.
// For example: rate-limited API calls, invalid arguments or API keys, HTTP 4xx or 5xx errors.
// Errors reading or writing the response itself may be of type http.ProtocolError, any error type from the
// net package, or any other error type that may be returned by the I/O processes used in forming your request.
type ErrorStatus struct {
	*Status
	HTTPError *int
}

func (e ErrorStatus) Error() string {
	if e.Status != nil {
		return e.Message
	}
	if e.HTTPError != nil {
		return fmt.Sprintf("%d %s", *(e.HTTPError), http.StatusText(*(e.HTTPError)))
	}
	return "Unknown error"
}
