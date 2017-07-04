package helper

func SafeMap(m map[string]interface{}, key string, def interface{}) interface{} {
	var outp interface{} = def
	switch def.(type) {
	case string:
		if val, ok := m[key].(string); ok {
			outp = val
		}
	case int:
		if val, ok := m[key].(int); ok {
			outp = val
		}
	case []interface{}:
		if val, ok := m[key].([]interface{}); ok {
			outp = val
		}
	}
	return outp
}
