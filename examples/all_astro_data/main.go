package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Planet data", func() (any, error) { return r.Client.AllPlanetData(r.Context, vedastro.PlanetSun, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("House data", func() (any, error) { return r.Client.AllHouseData(r.Context, vedastro.House1, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("Zodiac data", func() (any, error) { return r.Client.AllZodiacSignData(r.Context, vedastro.ZodiacGemini, r.Birth) }); err != nil {
			return err
		}
		return nil
	})
}
