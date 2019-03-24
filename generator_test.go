package feeder

import "testing"

func TestItemConvert(t *testing.T) {
	// Success empty struct
	item := Item{}
	item.convert()
}

func TestFeedConvert(t *testing.T) {
	// Success empty struct
	item := Feed{}
	item.convert()
}
