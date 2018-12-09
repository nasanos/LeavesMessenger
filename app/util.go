package main

import (
	"math/rand"
    "time"
)

func generateHex(n int) (string) {
	seededRandom := rand.New(rand.NewSource(time.Now().UnixNano()))

	hex := make([]byte, n)
	for i := range hex {
		hex[i] = hexCharSet[seededRandom.Intn(len(hexCharSet))]
	}

	return string(hex)
}

func checkAuthentication(token string, branchKey string) string {
    var authentication string = ""

    if branchKey == "" {
        authentication = selectAuthentication(token)
    } else {
        authentication = selectAuthenticationForBranch(token, branchKey)
    }

    return authentication
}
