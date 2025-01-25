package pb

import "strings"

// Convert back = and ' to normal characters so that pocketbase can parse them
func updateParams(encodedParams string) (updatedQuery string) {
	updatedQuery = strings.ReplaceAll(encodedParams, "%3D", "=")
	updatedQuery = strings.ReplaceAll(updatedQuery, "%27", "'")
	return updatedQuery
}
