package flagenv

import (
	"flag"
	"io/ioutil"
	"testing"
)

func TestParseFrom_Env(t *testing.T) {
	fs, have := mkFlagSetAndIntFlag(t)
	if err := ParseWithEnv(fs, nil, withEnv(map[string]string{
		"INT": "666",
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 666)
}

func TestParseFromEnv_Empty(t *testing.T) {
	fs, have := mkFlagSetAndIntFlag(t)
	if err := ParseWithEnv(fs, nil, withEnv(map[string]string{
		// empty
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 15)
}

func TestParseFromEnv_Error(t *testing.T) {
	fs, _ := mkFlagSetAndIntFlag(t)
	fs.SetOutput(ioutil.Discard) // no need to print usage to stderr
	if err := ParseWithEnv(fs, nil, withEnv(map[string]string{
		"INT": "asdf",
	})); err == nil {
		t.Fatalf("err = nil, want any")
	}
}

func TestParseFromEnv_Argv(t *testing.T) {
	fs, have := mkFlagSetAndIntFlag(t)
	if err := ParseWithEnv(fs, []string{"-int", "333"}, withEnv(map[string]string{
		"INT": "666",
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 333)
}

func TestParseFromEnv_Ignore(t *testing.T) {
	fs, have := mkFlagSetAndIntFlag(t)
	if err := ParseWithEnv(fs, nil, WithMap(func(name string) string {
		return ""
	}), withEnv(map[string]string{
		"INT": "111",
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 15)
}

func mkFlagSetAndIntFlag(t *testing.T) (*flag.FlagSet, *int) {
	t.Helper()
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	return fs, fs.Int("int", 15, "int var")
}

func withEnv(m map[string]string) Option {
	return WithGetenv(func(key string) string {
		return m[key]
	})
}

func testFlag(t *testing.T, have *int, want int) {
	t.Helper()
	if *have != want {
		t.Fatalf("flag = %d, want %d", *have, want)
	}
}
