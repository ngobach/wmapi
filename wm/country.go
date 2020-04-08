package wm

type Country string

func CountryFrom(s string) Country {
	// TODO: Validation here
	return Country(s)
}
