package flagenv

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestParseFrom_Env(t *testing.T) {
	fs, have := mkFlagSet(t, "666")
	if err := ParseWithEnv(fs, []string{}, WithMap(envName)); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 666)
}

func TestParseFromEnv_Empty(t *testing.T) {
	fs, have := mkFlagSet(t, "")
	if err := ParseWithEnv(fs, []string{}, WithMap(envName)); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 15)
}

func TestParseFromEnv_Error(t *testing.T) {
	fs, _ := mkFlagSet(t, "asdf")
	fs.SetOutput(ioutil.Discard) // no need to print usage to stderr
	if err := ParseWithEnv(fs, []string{}, WithMap(envName)); err == nil {
		t.Fatalf("err = nil, want any")
	}
}

func TestParseFromEnv_Argv(t *testing.T) {
	fs, have := mkFlagSet(t, "666")
	if err := ParseWithEnv(fs, []string{"-int", "333"}, WithMap(envName)); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 333)
}

func TestParseFromEnv_Ignore(t *testing.T) {
	fs, have := mkFlagSet(t, "111")
	if err := ParseWithEnv(fs, []string{}, WithMap(func(name string) string {
		return ""
	})); err != nil {
		t.Fatal(err)
	}
	testFlag(t, have, 15)
}

func mkFlagSet(t *testing.T, s string, opts ...Option) (*flag.FlagSet, *int) {
	t.Helper()
	if err := os.Setenv(envName("int"), s); err != nil {
		t.Fatal(err)
	}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	return fs, fs.Int("int", 15, "int var")
}

func envName(s string) string {
	return "__test_" + strings.ToUpper(s)
}

func testFlag(t *testing.T, have *int, want int) {
	t.Helper()
	if *have != want {
		t.Fatalf("flag = %d, want %d", *have, want)
	}
}
