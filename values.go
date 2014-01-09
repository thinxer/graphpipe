package graphpipe

func init() {
	Regsiter("Int", func(value *int) (int, error) {
		return *value, nil
	})
	Regsiter("Float64", func(value *float64) (float64, error) {
		return *value, nil
	})
	Regsiter("String", func(value *string) (string, error) {
		return *value, nil
	})
}
