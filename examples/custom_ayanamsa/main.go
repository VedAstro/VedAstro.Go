package main

import (
	vedastro "github.com/VedAstro/VedAstro.Go"
	"github.com/VedAstro/VedAstro.Go/internal/exampleutil"
)

func main() {
	exampleutil.Run(func(r *exampleutil.Runner) error {
		for _, name := range []vedastro.Ayanamsa{vedastro.AyanamsaLahiri, vedastro.AyanamsaRaman, vedastro.AyanamsaKrishnamurti} {
			ctx := vedastro.WithAyanamsa(r.Context, name)
			if err := r.Show(string(name), func() (any, error) { return r.Client.AyanamsaDegree(ctx, r.Birth) }); err != nil {
				return err
			}
		}
		return nil
	})
}
