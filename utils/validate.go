package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	validateRegs = map[string]*regexp.Regexp{
		"email":     regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`),
		"alphaNum":  regexp.MustCompile("^[\\p{L}\\p{N}]+$"),
		"alphaNumu": regexp.MustCompile("^[\\p{L}\\p{N}\\._-]+$"),
	}

	validateRegsMsgs = map[string]string{
		"email":     "is not a valid email address",
		"alphaNum":  "is not a valid alphanumeric characters",
		"alphaNumu": "is not a valid alphanumeric andOr ._- characters",
	}

	reqErr = fmt.Errorf("required field")
)

// ValidateStruct validate struct return error array with key field name
func ValidateStruct(s interface{}) map[string]error {
	errs := map[string]error{}

	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get("valid")

		if tag == "" || tag == "-" {
			continue
		}

		tag = strings.Replace(tag, " ", "", -1)

		tags := strings.Split(tag, ",")
		req := false
		ty := tags[0]
		it := 0

		if ty == "req" {
			req = true
			ty = tags[1]
			it = 2
		} else {
			it = 1
		}

		var err error
		switch ty {
		case "alphaNum":
			min := 0
			max := 0
			fmt.Sscanf(strings.Join(tags[it:], ","), "min=%d,max=%d", &min, &max)
			err = validateString(v.Field(i).Interface().(string), min, max, req, "alphaNum")
		case "alphaNumu":
			min := 0
			max := 0
			fmt.Sscanf(strings.Join(tags[it:], ","), "min=%d,max=%d", &min, &max)
			err = validateString(v.Field(i).Interface().(string), min, max, req, "alphaNumu")
		case "email":
			err = validateEmail(v.Field(i).Interface().(string), req)
		case "num":
			min := -1
			max := -1
			fmt.Sscanf(strings.Join(tags[it:], ","), "min=%d,max=%d", &min, &max)
			err = validateNumber(v.Field(i).Interface().(int), min, max, req)
		}

		if err != nil {
			errs[v.Type().Field(i).Name] = err
		}
	}

	return errs
}

func validateString(val string, min int, max int, req bool, regkey string) error {
	val = strings.TrimSpace(val)
	l := utf8.RuneCountInString(val)

	if req && l == 0 {
		return reqErr
	}

	if min > 0 && l < min {
		return fmt.Errorf("should be at least %v chars long", min)
	}

	if min > 0 && max >= min && l > max {
		return fmt.Errorf("should be less than %v chars long", max)
	}

	if !validateRegs[regkey].MatchString(val) {
		return fmt.Errorf(validateRegsMsgs[regkey])
	}

	return nil
}

func validateEmail(val string, req bool) error {
	val = strings.TrimSpace(val)
	l := utf8.RuneCountInString(val)

	if req && l == 0 {
		return reqErr
	}

	if !validateRegs["email"].MatchString(val) {
		return fmt.Errorf(validateRegsMsgs["email"])
	}

	return nil
}

func validateNumber(val int, min int, max int, req bool) error {
	if req && val <= 0 {
		return reqErr
	} else if !req && (min == -1 && max == -1) {
		return nil
	}

	if val < min {
		return fmt.Errorf("should be greater than %v", min)
	}

	if max >= min && val > max {
		return fmt.Errorf("should be less than %v", max)
	}

	return nil
}
