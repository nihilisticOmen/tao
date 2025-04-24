package jwts

import "testing"

func TestParseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDYwMjgwODksInRva2VuS2V5IjoiMTAwMCJ9.BKlhvCiEeZ99nNzmBgU69s_qtzNULBu6vnNzVdldblQ"
	ParseToken(tokenString, "msproject")
}
