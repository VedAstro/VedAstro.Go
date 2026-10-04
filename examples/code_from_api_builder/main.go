package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("House occupied by the Sun", func() (any, error) {
			return r.Client.HousePlanetOccupiesBasedOnSign(r.Context, vedastro.PlanetSun, r.Birth)
		}); err != nil {
			return err
		}
		return nil
	})
}
