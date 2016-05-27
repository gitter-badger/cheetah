// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package cheetah

import (
	"encoding/base64"
	"github.com/HeadwindFly/cheetah/utils/string"
	"strings"
)

// Generate CSRF token.
// @maskLen length of mask.
// @token true token.
func GenerateCsrfToken(maskLen int, token []byte) string {
	// Generate mask string.
	mask := stringutil.GenerateRandomByte(maskLen)

	// XOR
	tokenByte := xorCsrfTokens(token, mask)

	// Base64 encoding.
	tokenStr := base64.StdEncoding.EncodeToString([]byte(string(mask) + string(tokenByte)))

	return strings.Replace(tokenStr, "+", ".", -1)
}

// Validate CSRF token.
// @maskLen length of mask.
// @token the token was generated by method named GenerateCsrfToken().
// @trueToken true token.
func ValidateCsrfToken(maskLen int, token, trueToken string) bool {
	// Restore the original base64 encoding string.
	token = strings.Replace(token, ".", "+", -1)

	// Base64 decoding.
	tokenByte, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		panic(err)
	}

	// If the token is invalid， returns false.
	token = string(tokenByte)
	if len(token) <= maskLen {
		return false
	}

	// Get mask by maskLen.
	mask := []byte(token)[0:maskLen]
	tokenByte = []byte(token)[maskLen:]

	// XOR
	token = string(xorCsrfTokens(mask, tokenByte))

	// Return true if the token is equals to trueToken.
	if 0 == strings.Compare(token, trueToken) {
		return true
	}

	// Otherwise false will be returned.
	return false
}

// XOR
func xorCsrfTokens(token1, token2 []byte) []byte {
	len1 := len(token1)
	len2 := len(token2)
	if len1 > len2 {
		for i := 0; i < len1-len2; i++ {
			token2 = append(token2, token2[i%len2])
		}
	} else {
		for i := 0; i < len2-len1; i++ {
			if len1 == 0 {
				token1 = append(token1, ' ')
			} else {
				token1 = append(token1, token1[i%len1])
			}
		}
	}
	token := []byte{}
	for i := 0; i < len(token1); i++ {
		token = append(token, token1[i]^token2[i])
	}
	return token
}
