package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Ascendant sign", func() (any, error) { return r.Client.HouseSignName(r.Context, vedastro.House1, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("Planets in house 1", func() (any, error) { return r.Client.PlanetsInHouseBasedOnSign(r.Context, vedastro.House1, r.Birth) }); err != nil {
			return err
		}
		return nil
	})
}
