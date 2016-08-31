package main

import "testing"

func TestSignParam(t *testing.T) {
	appKey := "test_app_key"
	appSecret := "test_app_secret"
	params := map[string]string{
		"foo":     "bar",
		"egg":     "spam",
		keyDate:   "2016-05-27T09:15:39.559Z",
		keyAppKey: appKey,
	}
	expected := "89985e2420899196db5bdf16b3c2ed0922c0c221"
	actual := generateSignature(params, appSecret)
	if actual != expected {
		t.Error(
			"For", params,
			"expected", expected,
			"got", actual,
		)
	}
}
