package paginator

import (
	"fmt"
	"net/url"
	"strings"
)

// GetArrayFilters returns urlValues without the parameters that need to
// be used for Postgres array queries, and return the Postgres array
// query to be used
func GetArrayFilters(urlValues map[string][]string, filterFields []string) (map[string][]string, string) {
	arrayQuery := CreateArrayQuery(urlValues, filterFields)

	for _, field := range filterFields {
		delete(urlValues, field)
	}

	return urlValues, arrayQuery
}

// CreateArrayQuery will create the Postgres array query based on url.Values
// and the filterFields
func CreateArrayQuery(urlValues url.Values, filterFields []string) string {
	var stmts []string

	for _, field := range filterFields {
		val, ok := urlValues[field]
		if !ok {
			continue
		}

		stmts = append(stmts, createStatement(field, val))
	}

	// Create the query
	if len(stmts) == 1 {
		return stmts[0]
	} else if len(stmts) > 1 {
		return strings.Join(stmts, " AND ")
	} else {
		return ""
	}
}

// createStatement will create the sql part that can be used in
// a GORM Where clause
func createStatement(key string, value []string) string {
	var stmt string
	if len(value) > 1 {
		stmt = fmt.Sprintf("%s @> '{%s}'", key, strings.Join(value, "', '"))
	} else {
		stmt = fmt.Sprintf("%s @> '{%s}'", key, value[0])
	}
	return stmt
}

// transformArrayFiltering will transform:
//
// 	map[string][]string{"source": ["test1", "test2"], "status": ["crawled"]}
//
// to:
//
// 	map[string][]string{"source": ["test1", "test2"]}
func transformArrayFiltering(queryparams map[string][]string) (map[string][]string, map[string][]string) {
	arrayFilters := make(map[string][]string)

	for k, v := range queryparams {
		if len(v) > 0 {
			arrayFilters[k] = v
			delete(queryparams, k)
		}
	}

	return queryparams, arrayFilters
}

func GetSearchFilters(searchFields []string, searchTerm string) string {
	if searchTerm != "" && len(searchFields) > 0 {
		return CreateSearchQuery(searchFields, searchTerm)
	} else {
		return ""
	}
}

// CreateSearchQuery is a simple implementation of a LIKE statement based
// on provided searchFields.
func CreateSearchQuery(searchFields []string, searchTerm string) string {
	var stmts []string
	for _, field := range searchFields {
		stmts = append(stmts, fmt.Sprintf("lower(%s) LIKE lower('%%%s%%')", field, searchTerm))
	}

	if len(stmts) > 1 {
		return strings.Join(stmts, " OR ")
	} else if len(stmts) == 1 {
		return stmts[0]
	} else {
		return ""
	}
}
