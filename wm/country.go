package wm

type Country string

func CountryFrom(s string) *Country {
	// TODO: Validation here
	if len(s) == 0 {
		s = "viet-nam"
	}
	if s == "world" {
		return nil
	}
	c := Country(s)
	return &c
}
