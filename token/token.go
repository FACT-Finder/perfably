package token

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

var (
	tokenCharacters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_")

	randReader = rand.Reader
)

func CreatePassword(pw string, strength int) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), strength)
	if err != nil {
		panic(err)
	}
	return hashedPassword
}

func ComparePassword(hashedPassword, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}

func GenerateRandomString(length int) string {
	res := make([]byte, length)
	for i := range res {
		index := randIntn(len(tokenCharacters))
		res[i] = tokenCharacters[index]
	}
	return string(res)
}

func randIntn(n int) int {
	max := big.NewInt(int64(n))
	res, err := rand.Int(randReader, max)
	if err != nil {
		panic("random source is not available")
	}
	return int(res.Int64())
}

func init() {
	randIntn(2)
}
