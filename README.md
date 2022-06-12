# retryablehttp
Retryable HTTP Client In Go

[![Go](https://github.com/ermanimer/retryablehttp/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/ermanimer/retryablehttp/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ermanimer/retryablehttp)](https://goreportcard.com/report/github.com/ermanimer/retryablehttp)

Simple HTTP client interface with automatic retries and constant backoff. Inspired by [HashiCorp](https://github.com/hashicorp)'s [go-retryablehttp](https://github.com/hashicorp/go-retryablehttp) library.

# Usage

NewClient() creates and returns a retryable HTTP client instance with provided options.

```go
c, err := retryablehttp.NewClient(
	retryablehttp.WithHTTPClient(http.DefaultClient),
	retryablehttp.WithMaxReqCount(3),
	retryablehttp.WithBackoff(100*time.Millisecond),
	retryablehttp.WithResHandler(func(res *http.Response) error {
		if res == nil {
            return ErrNilRes
        }

        statusCode := res.StatusCode
        if statusCode < 200 || statusCode > 299 {
            return ErrStatusCode(statusCode)
        }

        return nil
	}),
)
```

**WithHTTPClient** option configures underlying http client.

**WithMaxReqCount** option configures maximum request count.

**WithBackoff** option configures backoff duration which represents sleeping intervals between requests.

**WithResponseHandler** option configures response handler which handles responses.

Client has `Do(*http.Request) (*http.Response, error)` function which is identical to `*http.Client`. This makes our client broadly applicable with minimal effort.

```go
res, err := c.Do(req)
```

# Contribution

Any contribution or feedback is welcome.