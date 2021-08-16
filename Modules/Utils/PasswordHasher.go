package Utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

//Source: https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72

func PasswordHasher(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
