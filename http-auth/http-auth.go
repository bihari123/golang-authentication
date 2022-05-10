package httpauth

import (
	"encoding/base64"
	"fmt"
)

func ConvertToBase64(){
  fmt.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))
}
