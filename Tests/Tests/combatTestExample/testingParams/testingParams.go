package testingParams

import (
	"TestBot/aas"
)

type TestingParams struct {
	SessionName string
	Url         string
	Locales     []string
	Tags        []string
	Browser     *aas.AAS
}

type TestingParamsForATest struct {
	SessionName string
	Url         string
	Locale      string
	Browser     *aas.AAS
}
