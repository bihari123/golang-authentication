package httpauth

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func ConvertToBase64() {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))
}
