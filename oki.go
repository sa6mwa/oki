// Provides generic function oki.DOKI wrapping WithResponses receiver functions
// by an oapi-codegen generated ClientWithResponses.
//
// (C) 2025 SA6MWA Michel https://pkt.systems
// License: MIT
package oki

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

var (
	ERR_NO_CONTENT error = fmt.Errorf("%d %s", http.StatusNoContent, http.StatusText(http.StatusNoContent))
)

var jsonFieldRe = regexp.MustCompile(`^JSON(\d{3})$`)

// DOKI: ensure 2xx AND that the corresponding JSON2xx field exists & is non-zero.
// Returns the original response so you can still access headers, etc. If there
// is only one JSON2xx field in resp and error is nil, it is guaranteed to be
// valid/non-nil.
func DOKI[R any](resp R, err error) (R, error) {
	if err != nil {
		return resp, err
	}
	code := httpStatus(resp)
	if code == 0 {
		// no HTTPResponse; try inferring from JSON fields
		if code = findFirst2xxCode(resp); code == 0 {
			return resp, fmt.Errorf("could not determine HTTP status for %T", resp)
		}
	}
	if code < 200 || code >= 300 {
		return resp, fmt.Errorf("HTTP %d %s", code, http.StatusText(code))
	}
	// For 204 No Content: allow empty body (no JSON204 field expected)
	if code == http.StatusNoContent {
		return resp, nil
	}
	// For other 2xx: ensure matching JSON<code> field exists and is non-zero.
	if !ensureJSONForCode(resp, code) {
		return resp, fmt.Errorf("missing or empty JSON%d on %T", code, resp)
	}
	return resp, nil
}

// Returns http.Response StatusCode if HTTPResponse field in resp exist and has
// valid values, otherwise 0.
func httpStatus[R any](resp R) int {
	rv := reflect.ValueOf(resp)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return 0
	}
	f := rv.FieldByName("HTTPResponse")
	if f.IsValid() && !f.IsNil() {
		if hr, ok := f.Interface().(*http.Response); ok && hr != nil {
			return hr.StatusCode
		}
	}
	return 0
}

// findFirst2xxCode finds a JSON2xx field in resp and returns it's 2xx number as
// int. This does not mean that the status code of the http response was 2xx,
// you need to check that later with ensureJSONForCode with the return value by
// this function (unless 0).
func findFirst2xxCode[R any](resp R) int {
	rv := reflect.ValueOf(resp)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return 0
	}
	rt := rv.Type()
	for i := range rt.NumField() {
		name := rt.Field(i).Name
		m := jsonFieldRe.FindStringSubmatch(name)
		if len(m) != 2 {
			continue
		}
		code, _ := strconv.Atoi(m[1])
		if code >= 200 && code < 300 {
			if !rv.Field(i).IsNil() {
				return code
			}
		}
	}
	return 0
}

// ensureJSONForCode ensures the JSON{code} field in resp exists and is
// non-zero/non-nil.
func ensureJSONForCode[R any](resp R, code int) bool {
	rv := reflect.ValueOf(resp)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return false
	}
	field := rv.FieldByName(fmt.Sprintf("JSON%d", code))
	if !field.IsValid() {
		return false
	}
	// Non-zero means: non-nil pointer/map/slice/interface and non-zero value otherwise.
	return !field.IsZero()
}
