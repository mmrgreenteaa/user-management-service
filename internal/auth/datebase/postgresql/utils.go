package postgresql

import (
	"crypto/sha256"
	"fmt"
)

// GenHash generate hash return hash string 
func GenHash(refresh string) string {

	data := []byte(refresh)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x",hash)
}
