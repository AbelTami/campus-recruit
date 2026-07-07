package hash

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string, pepper string) (string, error) {
	peppered := password + pepper
	hash, err := bcrypt.GenerateFromPassword([]byte(peppered), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password string, hash string, pepper string) bool {
	peppered := password + pepper
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(peppered))
	return err == nil
}
