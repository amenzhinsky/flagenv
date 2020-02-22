# flagenv

Environment variables support for the standard `flag` package.

## Usage

The library helps to reduce amount of code usually written to use environment variables as fallbacks for command-line options.

Simply replace `flag.Parse()` in your main application with `flagenv.Parse()` to add the functionality (the module also supports passing `*flag.FlagSet` to `flagenv.ParseWithEnv` for lower-level usage).

Environment is only applied to the flags that are not set through the command-line arguments.

Environment variable names are added to each flag in the help message:

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
