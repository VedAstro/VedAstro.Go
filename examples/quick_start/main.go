package main

import (
	"context"
	"fmt"
	"log"

	vedastro "github.com/VedAstro/VedAstro.Go"
)

func main() {
	client := vedastro.NewClient("FreeAPIUser")
	birth, err := vedastro.NewTime("14:30 25/10/1992 +05:30",
		vedastro.NewGeoLocation("Mumbai", 72.8777, 19.0760))
	if err != nil {
		log.Fatal(err)
	}
	result, err := client.PlanetRasiD1Sign(context.Background(), vedastro.PlanetSun, birth)
	if err != nil {
		log.Fatal(err)
	}
	sign, ok := result.(map[string]any)
	if !ok {
		log.Fatal("unexpected Sun sign response")
	}
	fmt.Println("Sun sign:", sign["Name"])
}
