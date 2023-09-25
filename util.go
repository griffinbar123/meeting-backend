package main

import (
	"math/rand"
)

func RangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}

func remove(s []ClientInRoom, i int) []ClientInRoom {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeRoom(s []Room, i int) []Room {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func removeActiveClient(s []ActiveClient, i int) []ActiveClient {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
