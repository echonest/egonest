package egonest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGenericUnmarshal(t *testing.T) {
	resp := &http.Response{}
	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`
{"response": {"status": {"version": "4.2", "code": 0, "message": "Success"}, "artist": {"hotttnesss": 0.863645, "id": "ARH6W4X1187B99274F", "name": "Radiohead"}}}`)))
	m, err := GenericUnmarshal(resp, true)
	if err != nil {
		t.Fatal(err)
	}

	var ok bool
	var r map[string]interface{}
	if r, ok = m["response"].(map[string]interface{}); !ok {
		t.Fatal(m)
	}

	var s map[string]interface{}
	if s, ok = r["status"].(map[string]interface{}); !ok {
		t.Fatal(r)
	}

	if c, ok := s["code"].(json.Number); !ok {
		t.Fatal(s)
	} else if ci, err := c.Int64(); ci != 0 || err != nil {
		t.Fatal(ci, err)
	}

	var a map[string]interface{}
	if a, ok = r["artist"].(map[string]interface{}); !ok {
		t.Fatal(err)
	}

	if _, ok = a["id"].(string); !ok {
		t.Fatal(err)
	}

	if _, ok := a["name"].(string); !ok {
		t.Fatal(err)
	}

	if _, ok := a["hotttnesss"].(json.Number); !ok {
		t.Fatal(err)
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`
{"response": "status": {"version": "4.2", "code": 0, "message": "Success"}, "artist": {"hotttnesss": 0.863645, "id": "ARH6W4X1187B99274F", "name": "Radiohead"}}}`)))
	m, err = GenericUnmarshal(resp, true)
	if err == nil {
		t.Log("Should have failed")
		t.FailNow()
	}
	t.Log(err)
}

func TestCustomUnmarshal(t *testing.T) {
	resp := &http.Response{}
	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`
{"response": {"status": {"version": "4.2", "code": 0, "message": "Success"}, "artist": {"hotttnesss": 0.863645, "id": "ARH6W4X1187B99274F", "name": "Radiohead"}}}`)))
	var r struct {
		Response struct {
			Status `json:"status"`
			Artist struct {
				Id string `json:"id"`
			} `json:"artist"`
		} `json:"response"`
	}
	err := CustomUnmarshal(resp, &r)
	if err != nil {
		t.Log("JSON failure", err)
		t.Fail()
	}
	if r.Response.Artist.Id != "ARH6W4X1187B99274F" {
		t.Log("Wrong artist ID")
		t.Fail()
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`
{"response": {"status": {"version" "4.2", "code": 0, "message": "Success"}, "artist": {"hotttnesss": 0.863645, "id": "ARH6W4X1187B99274F", "name": "Radiohead"}}}`)))
	err = CustomUnmarshal(resp, &r)
	if err == nil {
		t.Log("Should have failed")
		t.FailNow()
	}
	t.Log(err)
}

func TestHTTPErrorUnmarshal(t *testing.T) {
	makeresp := func() *http.Response {
		resp := &http.Response{}
		resp.Body = ioutil.NopCloser(bytes.NewReader([]byte(`{"response": {"status": {"version": "4.2", "code": 1, "message": "1|Invalid key: Unknown"}}}`)))
		resp.StatusCode = 400
		resp.Status = "400 Bad Request"
		return resp
	}
	res, err := GenericUnmarshal(makeresp(), true)
	if err == nil {
		t.Log("should have failed:", res)
		t.Fail()
	}
	var r struct {
		Response struct {
			Status `json:"status"`
		} `json:"response"`
	}
	err = CustomUnmarshal(makeresp(), &r)
	if err == nil {
		t.Log("should have failed:", r)
		t.Fail()
	}
}

func TestDig(t *testing.T) {
	mtxt := `{"list":[0,1,2], "object":{"foo":"bar","bazzznesss":0.12312}}`
	var m map[string]interface{}
	err := json.Unmarshal([]byte(mtxt), &m)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	i := Dig(m, "list", 2)
	if i != float64(2) {
		t.Logf("Wrong value from map->list: %v != %v", i, float64(2))
		t.Fail()
	}
	i = Dig(m, "object", "foo")
	if i != "bar" {
		t.Logf("Wrong value from map->list: %v != \"bar\"", i)
		t.Fail()
	}
}
