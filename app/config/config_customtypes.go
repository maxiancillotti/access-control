package config

import (
	"time"

	"github.com/pkg/errors"
)

// Purpose:
// 1) Reap the benefit of being able to use time duration string format in config variables.
// 2) Validate format when parsing the config file / env var and not having to do it
// always before passing each value of the type to a builder or contructor func.
type StrTimeDuration string

// Implements Setter interface needed by library cleanenv to set the custom type
func (d *StrTimeDuration) SetValue(s string) error {
	if s == "" {
		return errors.New("value of a type StrTimeDuration cannot be empty")
	}
	_, err := time.ParseDuration(s)
	if err != nil {
		return errors.Wrap(err, "string format invalid, could not be parsed to time.Duration")
	}
	*d = StrTimeDuration(s)
	return nil
}

func (d *StrTimeDuration) GetDuration() time.Duration {
	parsed, _ := time.ParseDuration(string(*d))
	// error already checked at SetValue
	return parsed
}
