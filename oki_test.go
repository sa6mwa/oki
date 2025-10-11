package oki

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

type jsonPayload struct {
	Message string
}

func TestDOKI_ErrorPassthrough(t *testing.T) {
	resp := struct{}{}
	wantErr := errors.New("boom")

	gotResp, gotErr := DOKI(resp, wantErr)
	if gotErr == nil || !errors.Is(gotErr, wantErr) {
		t.Fatalf("DOKI(_, boom) error = %v, want boom", gotErr)
	}
	if gotResp != resp {
		t.Fatalf("DOKI(resp, boom) response = %v, want original value", gotResp)
	}
}

func TestDOKI_InvalidOrMissingStatus(t *testing.T) {
	t.Run("missing status information", func(t *testing.T) {
		resp := &struct {
			JSON200 *jsonPayload
		}{}

		if _, err := DOKI(resp, nil); err == nil || !strings.Contains(err.Error(), "could not determine HTTP status") {
			t.Fatalf("DOKI(_, nil) error = %v, want status inference failure", err)
		}
	})

	t.Run("non-2xx status code", func(t *testing.T) {
		resp := &struct {
			HTTPResponse *http.Response
			JSON200      *jsonPayload
		}{
			HTTPResponse: &http.Response{StatusCode: 500},
			JSON200:      &jsonPayload{Message: "ignored"},
		}

		if _, err := DOKI(resp, nil); err == nil || err.Error() != "HTTP 500 Internal Server Error" {
			t.Fatalf("DOKI(_, nil) error = %v, want HTTP 500", err)
		}
	})
}

func TestDOKI_NoContentAllowed(t *testing.T) {
	resp := &struct {
		HTTPResponse *http.Response
	}{
		HTTPResponse: &http.Response{StatusCode: http.StatusNoContent},
	}

	if _, err := DOKI(resp, nil); err != nil {
		t.Fatalf("DOKI(_, nil) error = %v, want nil for 204", err)
	}
}

func TestDOKI_JSONValidation(t *testing.T) {
	t.Run("missing JSON field", func(t *testing.T) {
		resp := &struct {
			HTTPResponse *http.Response
			JSON200      *jsonPayload
		}{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200:      nil,
		}

		if _, err := DOKI(resp, nil); err == nil || !strings.Contains(err.Error(), "missing or empty JSON200") {
			t.Fatalf("DOKI(_, nil) error = %v, want missing JSON200", err)
		}
	})

	t.Run("valid JSON field", func(t *testing.T) {
		resp := &struct {
			HTTPResponse *http.Response
			JSON200      *jsonPayload
		}{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200:      &jsonPayload{Message: "ok"},
		}

		gotResp, err := DOKI(resp, nil)
		if err != nil {
			t.Fatalf("DOKI(_, nil) error = %v, want nil", err)
		}
		if gotResp != resp {
			t.Fatalf("DOKI(_, nil) response = %v, want original pointer", gotResp)
		}
	})
}

func TestDOKI_InferStatusFromJSON(t *testing.T) {
	resp := &struct {
		JSON201 *jsonPayload
	}{
		JSON201: &jsonPayload{Message: "created"},
	}

	if _, err := DOKI(resp, nil); err != nil {
		t.Fatalf("DOKI(_, nil) error = %v, want nil when status inferred from JSON201", err)
	}
}

func TestHTTPStatus(t *testing.T) {
	resp := struct {
		HTTPResponse *http.Response
	}{
		HTTPResponse: &http.Response{StatusCode: 201},
	}
	if got := httpStatus(resp); got != 201 {
		t.Fatalf("httpStatus(struct) = %d, want 201", got)
	}

	ptr := &resp
	if got := httpStatus(ptr); got != 201 {
		t.Fatalf("httpStatus(pointer) = %d, want 201", got)
	}

	resp.HTTPResponse = nil
	if got := httpStatus(resp); got != 0 {
		t.Fatalf("httpStatus(nil) = %d, want 0", got)
	}
}

func TestFindFirst2xxCode(t *testing.T) {
	resp := struct {
		JSON200 *jsonPayload
		JSON201 *jsonPayload
	}{
		JSON200: nil,
		JSON201: &jsonPayload{Message: "created"},
	}

	if got := findFirst2xxCode(resp); got != 201 {
		t.Fatalf("findFirst2xxCode = %d, want 201", got)
	}

	resp.JSON201 = nil
	if got := findFirst2xxCode(resp); got != 0 {
		t.Fatalf("findFirst2xxCode (all nil) = %d, want 0", got)
	}
}

func TestEnsureJSONForCode(t *testing.T) {
	value := jsonPayload{Message: "ok"}
	resp := struct {
		JSON200 *jsonPayload
		JSON201 string
	}{
		JSON200: &value,
	}

	if !ensureJSONForCode(resp, 200) {
		t.Fatalf("ensureJSONForCode(resp, 200) = false, want true")
	}
	if ensureJSONForCode(resp, 201) {
		t.Fatalf("ensureJSONForCode(resp, 201) = true, want false for zero string")
	}
	if ensureJSONForCode(resp, 202) {
		t.Fatalf("ensureJSONForCode(resp, 202) = true, want false for missing field")
	}
}
