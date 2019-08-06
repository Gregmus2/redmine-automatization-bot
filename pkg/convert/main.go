package convert

// convert slice of strings to slice of interfaces
func Iface(list []string) []interface{} {
	vals := make([]interface{}, len(list))
	for i, v := range list {
		vals[i] = v
	}

	return vals
}
