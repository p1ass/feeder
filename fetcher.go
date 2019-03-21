package feeder

// Fetcher is ...
type Fetcher interface {
	Fetch() (*Items, error)
}
