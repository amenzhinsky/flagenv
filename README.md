# flagenv

Environment variables support for golang's `flag` package.

## Usage

Tired of writing something like this in each project:

```go
flag.UintVar(&timeoutFlag, "timeout", 10, "connection timeout in seconds [$TIMEOUT]")
flag.Parse()

if s := os.Getenv("TIMEOUT"); s != "" &&
	timeoutFlag != 10 { // silly way to detect unset flags
	timeout, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		return fmt.Errorf("invalid value %q for flag -%s [$%s]: %s",
			s, "timeout", "TIMEOUT", err,
		)
	}
	timeoutFlag = uint(timeout)
}
```

Simply call `flagenv.Parse()` function after registering all necessary flags that wraps around `flag.Parse()` enabling environment variables as fallback values:

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
		// fall back to the standard behaviour
		return strings.ToUpper(name)
	}
}))
```
