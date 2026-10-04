package exampleutil

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	vedastro "github.com/VedAstro/VedAstro.Go"
)

// Runner shares explicit sample inputs and paces free-tier calls.
type Runner struct {
	Client            *vedastro.Client
	Context           context.Context
	Birth, Other, Now vedastro.Time
	lastCall          time.Time
	delay             time.Duration
}

// Run reports errors and supports Ctrl+C cancellation for long examples.
func Run(example func(*Runner) error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	key := os.Getenv("VEDASTRO_API_KEY")
	delay := time.Duration(0)
	if key == "" || key == "FreeAPIUser" {
		key = "FreeAPIUser"
		delay = 13 * time.Second
	}
	geo := vedastro.NewGeoLocation("Mumbai", 72.8777, 19.076)
	birth, err := vedastro.NewTime("14:30 25/10/1992 +05:30", geo)
	if err != nil {
		panic(err)
	}
	other, err := vedastro.NewTime("09:00 01/01/1990 -05:00", vedastro.NewGeoLocation("New York", -74.006, 40.7128))
	if err != nil {
		panic(err)
	}
	now, err := vedastro.NewTimeFromDate(time.Now().In(time.FixedZone("IST", 19800)), geo)
	if err != nil {
		panic(err)
	}
	r := &Runner{Client: vedastro.NewClient(key), Context: ctx, Birth: birth, Other: other, Now: now, delay: delay}
	if err := example(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Fetch waits before starting the next free-tier request.
func (r *Runner) Fetch(call func() (any, error)) (any, error) {
	if wait := r.delay - time.Since(r.lastCall); wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-r.Context.Done():
			return nil, r.Context.Err()
		}
	}
	r.lastCall = time.Now()
	return call()
}

// Show prints a labeled, indented API result.
func (r *Runner) Show(label string, call func() (any, error)) error {
	value, err := r.Fetch(call)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s:\n%s\n", label, data)
	return nil
}
