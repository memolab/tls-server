package tests

import (
	"strings"
	"testing"
	"tls-server/utils"
)

func TestValidationStruct(t *testing.T) {

	usr1 := struct {
		ID       int    `valid:"num"`
		Username string `valid:"req,alphaNum, min=5,max=60"`
		Email    string `valid:"req,email"`
		Password string `valid:"req,alphaNumu, min=5,max=60"`
	}{Email: "@m.com",
		ID:       -1,
		Username: "meomeo",
		Password: "Passs%"}

	errs := utils.ValidateStruct(usr1)

	if errs["ID"] != "" {
		t.Errorf("got %s, Expected empty", errs["ID"])
	}

	if errs["Username"] != "" {
		t.Errorf("got %s, Expected empty", errs["username"])
	}

	if errs["Email"] == "" {
		t.Errorf("Expected error email field")
	} else if errs["Email"] != "is not a valid email address" {
		t.Errorf("Expected error email (is not a valid email address)")
	}

	if errs["Password"] == "" {
		t.Errorf("Expected error Password field")
	} else if !strings.HasPrefix(errs["Password"], "is not a valid alphanumeric") {
		t.Errorf("got %s, Expected has (is not a valid alphanumeric)", errs["Password"])
	}

}
