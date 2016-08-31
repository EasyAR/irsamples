package main

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	keyDate      = "date"
	keyAppKey    = "appKey"
	keySignature = "signature"
)

func sha1Hex(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func generateSignature(params map[string]string, appSecret string) string {
	keys := []string{}
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	parts := keys
	for i, k := range keys {
		parts[i] = k + params[k]
	}

	paramStr := strings.Join(parts, "")

	return sha1Hex(paramStr + appSecret)
}

func signParam(params map[string]string, appKey string, appSecret string) map[string]string {
	params[keyDate] = time.Now().UTC().Format(time.RFC3339Nano)
	params[keyAppKey] = appKey
	params[keySignature] = generateSignature(params, appSecret)
	return params
}
