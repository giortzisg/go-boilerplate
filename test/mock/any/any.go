package any

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type Time struct{}

func (a Time) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (a Time) Value() (driver.Value, error) {
	return time.Now(), nil
}

type UUID struct{}

func (a UUID) Match(v driver.Value) bool {
	_, ok := v.(uuid.UUID)
	return ok
}

func (a UUID) Value() (driver.Value, error) {
	return uuid.New(), nil
}
