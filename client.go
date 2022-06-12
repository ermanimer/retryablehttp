package retryablehttp

import (
	"errors"
	"net/http"
	"time"
)

// errors
var (
	ErrNilHTTPClient          = errors.New("http client is nil")
	ErrInvalidMaxReqCount     = errors.New("maximum request count is not valid")
	ErrInvalidBackoff         = errors.New("backoff is not valid")
	ErrNilResHandler          = errors.New("response handler is nil")
	ErrNilRes                 = errors.New("response is nil")
	ErrUnsuccessfulStatusCode = errors.New("unsuccessful status code")
)

// default options
const (
	defaultMaxReqCount = 1
	defaultBackoff     = 0
)

var (
	defaultResHandler = func(res *http.Response) error {
		if res == nil {
			return ErrNilRes
		}

		statusCode := res.StatusCode
		if statusCode < 200 || statusCode > 299 {
			return ErrUnsuccessfulStatusCode
		}

		return nil
	}
)

// Client represents retryable http client.
type Client struct {
	httpClient  *http.Client
	maxReqCount int
	backoff     time.Duration
	resHandler  func(res *http.Response) error
}

// Option configures client options.
type Option func(c *Client) error

// WithHTTPClient configures client's http client.
// Default http client is http.DefaultClient{}.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return ErrNilHTTPClient
		}

		c.httpClient = httpClient

		return nil
	}
}

// WithMaxReqCount configures client's max request count.
// Default maximum request count is 1.
func WithMaxReqCount(maxReqCount int) Option {
	return func(c *Client) error {
		if maxReqCount < 1 {
			return ErrInvalidMaxReqCount
		}

		c.maxReqCount = maxReqCount

		return nil
	}
}

// WithBackoff configures client's backoff duration, which represents sleeping intervals between retries.
// Default backoff duration is 0.
func WithBackoff(backoff time.Duration) Option {
	return func(c *Client) error {
		if backoff < 0 {
			return ErrInvalidBackoff
		}

		c.backoff = backoff

		return nil
	}
}

// WithResHandler configures client's response handler function which handles http response.
// Default response handler:
//
//  func defaultResHandler(res *http.Response) error {
//  	if res == nil {
//  		return ErrNilRes
//  	}
//
//  	statusCode := res.StatusCode
//  	if statusCode < 200 || statusCode > 299 {
//  		return ErrUnsuccessfulStatusCode
//  	}
//
//  	return nil
//  }
func WithResHandler(resHandler func(res *http.Response) error) Option {
	return func(c *Client) error {
		if resHandler == nil {
			return ErrNilResHandler
		}

		c.resHandler = resHandler

		return nil
	}
}

// NewClient creates and returns new retryable http client instance.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		httpClient:  http.DefaultClient,
		maxReqCount: defaultMaxReqCount,
		backoff:     defaultBackoff,
		resHandler:  defaultResHandler,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Do sends http request with automatic retries returns first successful or last unsuccessful response.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error
	for i := 0; i < c.maxReqCount; i++ {
		res, err = c.httpClient.Do(req)

		if err == nil {
			err = c.resHandler(res)
		}

		if err != nil {
			time.Sleep(c.backoff)

			continue
		}

		break
	}

	return res, err
}
