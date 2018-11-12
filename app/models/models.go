package models

// ErrMissingField error for missing field in model
type ErrMissingField string

func (e ErrMissingField) Error() string {
	return string(e) + " is required"
}
