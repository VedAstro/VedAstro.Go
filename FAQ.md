# FAQ

## Is this an astrology calculation engine?

It is a customer client for the remote VedAstro API. Calculations run on the server; you need network access. The SDK has no third-party Go dependencies.

## Why isn't there an npm or PyPI package?

Go downloads modules by their repository path and version tag. Install with `go get github.com/VedAstro/VedAstro.Go@v1.0.0` and import package `vedastro`. GitHub hosts the source; Go's proxy caches releases.

## How does it compare with Python?

It exposes the same 684 calculator names and required argument order, symbolic enums, Time/GeoLocation JSON, and payload unwrapping. Go adds explicit contexts and errors, instance-specific settings, and method-specific optional structs. Four Python optional defaults incorrectly use None instead of the server's empty string; Go omits those fields by default.

## Why do methods return any?

Results vary by calculator and may be objects, arrays, scalar values, or SVG strings. Check the error, then use a type assertion or type switch. JSON numbers are float64. See the README's result table.

## How do I cancel a call?

Pass a context from `context.WithCancel` or `context.WithTimeout`. Use `errors.Is` to detect cancellation or deadline expiry. There is **no deadline by default**: a calculation can take milliseconds or minutes, and the client cannot know what is acceptable for your workload, so a built-in limit would only ever truncate a valid answer. Add `WithTimeout` only when your own code has decided a call has run too long. HTTP and API failures use `*APIError`. Calls do not automatically retry.

## How do optional parameters work?

Omit the options struct to use server defaults, or pass one struct with pointer fields. Use `vedastro.Ptr(0)`, `vedastro.Ptr(false)`, or `vedastro.Ptr("")` to send explicit values. A nil field is omitted.

## Can I use multiple keys or ayanamsas concurrently?

Yes. Create separate clients for separate keys. Set a default with `WithDefaultAyanamsa` and override per call with `WithAyanamsa(ctx, value)`. Settings do not mutate global state. Free-tier limits still apply across calls.

## Does creating a client contact the network?

No. It prints the existing ASCII banner once per process. Only calculations and an explicit `CheckForUpdate` contact remote services.

## Can the SDK update itself?

No. `CheckForUpdate` returns release information and a `go get` command. Update your application's dependency, rebuild, and restart it. Already-running executables retain their compiled SDK version.

## Where are all the examples?

See [examples/README.md](examples/README.md). They cover the Python demo topics, including charts, panchanga, compatibility, dasa, transits, numerology, source-text search, and exports.
