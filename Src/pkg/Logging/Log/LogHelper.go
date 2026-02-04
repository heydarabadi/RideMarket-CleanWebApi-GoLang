package Log

func MapToZapParams(extra map[Extrakey]interface{}) []interface{} {
	params := make([]interface{}, 0)
	for k, v := range extra {
		params = append(params, string(k))
		params = append(params, v)
	}
	return params
}
