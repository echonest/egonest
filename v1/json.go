package egonest

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Unmarshals the JSON in a request to a map[string]interface{}. Values may be of type string,
// []interface{}, map[string]interface{}, or float64 (or json.Number if useNumber is true)

// GenericUnmarshal will call resp.Body.Close(). To retrieve rate limit information, you may still inspect
// the headers after the body is closed.

// If err is not nil, the result may still have useful information about the failure.
func GenericUnmarshal(resp *http.Response, useNumber bool) (result map[string]interface{}, err error) {
	defer resp.Body.Close()
	j := json.NewDecoder(resp.Body)
	if useNumber {
		j.UseNumber()
	}

	err = j.Decode(&result)
	if err != nil {
		return
	}
	if resp.StatusCode >= http.StatusBadRequest {
		err = errors.New(resp.Status)
	}
	return
}

// Dig takes a series of arguments for drilling into an unmarshalled map[string]interface{} or []interface{}. Returns nil if the object is not found.
func Dig(subject interface{}, args ...interface{}) (result interface{}) {
	debugLogger.Printf("subject: %v %T args: %v", subject, subject, args)
	switch s := subject.(type) {
	case map[string]interface{}:
		a := args[0]
		st, ok := a.(string)
		if !ok {
			debugLogger.Println("not found")
			return nil
		}
		result = s[st]
		if len(args) == 1 {
			return
		}
	case []interface{}:
		a := args[0]
		in, ok := a.(int)
		if !ok {
			debugLogger.Println("not found")
			return nil
		}
		if in >= len(s) {
			debugLogger.Println("not found")
			return nil
		}
		result = s[in]
		if len(args) == 1 {
			return
		}
	default:
		debugLogger.Println("not found")
		return nil
	}
	return Dig(result, args[1:]...)

}

// CustomUnmarshal will unmarshal the JSON in resp's Body into dest. dest must be a pointer to a struct with
// a single struct field named Response or tagged json:"response".
// CustomUnmarshal will only return a non-nil error if resp's HTTP status code is not 200.
func CustomUnmarshal(resp *http.Response, dest interface{}) (err error) {
	defer resp.Body.Close()
	j := json.NewDecoder(resp.Body)
	err = j.Decode(dest)
	if err == nil {
		if resp.StatusCode >= http.StatusBadRequest {
			err = errors.New(resp.Status)
		}
	}
	return err
}
