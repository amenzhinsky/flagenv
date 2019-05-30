# flagenv

Environment variables support for golang's `flag` package.

## Usage

If you're tired of writing something like this in each project:

```go
flag.UintVar(&timeoutFlag, "timeout", 10, "connection timeout in seconds [$TIMEOUT]")
flag.Parse()

if s := os.Getenv("TIMEOUT"); s != "" &&
	timeoutFlag != 10 { // silly way to detect unset flags
	timeout, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return fmt.Errorf("infalid value %q for flag -%s [$%s]: %s",
			s, "timeout", "TIMEOUT", err,
		)
	}
	timeoutFlag = uint(timeout)
}
```

Simply after registering all necessary flags call `flagenv.Parse()` function that wraps around `flag.Parse()`:

```go
flag.UintVar(&timeoutFlag, "timeout", 10, "connection timeout")
flagenv.Parse()
```

Environment variables are only applied to the flags that are not set through the command-line arguments.

It also adds environment variable names to each flag in the help message:

```
Usage of main:
  -timeout uint
        connection timeout [$TIMEOUT] (default 10)
```

For using the package with `*flag.FlagSet` see `ParseWithEnv` function that supports mapping customization of flag names to env names.
