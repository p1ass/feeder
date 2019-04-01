package feeder_test

import (
	"github.com/naoki-kishi/feeder"
	"testing"
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
