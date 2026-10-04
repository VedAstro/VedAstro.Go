# VedAstro Go

Customer SDK for the [VedAstro API](https://vedastro.org). All 684 calculator methods run on the remote service. The SDK uses only the Go standard library and requires Go 1.26 or later.

## Install

```sh
go get github.com/VedAstro/VedAstro.Go@v1.0.0
```

Import as `vedastro "github.com/VedAstro/VedAstro.Go"`. See [QUICKSTART](QUICKSTART.md), [FAQ](FAQ.md), and [21 runnable examples](examples/README.md).

## First calculation

```go
package main

import (
    "context"
    "fmt"
    "log"

    vedastro "github.com/VedAstro/VedAstro.Go"
)

func main() {
    client := vedastro.NewClient("FreeAPIUser")
    birth, err := vedastro.NewTime("14:30 25/10/1992 +05:30",
        vedastro.NewGeoLocation("Mumbai", 72.8777, 19.0760))
    if err != nil {
        log.Fatal(err)
    }
    result, err := client.PlanetRasiD1Sign(context.Background(), vedastro.PlanetSun, birth)
    if err != nil {
        log.Fatal(err)
    }
    sign, ok := result.(map[string]any)
    if !ok {
        log.Fatal("unexpected Sun sign response")
    }
    fmt.Println("Sun sign:", sign["Name"])
}
```

Every calculator takes `context.Context` first and returns `(any, error)`. Required arguments retain the Python calculator's order and names. Browse [package documentation](https://pkg.go.dev/github.com/VedAstro/VedAstro.Go) or `go doc` for signatures.

## Keys, contexts, and configuration

Use `NewClient(os.Getenv("VEDASTRO_API_KEY"))` for your own key. An empty key leaves authentication to the service's free tier. Never embed paid keys in distributed source. Free access currently allows five requests per minute; space requests at least 13 seconds apart, including requests across clients. The SDK does not automatically retry or throttle. Most examples include pacing.

Each client owns its settings and supports concurrent requests. The default calculation timeout is 120 seconds; an earlier context deadline wins.

```go
client := vedastro.NewClient(apiKey,
    vedastro.WithTimeout(30*time.Second),
    vedastro.WithDefaultAyanamsa(vedastro.AyanamsaLahiri),
)
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
result, err := client.PlanetRasiD1Sign(ctx, vedastro.PlanetSun, birth)
```

`WithHTTPClient(*http.Client)` accepts a custom HTTP client, retaining a shallow copy. Shared custom transports must be safe for concurrent use. `WithBaseURL(url)` selects a different Calculate service URL. Invalid configuration is returned as an error when a calculator is called.

Use `errors.Is(err, context.Canceled)` or `errors.Is(err, context.DeadlineExceeded)` for interruption. Remote HTTP/API and malformed-response failures use `*vedastro.APIError`, available through `errors.As`. Transport errors omit request credentials; API failure messages redact the configured key.

## Decoded results

The API's `Payload` is decoded into standard Go JSON values:

| JSON | Go |
| --- | --- |
| Object | `map[string]any` |
| Array | `[]any` |
| Number | `float64` |
| String / Boolean / null | `string` / `bool` / `nil` |
| SVG response | `string` |

As in Python, a payload object with exactly one property is unwrapped to that property's value. Zero, false, empty strings, empty collections, and null are valid results. Check the error before asserting a result type. Use `json.MarshalIndent` to display or save dynamic results.

## Optional arguments

Methods with optional parameters accept zero or one method-specific options struct. Nil fields are omitted so the server supplies its defaults. `Ptr` preserves explicit zero, false, and empty values:

```go
result, err := client.SearchSourceText(ctx, "Saturn",
    vedastro.SearchSourceTextOptions{
        TopK: vedastro.Ptr(3),
        SourceName: vedastro.Ptr("Brihat-Parashara-Hora-Shastra"),
    })
```

Passing more than one options value returns an error without sending a request. Defaults are documented on generated fields and in [api_manifest.json](api_manifest.json). Four Prashna Marga methods have an empty-string server default that Python represents as None; Go intentionally uses the server default when omitted.

## Time, location, and enums

`NewGeoLocation(name, longitude, latitude)` takes longitude first. `NewTime` validates the explicit UTC-offset format `HH:MM DD/MM/YYYY ±HH:MM`. `NewTimeFromDate(time.Time, location)` preserves the supplied date's offset. Locations are validated when serialized. Helpers serialize the same `StdTime`, `Location`, `Name`, `Longitude`, and `Latitude` fields as Python.

Typed constants include `PlanetSun`, `House1`, `ZodiacAries`, and `AyanamsaLahiri`. They serialize the API's symbolic strings.

A context override affects only that call and its children, keeping concurrent requests independent:

```go
raman := vedastro.WithAyanamsa(ctx, vedastro.AyanamsaRaman)
result, err := client.AyanamsaDegree(raman, birth)
```

An empty override selects the server default even if the client has a default. Omission currently selects Lahiri on the server.

## Updates

The existing ASCII banner prints once when the first client is created. No network update check occurs at startup.

```go
info, err := client.CheckForUpdate(ctx)
if err == nil && info.Available {
    fmt.Println(info.UpdateCommand)
    fmt.Println(info.RebuildInstructions)
}
```

This explicit check queries Go's public module proxy. It never sends the API key or modifies files. Run the returned `go get` command, rebuild your application, and restart it to use a newer SDK.

## Development and generation

`calculate_generated.go`, `enums_generated.go`, and `api_manifest.json` come from **StaticTableGenerator Task 11** in the main VedAstro repository. Do not edit generated files. Coverage is derived from metadata, with first-overload selection matching the other SDKs; the initial baseline is 684 methods.

From the main repository:

```sh
dotnet run --project StaticTableGenerator -- --task 11
```

Output defaults to the sibling `VedAstro.Go-PUBLIC` folder. Set `VEDASTRO_GO_OUT` to override it and optionally `VEDASTRO_GOFMT` to the gofmt executable. Isolated Task 11 uses local metadata only: no API, Cosmos, or LLM calls. Output is deterministic and formatted before writing.

From this repository:

```sh
gofmt -w .
go vet ./...
go test -timeout 90s ./...
go build ./...
go test -race -timeout 120s ./...
python scripts/check_python_parity.py ../VedAstro.Python-PUBLIC/vedastro/calculate.py
```

Race checks require a supported platform and C toolchain. CI tests Go 1.26 and 1.27 on Windows and Linux, including Linux race checks. Tests use local HTTP servers; examples call the live API only when explicitly run.

Release with a new semantic version tag after validation. Never move or replace an existing release tag. GitHub hosts source, the Go module proxy distributes tagged versions, and pkg.go.dev indexes documentation. See [Go's publishing guide](https://go.dev/doc/modules/publishing).

MIT licensed.
