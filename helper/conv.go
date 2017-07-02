package helper

func ConvInterfaceSliceToStringSlice(slice []interface{}) []string {
	outp := make([]string, len(slice))
	for k, v := range slice {
		outp[k], _ = v.(string)
	}
	return outp
}
