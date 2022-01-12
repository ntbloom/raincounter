package fetch

type Fetcher interface {
	// Fetch blocks until all of the data have been retrieved from their respective places
	Fetch() interface{}
}
