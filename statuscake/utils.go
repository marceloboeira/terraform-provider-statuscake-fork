package statuscake

import (
	"fmt"
	"strings"
)

func toStringSlice(a []interface{}) []string {
	return strings.Split(strings.Trim(strings.Replace(fmt.Sprint(a), " ", ",", -1), "[]"), ",")
}

func getUserAgent(currentUserAgent string) string {
	return "terraform-provider-statuscake " + currentUserAgent
}
