package flagenv

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Parse should be used in place of `flag.Parse()` to add support of
// environment variables to the default command-line flags parser.
func Parse() {
	// CommandLine is set for ExitOnError, no need to handle errors.
	_ = ParseWithEnv(flag.CommandLine, os.Args[1:], nil)
}

// MapFunc is mapping function from flag name to environment variable name.
type MapFunc func(s string) string

// ParseWithEnv enables environment variables support for the given flag set.
// It adds environment variable names to each value usage string.
// Values from the environment only applied when flag hasn't been set from the
// command-line arguments and the value is not empty.
//
// Panics when fs is already parsed.
//
// If fn is nil `strings.ToUpper` is used.
func ParseWithEnv(fs *flag.FlagSet, argv []string, fn MapFunc) error {
	if fs.Parsed() {
		panic("already parsed")
	}
	if fn == nil {
		fn = strings.ToUpper
	}

	// collect all the flags before parsing and remove ones that have been set
	m := map[string]*flag.Flag{}
	fs.VisitAll(func(f *flag.Flag) {
		m[f.Name] = f
		f.Usage = f.Usage + " [$" + fn(f.Name) + "]"
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
		s := os.Getenv(fn(f.Name))
		if s == "" {
			continue
		}
		if err := f.Value.Set(s); err != nil {
			err = failf(fs, "invalid value %q for flag -%s [$%s]: %v",
				s, f.Name, fn(f.Name), err,
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
