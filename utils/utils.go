package utils

import (
	"time"

	"github.com/google/uuid"
)

var TimeNowUTC = func() time.Time {
	return time.Now().UTC()
}

var NewUUIDString = func() string {
	return uuid.New().String()
}
