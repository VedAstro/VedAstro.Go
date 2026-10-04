package main

import (
	"context"
	"errors"
	"fmt"
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		// Cancellation is checked before the HTTP request; this intentionally makes no API call.
		ctx, cancel := context.WithCancel(r.Context)
		cancel()
		_, err := r.Client.PlanetRasiD1Sign(ctx, vedastro.PlanetSun, r.Birth)
		if errors.Is(err, context.Canceled) {
			fmt.Println("Canceled request handled")
		} else {
			return fmt.Errorf("expected cancellation, got %v", err)
		}
		err = r.Show("Sun sign", func() (any, error) { return r.Client.PlanetRasiD1Sign(r.Context, vedastro.PlanetSun, r.Birth) })
		var apiErr *vedastro.APIError
		if errors.As(err, &apiErr) {
			fmt.Printf("API endpoint %s returned HTTP %d\n", apiErr.Endpoint, apiErr.StatusCode)
		}
		return err
	})
}
