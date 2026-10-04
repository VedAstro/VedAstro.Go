# Quickstart

Install Go 1.26 or later, then create a consumer module:

```sh
mkdir my-vedastro-app
cd my-vedastro-app
go mod init example.com/my-vedastro-app
go get github.com/VedAstro/VedAstro.Go@v1.0.0
```

Copy [examples/quick_start/main.go](examples/quick_start/main.go) into `main.go` and run:

```sh
go run .
```

The example sends one live API request and prints the Sun's sign for the supplied birth time. A network connection is required.

Replace `"FreeAPIUser"` with `os.Getenv("VEDASTRO_API_KEY")` (add the `os` import) to use your key. For free access, wait at least 13 seconds between calls, across all clients.

To run the other examples, clone this repository and run `go run ./examples/daily_panchanga` from its root. Most examples handle Ctrl+C and free-tier pacing. Set `VEDASTRO_API_KEY` in your shell to use a paid key.

Read [README](README.md) for decoded result types, optional settings, ayanamsa, deadlines, and errors.
