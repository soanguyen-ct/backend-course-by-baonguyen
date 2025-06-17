package hashpass

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
)

var saltRegex = regexp.MustCompile("(?:[$]([1-9][0-9]*)[$])?(.*)")
var saltHashRegex = regexp.MustCompile(`^\$(\d+)\$(.*)`)

const STRETCH_LOOPS = 1024
const SALT_LENGTH = 16
const ALPHANUMERIC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genSalt() string {
	b := make([]byte, SALT_LENGTH)
	for i := range b {
		b[i] = ALPHANUMERIC[rand.Intn(len(ALPHANUMERIC))]
	}
	return string(b)
}

func blocketSalt(salt, password string) string {

	group := saltRegex.FindStringSubmatch(salt)

	stretchLoops := STRETCH_LOOPS
	if len(group) > 3 {
		stretchLoops, _ = strconv.Atoi(group[1])
		salt = group[2]
	}
	hashPassword := password

	for i := 0; i < stretchLoops; i++ {
		hs := sha1.New()
		hs.Write([]byte(salt + hashPassword))
		hashPassword = hex.EncodeToString(hs.Sum(nil))
	}

	hashPassword = fmt.Sprintf("$%d$%s%s", stretchLoops, salt, hashPassword)

	return hashPassword
}

func HashPassword(password string) string {
	salt := genSalt()
	return blocketSalt(salt, password)
}

func HashPasswordLogin(password, hashPassword string) string {
	groups := saltHashRegex.FindStringSubmatch(hashPassword)
	if len(groups) <= 2 {
		return "-1"
	}
	salt := groups[2][0:SALT_LENGTH]

	return blocketSalt(salt, password)
}
