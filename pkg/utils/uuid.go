package utils

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func GenerateTrackingID() string {
	val, err := uuid.NewRandom()
	if err != nil {
		// fallback should uuid return an error for any reason
		rand.Seed(time.Now().UnixNano())
		min := 1000000000
		max := 9999999999
		randomNumber := rand.Intn(max-min+1) + min
		return fmt.Sprintf("%d", randomNumber)
	}
	return val.String()
}
