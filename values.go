package graphpipe

func init() {
	Register("Int", func(value *int) (int, error) {
		return *value, nil
	})
	Register("Float64", func(value *float64) (float64, error) {
		return *value, nil
	})
	Register("String", func(value *string) (string, error) {
		return *value, nil
	})
}
