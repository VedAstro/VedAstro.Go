package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("All Sun data", func() (any, error) { return r.Client.AllPlanetData(r.Context, vedastro.PlanetSun, r.Birth) }); err != nil {
			return err
		}
		return nil
	})
}
