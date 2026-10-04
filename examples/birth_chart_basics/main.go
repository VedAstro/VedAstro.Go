package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Sun sign", func() (any, error) { return r.Client.PlanetRasiD1Sign(r.Context, vedastro.PlanetSun, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("Moon sign", func() (any, error) { return r.Client.PlanetRasiD1Sign(r.Context, vedastro.PlanetMoon, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("Ascendant", func() (any, error) { return r.Client.HouseSignName(r.Context, vedastro.House1, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("Birth constellation", func() (any, error) { return r.Client.MoonConstellation(r.Context, r.Birth) }); err != nil {
			return err
		}
		return nil
	})
}
