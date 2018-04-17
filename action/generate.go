package action

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/malnick/cryptorious/config"
)

func checkFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

// GenerateKeys creates public private keys for a $USER
func GenerateKeys(c config.Config) error {
	// not implemented
	return nil
}

var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")

// NewPassword accepts a length and returns a randomized string value
func NewPassword(length int) error {
	p, err := randomPassword(length, stdChars)
	fmt.Println(p)
	return err

}

func randomPassword(length int, chars []byte) (string, error) {
	// https://raw.githubusercontent.com/cmiceli/password-generator-go/master/gen.go

	newPassword := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
			return string(newPassword), err
		}
		for _, c := range randomData {
			if c >= maxrb {
				continue
			}
			newPassword[i] = chars[c%clen]
			i++
			if i == length {
				return string(newPassword), nil
			}
		}
	}
}
