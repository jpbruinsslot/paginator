package paginator

import (
	"fmt"
	"net/http"
	"strconv"
)

// ParseQueryParams will return sensible values for pagination, ordering and
// filtering from the queryparameters that are passed from a http request. It
// will tranform the input values to values that can be used for querying the
// storage backend you're using.
func ParseQueryParams(r *http.Request) (filter map[string]interface{}, search string, ordering string, offset int, limit int) {
	// r.URL.Query() returns a map[string][]string which contain the values. We
	// will first try to get the pagination, and ordering parameters and remove
	// them from the map. What is left will be the filtering parameters.

	// Copy queryparams, otherwise we would alter the original map
	queryparams := make(map[string][]string)
	for k, v := range r.URL.Query() {
		queryparams[k] = v
	}

	// Pagination
	queryparams, offset, limit = transformPagination(queryparams)

	// Ordering
	queryparams, ordering = transformOrdering(queryparams)

	// Search
	queryparams, search = transformSearching(queryparams)

	// Filtering
	filter = transformFiltering(queryparams)

	return
}

// transformPagination will transform:
//
// 		map[string][]string{"offset": ["0"], "limit": ["100"]
//
// to:
//
// 		var retOffset int = 0
// 		var retLimit int = 100
func transformPagination(queryparams map[string][]string) (map[string][]string, int, int) {
	var retOffset int
	var retLimit int

	// Pagination, offset (?offset=10)
	offset := queryparams["offset"]
	if offset != nil {
		retOffset, _ = strconv.Atoi(queryparams["offset"][0])
	} else {
		retOffset = 0
	}
	delete(queryparams, "offset")

	// Pagination, limit (?limit=10)
	limit := queryparams["limit"]
	if limit != nil {
		retLimit, _ = strconv.Atoi(queryparams["limit"][0])
	} else {
		// default, could be in config
		retLimit = 100
	}
	delete(queryparams, "limit")

	return queryparams, retOffset, retLimit
}

// transformOrdering will transform:
//
//		map[string][]string{"ordering": ["-id"]}}
//
// to:
//
//		var retOrdering string = "id desc"
//
func transformOrdering(queryparams map[string][]string) (map[string][]string, string) {
	var retOrdering string
	ordering := queryparams["ordering"]

	if ordering != nil {
		strOrdering := ordering[0]

		var order string
		var field string

		// uncover ascending/descending
		if string(strOrdering[0]) == "-" {
			field = string(strOrdering[1:])
			order = "desc"
		} else {
			field = strOrdering
			order = "asc"
		}

		retOrdering = fmt.Sprintf("%s %s", field, order)

	}
	delete(queryparams, "ordering")

	return queryparams, retOrdering
}

// transformSearching will transform:
//
//		map[string][]string{"search": ["hello, world"]}}
//
// to:
//
//		var searchTerm string = "hello, world"
//
func transformSearching(queryparams map[string][]string) (map[string][]string, string) {

	var searchTerm string
	search := queryparams["search"]

	if search != nil {
		searchTerm = search[0]
	}
	delete(queryparams, "search")

	return queryparams, searchTerm
}

// transformOrdering will transform:
//
//		map[string][]string{"source": ["test"], "status": ["crawled"]}}
//
// to:
//
//		map[string]interface{}{"source": "test", "status": "crawled"}
//
func transformFiltering(queryparams map[string][]string) map[string]interface{} {
	outputMap := make(map[string]interface{})

	for k, v := range queryparams {
		outputMap[k] = v[0]
	}

	return outputMap
}
