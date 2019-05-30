package flagenv

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func TestParseFrom_Env(t *testing.T) {
	fs, have := mkFlagSet(t, "666")
	if err := ParseWithEnv(fs, []string{}, envName); err != nil {
		t.Fatal(err)
	}
	if *have != 666 {
		t.Fatalf("flag = %d, want %d", *have, 666)
	}
}

func TestParseFromEnv_Empty(t *testing.T) {
	fs, have := mkFlagSet(t, "")
	if err := ParseWithEnv(fs, []string{}, envName); err != nil {
		t.Fatal(err)
	}
	if *have != 15 {
		t.Fatalf("flag = %d, want %d", *have, 15)
	}
}

func TestParseFromEnv_Error(t *testing.T) {
	fs, _ := mkFlagSet(t, "asdf")
	if err := ParseWithEnv(fs, []string{}, envName); err == nil {
		t.Fatalf("err = nil, want any")
	}
}

func TestParseFromEnv_Argv(t *testing.T) {
	fs, have := mkFlagSet(t, "666")
	if err := ParseWithEnv(fs, []string{"-int", "333"}, envName); err != nil {
		t.Fatal(err)
	}
	if *have != 333 {
		t.Fatalf("flag = %d, want %d", *have, 333)
	}
}

func mkFlagSet(t *testing.T, s string) (*flag.FlagSet, *int) {
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
