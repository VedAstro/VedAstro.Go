package main

import (
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Tithi", func() (any, error) { return r.Client.LunarDay(r.Context, r.Now) }); err != nil {
			return err
		}
		if err := r.Show("Nakshatra", func() (any, error) { return r.Client.MoonConstellation(r.Context, r.Now) }); err != nil {
			return err
		}
		if err := r.Show("Yoga", func() (any, error) { return r.Client.NithyaYoga(r.Context, r.Now) }); err != nil {
			return err
		}
		if err := r.Show("Karana", func() (any, error) { return r.Client.Karana(r.Context, r.Now) }); err != nil {
			return err
		}
		return nil
	})
}
