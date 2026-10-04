package main

import (
	"fmt"
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		// Add customer records here. Every call is paced for the free API key.
		for index, birth := range []vedastro.Time{r.Birth, r.Other} {
			if err := r.Show(fmt.Sprintf("Chart %d", index+1), func() (any, error) {
				return r.Client.PlanetRasiD1Sign(r.Context, vedastro.PlanetSun, birth)
			}); err != nil {
				return err
			}
		}
		return nil
	})
}
