package main

import (
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		if err := r.Show("Compatibility report", func() (any, error) { return r.Client.MatchReport(r.Context, r.Birth, r.Other) }); err != nil {
			return err
		}
		return nil
	})
}
