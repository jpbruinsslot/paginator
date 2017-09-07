// Package paginator provides functions and wrappers to create paginated views
// for JSON API's
package paginator

import (
	"net/http"
	"strconv"
)

type Pagination struct {
	Count    int         `json:"count"`
	Next     *string     `json:"next"`     // pointer, to get null when empty
	Previous *string     `json:"previous"` // pointer, to get null when empty
	Results  interface{} `json:"results"`
}

// CreatePagination will wrap the results in a Pagination wrapper.
func CreatePagination(r *http.Request, results interface{}, limit, offset, count int) Pagination {
	return Pagination{
		Count:    count,
		Next:     createNextURL(r, offset, limit, count),
		Previous: createPreviousURL(r, offset, limit, count),
		Results:  results,
	}
}

// createNextURL will construct a formatted next url
func createNextURL(r *http.Request, offset, limit, count int) *string {
	var newOffset int
	if (count - (offset + limit)) > 0 {
		newOffset = offset + limit
	} else {
		newOffset = 0
	}

	// Create URL, use original http.Request
	if newOffset > 0 {
		values := r.URL.Query()
		values.Set("offset", strconv.Itoa(newOffset))
		r.URL.RawQuery = values.Encode()

		next := r.URL.String()
		return &next
	}

	return nil
}

// createPreviousURL will construct a formatted previous url
func createPreviousURL(r *http.Request, offset int, limit int, count int) *string {
	var newOffset int
	if (offset - limit) >= 0 {
		newOffset = (offset - limit)
	} else {
		newOffset = -1 // because 0 is a valid offset
	}

	// Create URL, use original http.Request
	if newOffset > -1 {
		values := r.URL.Query()
		values.Set("offset", strconv.Itoa(newOffset))
		r.URL.RawQuery = values.Encode()

		previous := r.URL.String()
		return &previous
	}

	return nil
}
