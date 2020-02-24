package flagenv

import (
	"flag"
	"io/ioutil"
	"testing"
)

func TestSimple(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	have := fs.Int("int", 15, "int var")
	if err := ParseWithEnv(fs, nil, withEnv(map[string]string{
		"INT": "666",
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 666)
}

func TestParseWithMapCustom(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	have := fs.Int("int", 15, "int var")
	if err := ParseWithEnv(fs, nil, withEnv(map[string]string{
		"CUSTOM": "666",
	}), WithMap(func(name string) string {
		switch name {
		case "int":
			return "CUSTOM"
		default:
			return ""
		}
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 666)
}

func TestParseWithMapIgnore(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	have := fs.Int("int", 15, "int var")
	if err := ParseWithEnv(fs, nil, WithMap(func(name string) string {
		return ""
	}), withEnv(map[string]string{
		"INT": "111",
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 15)
}

func TestParseEmpty(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	have := fs.Int("int", 15, "int var")
	if err := ParseWithEnv(fs, nil, withEnv(map[string]string{
		// empty
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 15)
}

func TestParseError(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_ = fs.Int("int", 15, "int var")
	fs.SetOutput(ioutil.Discard) // no need to print usage to stderr
	if err := ParseWithEnv(fs, nil, withEnv(map[string]string{
		"INT": "asdf",
	})); err == nil {
		t.Fatalf("err = nil, want any")
	}
}

func TestParseArgvOvertakes(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	have := fs.Int("int", 15, "int var")
	if err := ParseWithEnv(fs, []string{"-int", "333"}, withEnv(map[string]string{
		"INT": "666",
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 333)
}

func TestNameConflict(t *testing.T) {
	defer testPanic(t)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	_ = fs.Bool("d", false, "debug")
	_ = fs.Bool("D", false, "detach")
	if err := ParseWithEnv(fs, nil); err != nil {
		t.Fatal(err)
	}
}

func TestAlreadyParsed(t *testing.T) {
	defer testPanic(t)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	if err := fs.Parse(nil); err != nil {
		t.Fatal(err)
	}
	if err := ParseWithEnv(fs, nil); err != nil {
		t.Fatal(err)
	}
}

func withEnv(env map[string]string) Option {
	return WithLookupEnv(func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	})
}

func testFlag(t *testing.T, have *int, want int) {
	t.Helper()
	if *have != want {
		t.Fatalf("flag = %d, want %d", *have, want)
	}
}

func testPanic(t *testing.T) {
	t.Helper()
	if err := recover(); err == nil {
		t.Fatal("want a panic")
	}
}
