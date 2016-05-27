// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package stringutil

import (
	"crypto/rand"
	"io"
)

// Return part of a string.
func SubString(s string, start, length int) string {
	rs := []rune(s)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

// Make a string's first character is uppercase.
func UpperFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	c := s[0]
	if ('a' <= c) && (c <= 'z') {
		return string(rune(int(c) - 32)) + SubString(s, 1, len(s) - 1)
	}
	return s
}

// Make a string's first character is lowercase.
func LowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	c := s[0]
	if ('A' <= c) && (c <= 'Z') {
		return string(rune(int(c) + 32)) + SubString(s, 1, len(s) - 1)
	}
	return s
}

// Generate random string.
// @length length of random string.
func GenerateRandomString(len int) string {
	return string(GenerateRandomByte(len))
}

// Generate random []byte.
// @length length of []byte.
func GenerateRandomByte(len int) []byte {
	b := make([]byte, len)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil
	}
	return b
}