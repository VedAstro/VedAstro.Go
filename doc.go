// Package vedastro is the customer HTTP client for VedAstro's hosted astrology API.
//
// Calls require an internet connection. Results are JSON decoded into map[string]any,
// []any, string, float64, bool, or nil. SVG chart endpoints return strings.
// A payload containing a single object property is unwrapped to that property's value.
// Free-tier examples use FreeAPIUser and space calls at least 13 seconds apart.
//
// Clients can be shared by goroutines. Configure each client when constructing it;
// use WithAyanamsa to override the default for an individual context.
package vedastro
