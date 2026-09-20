package utils

import (
	"fmt"
	"math/rand/v2"
)

func num() int {
	return rand.IntN(8) + 1
}

// otp used for verifying emails on signup
func GetOtp() string {
	return fmt.Sprintf("%v%v%v%v%v%v", num(), num(), num(), num(), num(), num())
}
