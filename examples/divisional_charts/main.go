package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("D1", func() (any, error) { return r.Client.AllHouseRasiSigns(r.Context, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("D9", func() (any, error) { return r.Client.AllHouseNavamshaSign(r.Context, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("D10", func() (any, error) { return r.Client.AllHouseDashamamshaSign(r.Context, r.Birth) }); err != nil {
			return err
		}
		if err := r.Show("Sun D9", func() (any, error) { return r.Client.PlanetNavamshaD9Sign(r.Context, vedastro.PlanetSun, r.Birth) }); err != nil {
			return err
		}
		return nil
	})
}
