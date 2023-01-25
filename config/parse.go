package config

var parsers [3][]func()

// Register a parser. This follows compile-time plugin pattern. f is expected
// to directly read and write config values and won't be called concurrent to
// any other caller.
func Register(priority int, f func()) {
	parsers[priority] = append(parsers[priority], f)
}

// MustParse the config.
func MustParse() {
	for priority := range parsers {
		for _, f := range parsers[priority] {
			f()
		}
	}
}
