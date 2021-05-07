package env

import (
	"os"
	"strings"
)

func GetEnv() string {
	return strings.ToLower(os.Getenv("ENV"))
}

func IsTest() bool {
	if GetEnv() == "prod" || GetEnv() == "pre" {
		return false
	}
	return true
}

func IsProd() bool {
	if GetEnv() == "prod" {
		return true
	}
	return false
}

func IsPre() bool {
	if GetEnv() == "pre" {
		return true
	}
	return false
}
