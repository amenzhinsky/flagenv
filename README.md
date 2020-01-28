# flagenv

Environment variables support for the standard `flag` package.

## Usage

The library helps to reduce amount of code usually written to use environment variables as fallbacks for command-line options:

```go
flag.UintVar(&timeoutFlag, "timeout", 10, "connection timeout in seconds [$TIMEOUT]")
flag.Parse()

flagIsSet := func(name string) bool {
	var set bool
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "timeout" {
			set = true
		}
	})
	return set
}

if s := os.Getenv("TIMEOUT"); s != "" && !flagIsSet("timeout") { // awkward way to detect unset flags
	timeout, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return fmt.Errorf("invalid value %q for flag -%s [$%s]: %s",
			s, "timeout", "TIMEOUT", err,
		)
	}
	timeoutFlag = uint(timeout)
}
```

Simply call `flagenv.Parse()` function after registering all necessary flags that wraps around `flag.Parse()`:

```go
flag.UintVar(&timeoutFlag, "timeout", 10, "connection timeout")
flagenv.Parse()
```

Environment is only applied to the flags that are not set through the command-line arguments.

It also adds environment variable names to each flag in the help message:

```
Usage of main:
  -timeout uint
        connection timeout [$TIMEOUT] (default 10)
```

## Configuration

In case you're using short flag names you may employ `WithMap` option to override environment variable names with more descriptive alternatives:

```go
flag.StringVar(&addrFlag, "a", ":8080", "address to connect to")
flag.UintVar(&timeoutFlag, "t", 0, "connection timeout")
flag.BoolVar(&forceFlag, "f", false, "force the action")
flagenv.Parse(flagenv.WithMap(func(name string) string {
	switch name {
	case "a":
		return "ADDR"
	case "t":
		return "TIMEOUT"
	case "f":
		// ignore '-f' flag
		return "" 
	default:
		// fall back to the default behaviour
		return flagenv.DefaultMap(name)
	}
}))
```

## Contributing

All contributions are welcome. Please fill in an issue before submitting pull requests.
