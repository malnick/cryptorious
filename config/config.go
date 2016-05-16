package main

var (
	VERSION  string
	REVISION string
)

// Config{} is the complete application configuration
type Config struct {
	Version  string
	Revision string
}

//
func Configuration() (c Config, err error) {
	return c{
		Version:  VERSION,
		Revision: REVISION,
	}, nil
}
