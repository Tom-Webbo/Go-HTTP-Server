package main

import "strings"

func cleaner(body string) string {
	stringSlice := strings.Split(body, " ")
	for i, word := range stringSlice {
		if strings.ToLower(word) == "kerfuffle" {
			stringSlice[i] = "****"
		}
		if strings.ToLower(word) == "sharbert" {
			stringSlice[i] = "****"
		}
		if strings.ToLower(word) == "fornax" {
			stringSlice[i] = "****"
		}
	}
	msg := strings.Join(stringSlice, " ")
	return msg
}
