package retryablehttp

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// NewClient function should return ErrNilHTTPClient when nil http client is provided.
func TestNilHTTPClientOption(t *testing.T) {
	_, err := NewClient(
		WithHTTPClient(nil),
	)
	if err == nil {
		t.Error("unexpected nil error")
	}
	if err != ErrNilHTTPClient {
		t.Errorf("unexpected error, %s", err)
	}
}

// NewClient function should return ErrInvalidMaxReqCount when zero or negative maximum request count is provided.
func TestInvalidMaxReqCountOption(t *testing.T) {
	_, err := NewClient(
		WithMaxReqCount(0),
	)
	if err == nil {
		t.Error("unexpected nil error")
	}
	if err != ErrInvalidMaxReqCount {
		t.Errorf("unexpected error, %s", err)
	}
}

// NewClient function should return ErrInvalidBackoff when negative backoff duration is provided.
func TestInvalidBackoffOption(t *testing.T) {
	_, err := NewClient(
		WithBackoff(-1),
	)
	if err == nil {
		t.Error("unexpected nil error")
	}
	if err != ErrInvalidBackoff {
		t.Errorf("unexpected error, %s", err)
	}
}

// NewClient function should return ErrNilResHandler when nil response handler is provided.
func TestNilResponseHandler(t *testing.T) {
	_, err := NewClient(
		WithResHandler(nil),
	)
	if err == nil {
		t.Error("unexpected nil error")
	}
	if err != ErrNilResHandler {
		t.Errorf("unexpected error, %s", err)
	}
}

// Do method of a client with default options should return nil error and status code ok when test server always returns status code ok.
func TestWithDefaultOptionsAndStatusOK(t *testing.T) {
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s := httptest.NewServer(m)
	defer s.Close()

	c, err := NewClient()
	if err != nil {
		t.Errorf("creating client failed, %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, s.URL, http.NoBody)
	if err != nil {
		t.Errorf("creating http request failed, %s", err.Error())
	}

	res, err := c.Do(req)
	if err != nil {
		t.Errorf("doing http request failed, %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code, %d", res.StatusCode)
	}
}

// Do method of a client with default options should return error and status code bad request when test server always returns status code bad request.
func TestWithDefaultOptionsAndStatusBadRequest(t *testing.T) {
	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	s := httptest.NewServer(m)
	defer s.Close()

	c, err := NewClient()
	if err != nil {
		t.Errorf("creating client failed, %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, s.URL, http.NoBody)
	if err != nil {
		t.Errorf("creating http request failed, %s", err.Error())
	}

	res, err := c.Do(req)
	if err == nil {
		t.Error("unexpected nil error")
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code, %d", res.StatusCode)
	}
}

// Do method of a client should return nil error and status code ok when test server returns status code ok at last request. And request duration should be longer than sum of backoff durations.
func TestWithOptionsAndStatusOK(t *testing.T) {
	httpClient := http.DefaultClient
	maxReqCount := 3
	backoff := 100 * time.Millisecond
	responseHandler := func(res *http.Response) error {
		if res.StatusCode != http.StatusOK {
			return ErrUnsuccessfulStatusCode
		}

		return nil
	}

	m := http.NewServeMux()

	reqCount := 0
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		if reqCount < maxReqCount {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		w.WriteHeader(http.StatusOK)
	})

	s := httptest.NewServer(m)
	defer s.Close()

	c, err := NewClient(
		WithHTTPClient(httpClient),
		WithMaxReqCount(maxReqCount),
		WithBackoff(backoff),
		WithResHandler(responseHandler),
	)
	if err != nil {
		t.Errorf("creating client failed, %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, s.URL, http.NoBody)
	if err != nil {
		t.Errorf("creating http request failed, %s", err.Error())
	}

	beginning := time.Now()

	res, err := c.Do(req)
	if err != nil {
		t.Errorf("doing http request failed, %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code, %d", res.StatusCode)
	}

	duration := time.Since(beginning)
	if duration < time.Duration(maxReqCount-1)*backoff {
		t.Error("unexpected duration")
	}
}

// Do method of a client should return error and status code bad request when test server always returns status code bad request. And request duration should be longer than sum of backoff durations.
func TestWithOptionsAndStatusBadRequest(t *testing.T) {
	httpClient := http.DefaultClient
	maxReqCount := 3
	backoff := 100 * time.Millisecond
	responseHandler := func(res *http.Response) error {
		if res.StatusCode != http.StatusOK {
			return ErrUnsuccessfulStatusCode
		}

		return nil
	}

	m := http.NewServeMux()

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	s := httptest.NewServer(m)
	defer s.Close()

	c, err := NewClient(
		WithHTTPClient(httpClient),
		WithMaxReqCount(maxReqCount),
		WithBackoff(backoff),
		WithResHandler(responseHandler),
	)
	if err != nil {
		t.Errorf("creating client failed, %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, s.URL, http.NoBody)
	if err != nil {
		t.Errorf("creating http request failed, %s", err.Error())
	}

	beginning := time.Now()

	res, err := c.Do(req)
	if err == nil {
		t.Error("unexpected nil error")
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code, %d", res.StatusCode)
	}

	duration := time.Since(beginning)
	if duration < time.Duration(maxReqCount-1)*backoff {
		t.Error("unexpected duration")
	}
}
