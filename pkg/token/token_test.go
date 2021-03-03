package token

import "testing"

func TestGetToken(t *testing.T) {
	err, _ := GetToken("NEED TO ADD", "CREDS")
	if err != "" {
		t.Fail()
	}
}
