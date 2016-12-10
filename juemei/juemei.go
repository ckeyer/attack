package juemei

const (
	RootURL  = "http://www.juemei.com"
	MaxWorks = 1
)

var (
	BlackPrefix = []string{
		RootURL,
	}
	IgnorePrefix = []string{
		"#",
	}
)

type ResolveResult struct {
	URL      string
	Imgs     []string
	Links    []string
	OutLinks []string
}

type StoreResult struct {
	URLs                  []string
	Imgs, Links, OutLinks int64
}

type Store interface {
}
