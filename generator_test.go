package feeder_test

import (
	"testing"

	"github.com/p1ass/feeder"
)

func TestItemConvert(t *testing.T) {
	// Success empty struct
	item := feeder.Item{}
	feeder.ExportItemConvert(&item)
}

func TestFeedConvert(t *testing.T) {
	// Success empty struct
	item := feeder.Feed{}
	feeder.ExportFeedConvert(&item)
}
