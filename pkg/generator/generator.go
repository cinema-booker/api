package generator

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", r.Intn(10))
	}

	return code
}
