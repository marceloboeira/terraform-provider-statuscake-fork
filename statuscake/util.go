package statuscake

func castSetToSliceStrings(configured []interface{}) []string {
	res := make([]string, len(configured))

	for i, element := range configured {
		res[i] = element.(string)
	}
	return res
}
