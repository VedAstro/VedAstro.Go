package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Current Sun sign", func() (any, error) { return r.Client.PlanetRasiD1Sign(r.Context, vedastro.PlanetSun, r.Now) }); err != nil {
			return err
		}
		if err := r.Show("Current Saturn sign", func() (any, error) { return r.Client.PlanetRasiD1Sign(r.Context, vedastro.PlanetSaturn, r.Now) }); err != nil {
			return err
		}
		if err := r.Show("Saturn retrograde", func() (any, error) { return r.Client.IsPlanetRetrograde(r.Context, vedastro.PlanetSaturn, r.Now) }); err != nil {
			return err
		}
		return nil
	})
}
