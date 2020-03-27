package flagenv

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Option is a parser configuration option.
type Option func(*parser)

// WithMap sets callback that maps flag names to environment variable names.
func WithMap(fn MapFunc) Option {
	return func(p *parser) {
		p.mapping = fn
	}
}

// WithLookupEnv changes the default function that retrieves environment variables.
//
// Default is os.LookupEnv.
func WithLookupEnv(fn func(key string) (string, bool)) Option {
	return func(p *parser) {
		p.lookup = fn
	}
}

// Parse should be used in place of `flag.Parse()` to add support of
// environment variables to the default command-line flags parser.
func Parse(opts ...Option) {
	// CommandLine is set for ExitOnError, no need to handle errors.
	_ = ParseWithEnv(flag.CommandLine, os.Args[1:], opts...)
}

type parser struct {
	mapping MapFunc
	lookup  func(key string) (string, bool)
}

// DefaultMap is default variable mapping function,
// it replaces hyphens with underscores and upper-cases the result.
//
// Example: connect-timeout => CONNECT_TIMEOUT.
var DefaultMap = newDefaultMap()

func newDefaultMap() MapFunc {
	re := regexp.MustCompile("[^a-zA-Z0-9_]")
	return func(name string) string {
		return strings.ToUpper(re.ReplaceAllString(name, "_"))
	}
}

// MapFunc maps flag names to environment variable names.
//
// If returned value is an empty string the flag is ignored.
type MapFunc func(name string) string

// ParseWithEnv enables environment variables support for the given flag set.
// It adds environment variable names to each value usage string.
// Values from the environment only applied when flag hasn't been set from the
// command-line arguments and the value is not empty.
//
// Panics when fs is already parsed.
func ParseWithEnv(fs *flag.FlagSet, argv []string, opts ...Option) error {
	if fs.Parsed() {
		panic("already parsed")
	}

	p := &parser{mapping: DefaultMap, lookup: os.LookupEnv}
	for _, opt := range opts {
		opt(p)
	}

	// collect all the flags before parsing and remove ones that have been set
	m := map[string]*flag.Flag{}
	fs.VisitAll(func(f *flag.Flag) {
		name := p.mapping(f.Name)
		if name == "" {
			return
		}
		if existing, ok := m[name]; ok {
			panic(fmt.Sprintf("conflicting %q environment variable for -%s and -%s flags",
				name, existing.Name, f.Name))
		}

		m[f.Name] = f
		f.Usage = fmt.Sprintf("%s [$%s]", f.Usage, name)
	})
	if err := fs.Parse(argv); err != nil {
		return err
	}
	fs.Visit(func(f *flag.Flag) {
		delete(m, f.Name)
	})

	// repeat what `func (f *FlagSet) Parse(arguments []string) error` does,
	// only display env variable name next to flag name in error messages
	for _, f := range m {
		name := p.mapping(f.Name)
		s, ok := p.lookup(name)
		if !ok {
			continue
		}
		if err := f.Value.Set(s); err != nil {
			err = failf(fs, "invalid value %q for flag -%s [$%s]: %v",
				s, f.Name, name, err,
			)
			switch fs.ErrorHandling() {
			case flag.ContinueOnError:
				return err
			case flag.ExitOnError:
				os.Exit(2)
			case flag.PanicOnError:
				panic(err)
			}
		}
	}
	return nil
}

func failf(fs *flag.FlagSet, format string, v ...interface{}) error {
	err := fmt.Errorf(format, v...)
	fmt.Fprintln(fs.Output(), err)
	if fs.Usage != nil {
		fs.Usage()
	} else {
		if fs.Name() == "" {
			fmt.Fprintf(fs.Output(), "Usage:\n")
		} else {
			fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
		}
		fs.PrintDefaults()
	}
	return err
}
