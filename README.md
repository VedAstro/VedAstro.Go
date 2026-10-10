<h1 align="center">🪐 VedAstro for Go</h1>

<p align="center">
  <em>The most comprehensive Vedic astrology library for Go — 684 calculations, one <code>context</code> away.</em>
</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/VedAstro/VedAstro.Go"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?style=flat-square&logo=go&logoColor=white" alt="Go reference"/></a>
  <a href="https://pkg.go.dev/github.com/VedAstro/VedAstro.Go"><img src="https://img.shields.io/badge/go-1.26%2B-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go version"/></a>
  <a href="https://github.com/VedAstro/VedAstro.Go/blob/main/LICENSE"><img src="https://img.shields.io/github/license/VedAstro/VedAstro.Go?style=flat-square&color=fda085" alt="License"/></a>
  <a href="https://github.com/VedAstro/VedAstro.Go/stargazers"><img src="https://img.shields.io/github/stars/VedAstro/VedAstro.Go?style=flat-square&color=f6d365" alt="Stars"/></a>
  <a href="https://github.com/VedAstro/VedAstro.Go/actions"><img src="https://img.shields.io/github/actions/workflow/status/VedAstro/VedAstro.Go/ci.yml?style=flat-square&label=CI" alt="CI"/></a>
</p>

---

### 🎯 Built for Real Apps

Perfect for:
- 📱 **Horoscope services** - Get 200+ life predictions instantly
- 💑 **Marriage matching backends** - 16-factor Kuta compatibility analysis
- 📅 **Daily panchanga APIs** - Tithi, Nakshatra, Yoga, Karana
- 🔮 **AI astrology chatbots** - Natural language birth chart queries
- 📊 **Astrological research** - Batch process thousands of charts
- 🌟 **Numerology calculators** - Chaldean system with life aspect scores

**Standard library only.** The SDK has zero third-party dependencies, compiles into a single static binary, and ships typed constants for every enum — so calculators, planets, signs and ayanamsas autocomplete in your editor.

---

## 🏃 Quick Start (3 minutes to your first calculation)

### Installation (10 seconds)

```sh
go get github.com/VedAstro/VedAstro.Go@v1.0.0
```

Requires **Go 1.26 or newer**. Import the package as `vedastro`:

```go
import vedastro "github.com/VedAstro/VedAstro.Go"
```

**No C toolchain, no ephemeris files, no generated protobufs.** Calculations run on VedAstro's servers over HTTPS.

### Your First Calculation (20 seconds)

```go
package main

import (
	"context"
	"fmt"
	"log"

	vedastro "github.com/VedAstro/VedAstro.Go"
)

func main() {
	// 'FreeAPIUser' is the free-tier key; an empty key also uses the free tier.
	client := vedastro.NewClient("FreeAPIUser")

	// Time format: "HH:MM DD/MM/YYYY +TZ:TZ"
	// GeoLocation takes (name, longitude, latitude) — longitude first.
	birth, err := vedastro.NewTime("14:30 25/10/1992 +05:30",
		vedastro.NewGeoLocation("Mumbai", 72.8777, 19.0760))
	if err != nil {
		log.Fatal(err)
	}

	// Every calculator takes a context first and returns (any, error).
	result, err := client.PlanetRasiD1Sign(context.Background(), vedastro.PlanetSun, birth)
	if err != nil {
		log.Fatal(err)
	}

	sign := result.(map[string]any)
	fmt.Println("Sun sign:", sign["Name"]) // Output: Libra
}
```

**That's it! You just made your first Vedic astrology calculation.** 🎉

> **Results are dynamically typed.** Calculators return `any` because payloads vary — objects, arrays, scalars, or SVG strings. Check `err` first, then type-assert. See [Decoded results](#-decoded-results).

---

## 📚 Beginner-Friendly Examples

Every snippet below is complete and runnable: each one builds the `vedastro.Time` it needs, then calls `client`, which comes from the [Quick Start](#-quick-start-3-minutes-to-your-first-calculation). This helper reads one property from an object payload without repeating type assertions:

```go
// field reads one property from an object payload, or returns "" if absent.
func field(result any, key string) string {
	if object, ok := result.(map[string]any); ok {
		if value, ok := object[key].(string); ok {
			return value
		}
	}
	return ""
}
```

### Example 1: Get Birth Chart Basics

**Use Case:** Display the Sun, Moon, and Ascendant signs

```go
ctx := context.Background()

birth, err := vedastro.NewTime("14:30 25/10/1992 +05:30",
	vedastro.NewGeoLocation("Mumbai", 72.8777, 19.0760))
if err != nil {
	log.Fatal(err)
}

sun, err := client.PlanetRasiD1Sign(ctx, vedastro.PlanetSun, birth)
if err != nil {
	log.Fatal(err)
}

moon, err := client.PlanetRasiD1Sign(ctx, vedastro.PlanetMoon, birth)
if err != nil {
	log.Fatal(err)
}

// HouseSignName takes a typed house constant; it unwraps to a plain string.
rising, err := client.HouseSignName(ctx, vedastro.House1, birth)
if err != nil {
	log.Fatal(err)
}

fmt.Println("Sun:   ", field(sun, "Name"))   // e.g. Libra
fmt.Println("Moon:  ", field(moon, "Name"))  // e.g. Scorpio
fmt.Println("Rising:", rising)               // e.g. Capricorn
```

👉 **See full example:** [`examples/birth_chart_basics`](examples/birth_chart_basics/main.go)

---

### Example 2: Marriage Compatibility

**Use Case:** Check if two people are compatible for marriage

```go
person1, _ := vedastro.NewTime("23:40 31/12/1996 +09:00",
	vedastro.NewGeoLocation("Tokyo", 139.83, 35.65))
person2, _ := vedastro.NewTime("14:30 15/06/1997 -05:00",
	vedastro.NewGeoLocation("New York", -74.006, 40.7128))

report, err := client.MatchReport(context.Background(), person1, person2)
if err != nil {
	log.Fatal(err)
}

match := report.(map[string]any)
fmt.Printf("Compatibility: %v/100\n", match["KutaScore"])

summary := match["Summary"].(map[string]any)
fmt.Println(summary["ScoreSummary"])

// Each Kuta factor carries a Nature of Good, Bad, or Neutral.
for _, item := range match["PredictionList"].([]any) {
	kuta := item.(map[string]any)
	fmt.Printf("  • %v: %v\n", kuta["Name"], kuta["Nature"])
}
```

**Output** (shape and values depend on the two charts):

```
Compatibility: 65/100
Good match - Near perfect match, overall happiness
  • Graha Maitram: Good
  • Rajju: Good
  • Nadi Kuta: Good
  • Vasya Kuta: Bad
  • Dina Kuta: Good
```

👉 **See full example:** [`examples/marriage_compatibility`](examples/marriage_compatibility/main.go)

---

### Example 3: Current Planetary Positions

**Use Case:** Get today's planetary positions for any location

```go
mumbai := vedastro.NewGeoLocation("Mumbai", 72.8777, 19.0760)

// NewTimeFromDate keeps the supplied time's calendar fields and UTC offset.
now, err := vedastro.NewTimeFromDate(time.Now(), mumbai)
if err != nil {
	log.Fatal(err)
}

planets := []vedastro.PlanetName{
	vedastro.PlanetSun, vedastro.PlanetMoon, vedastro.PlanetMars,
	vedastro.PlanetMercury, vedastro.PlanetJupiter, vedastro.PlanetVenus,
	vedastro.PlanetSaturn, vedastro.PlanetRahu, vedastro.PlanetKetu,
}

fmt.Println("Current planetary positions:")
for _, planet := range planets {
	sign, err := client.PlanetRasiD1Sign(context.Background(), planet, now)
	if err != nil {
		log.Printf("skipping %s: %v", planet, err)
		continue
	}
	star, err := client.PlanetConstellation(context.Background(), planet, now)
	if err != nil {
		log.Printf("skipping %s: %v", planet, err)
		continue
	}
	fmt.Printf("  %s: %s in %v\n", planet, field(sign, "Name"), star)
}
```

`PlanetName` is a named `string` type, so `%s` prints it directly.

**Output:**
```
Current planetary positions:
  Sun: Taurus in Krithika - 2
  Moon: Sagittarius in Moola - 1
  ...
```

`PlanetConstellation` unwraps to a string that includes the pada, such as `Chitta - 3`.

> Free-tier pacing: this loop makes two calls per planet. See [Pacing requests](#-pacing-requests-and-rate-limits).

👉 **See full example:** [`examples/current_planets`](examples/current_planets/main.go)

---

### Example 4: Daily Panchanga

**Use Case:** Get today's Tithi, Nakshatra, Yoga, and Karana

```go
mumbai := vedastro.NewGeoLocation("Mumbai", 72.8777, 19.0760)
today, err := vedastro.NewTimeFromDate(time.Now(), mumbai)
if err != nil {
	log.Fatal(err)
}

ctx := context.Background()

tithi, err := client.LunarDay(ctx, today)
if err != nil {
	log.Fatal(err)
}
nakshatra, err := client.MoonConstellation(ctx, today)
if err != nil {
	log.Fatal(err)
}
yoga, err := client.NithyaYoga(ctx, today)
if err != nil {
	log.Fatal(err)
}
karana, err := client.Karana(ctx, today)
if err != nil {
	log.Fatal(err)
}

// Shapes differ: LunarDay and NithyaYoga return objects, the others unwrap to strings.
day := tithi.(map[string]any)
yogaInfo := yoga.(map[string]any)

fmt.Println("Tithi:    ", day["Name"], "("+day["Paksha"].(string)+")") // e.g. Amavasya (Krishna)
fmt.Println("Nakshatra:", nakshatra)                                   // e.g. Chitta - 3
fmt.Println("Yoga:     ", yogaInfo["Name"])                            // e.g. Vishkambha
fmt.Println("Karana:   ", karana)                                      // e.g. Chatushpada
```

**Output:**
```
Tithi:     Amavasya (Krishna)
Nakshatra: Chitta - 3
Yoga:      Vishkambha
Karana:    Chatushpada
```

**Note:** shapes vary per calculator — `LunarDay` (`Name`, `Paksha`, `Date`, `Day`, `Phase`) and `NithyaYoga` (`Name`, `Description`) decode to objects, while `MoonConstellation` and `Karana` unwrap to strings. Check a signature with `go doc`, and prefer the two-value type assertion while exploring.

**Note:** `NewTimeFromDate` preserves the *time.Time's* own offset. Pass a value in the target zone (`time.Now().In(time.FixedZone("IST", 19800))`) rather than relabelling a local time.

👉 **See full example:** [`examples/daily_panchanga`](examples/daily_panchanga/main.go)

---

### Example 5: Vimshottari Dasa Timeline

**Use Case:** Get planetary periods (Mahadasa → Bhukti → Antaram) over a range

```go
mumbai := vedastro.NewGeoLocation("Mumbai", 72.8777, 19.0760)
birth, _ := vedastro.NewTime("14:30 25/10/1992 +05:30", mumbai)
start, _ := vedastro.NewTime("00:00 01/01/2020 +05:30", mumbai)
end, _ := vedastro.NewTime("23:59 31/12/2030 +05:30", mumbai)

// Optional arguments use a method-specific options struct with pointer fields.
// (DasaAtRangeOptions.Levels and .PrecisionHours are both *int.)
dasa, err := client.DasaAtRange(context.Background(), birth, start, end,
	vedastro.DasaAtRangeOptions{
		Levels:         vedastro.Ptr(3),
		PrecisionHours: vedastro.Ptr(100),
	})
if err != nil {
	log.Fatal(err)
}

// Dynamic payloads print cleanly with json.MarshalIndent.
data, err := json.MarshalIndent(dasa, "", "  ")
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(data))
```

👉 **See full example:** [`examples/vimshottari_dasa`](examples/vimshottari_dasa/main.go)

---

### Example 6: Numerology Calculator

**Use Case:** Get Chaldean numerology analysis for a name

```go
name := "John Doe"
ctx := context.Background()

birthNumber, err := client.BirthNumber(ctx, birth)
if err != nil {
	log.Fatal(err)
}
destinyNumber, err := client.DestinyNumber(ctx, birth)
if err != nil {
	log.Fatal(err)
}

prediction, err := client.NameNumberPrediction(ctx, name)
if err != nil {
	log.Fatal(err)
}

result := prediction.(map[string]any)
fmt.Println("Birth number:  ", birthNumber)
fmt.Println("Destiny number:", destinyNumber)
fmt.Printf("%s: number %v (root %v, ruling planet %v)\n",
	name, result["Number"], result["RootNumber"], result["Planet"])

// "Prediction" is HTML-formatted interpretation text; drop the tags for a terminal.
html, _ := result["Prediction"].(string)
fmt.Println(regexp.MustCompile(`<[^>]+>`).ReplaceAllString(html, ""))
```

**Output:**
```
Birth number:   7
Destiny number: 2
John Doe: number 34 (root 7, ruling planet Ketu)
This number has the potential to be seen as lucky, ...
```

`Prediction` is HTML-formatted — strip the tags for terminal output (this snippet needs the `regexp` import), or render it in a browser.

👉 **See full example:** [`examples/numerology_calculator`](examples/numerology_calculator/main.go)

---

### Example 7: Semantic Search of Classical Texts (RAG)

**Use Case:** Ask a plain-English question and get the most relevant passages from classical Vedic books (BPHS, Phaladeepika, Hindu Predictive Astrology, …) — no exact keywords needed

```go
ctx := context.Background()

texts, err := client.GetAvailableSourceTexts(ctx)
if err != nil {
	log.Fatal(err)
}
fmt.Println("Searchable texts:", texts)

// Natural-language semantic search, narrowed to one book.
passages, err := client.SearchSourceText(ctx, "effects of Saturn in the 7th house",
	vedastro.SearchSourceTextOptions{
		TopK:       vedastro.Ptr(3),
		SourceName: vedastro.Ptr("Hindu-Predictive-Astrology"),
	})
if err != nil {
	log.Fatal(err)
}

for _, item := range passages.([]any) {
	passage := item.(map[string]any)
	// Lower score = closer match, so invert it into a relevance percentage.
	score, _ := passage["score"].(float64)
	// Some chunks carry no page number, so tolerate a nil.
	page, _ := passage["pageNumber"].(float64)
	fmt.Printf("%v p.%.0f (%.0f%%)\n", passage["sourceName"], page, (1-score)*100)
	fmt.Printf("   %v\n\n", passage["text"])
}
```

**Output:**
```
Searchable texts: [Brihat-Jataka Brihat-Parashara-Hora-Shastra Hindu-Predictive-Astrology Hora-Sara Jaimini-Sutras Phaladeepika Uttara-Kalamrita]
Brihat-Parashara-Hora-Shastra p.381 (28%)
   -------+ In the above case, the 2nd lord is Saturn. ...
```

> 🤖 This is the same retrieval step that powers VedAstro's RAG/AI features — feed the returned passages into an LLM prompt to build a **cited** astrology chatbot.

👉 **See full example:** [`examples/rag_vedic_books`](examples/rag_vedic_books/main.go)

---

## 🎓 Step-by-Step Tutorials

### Tutorial 1: Understanding Time Format

**The Most Common Beginner Issue: Time Format**

**Format:** `"HH:MM DD/MM/YYYY +TZ:TZ"` (24-hour, DD/MM/YYYY, explicit UTC offset — no timezone names)

| Location | Example | Timezone Offset |
|----------|---------|-----------------|
| India | `"14:30 25/10/1992 +05:30"` | IST = UTC+5:30 |
| USA (East) | `"09:30 25/10/1992 -05:00"` | EST = UTC-5:00 |
| USA (West) | `"06:30 25/10/1992 -08:00"` | PST = UTC-8:00 |
| Japan | `"23:30 25/10/1992 +09:00"` | JST = UTC+9:00 |
| UK | `"14:30 25/10/1992 +00:00"` | GMT = UTC+0:00 |
| Australia | `"00:30 26/10/1992 +10:00"` | AEST = UTC+10:00 |

**Two Ways to Create a Time:**

```go
mumbai := vedastro.NewGeoLocation("Mumbai", 72.8777, 19.0760)

// Method 1: API time string (validated immediately)
birth, err := vedastro.NewTime("14:30 25/10/1992 +05:30", mumbai)

// Method 2: from a Go time.Time, keeping its calendar fields and offset
birth, err = vedastro.NewTimeFromDate(time.Date(1992, 10, 25, 14, 30, 0, 0, ist), mumbai)
```

**Common Mistakes:**

```go
// ❌ Date order must be DD/MM/YYYY
vedastro.NewTime("14:30 1992-10-25 +05:30", mumbai)

// ❌ AM/PM is not supported (use 24-hour)
vedastro.NewTime("2:30 PM 25/10/1992 +05:30", mumbai)

// ❌ Missing UTC offset
vedastro.NewTime("14:30 25/10/1992", mumbai)

// ❌ Named zones are not accepted; use an offset
vedastro.NewTime("14:30 25/10/1992 IST", mumbai)

// ✅ Correct — NewTime returns an error you must handle
birth, err := vedastro.NewTime("14:30 25/10/1992 +05:30", mumbai)
```

> `NewGeoLocation(name, longitude, latitude)` takes **longitude first**. Coordinates and the time string are validated when the value is serialized, so an invalid `Time` fails at call time with a clear error rather than producing a wrong chart. `NewTime` itself validates eagerly.

---

### Tutorial 2: Choosing the Right Ayanamsa

**What is Ayanamsa?**

Ayanamsa is the difference between the tropical (Western) and sidereal (Vedic) zodiacs. Different systems shift positions by up to ~3°, so a planet near a sign boundary can change sign.

**47 systems are available.** Common choices:

| Constant | When to Use |
|----------|-------------|
| `vedastro.AyanamsaLahiri` | Indian government standard, most widely used — **the server default when none is set** |
| `vedastro.AyanamsaRaman` | Popular in South India |
| `vedastro.AyanamsaKrishnamurti` | KP (Krishnamurti Paddhati) system |
| `vedastro.AyanamsaFaganBradley` | Western sidereal astrology |
| `vedastro.AyanamsaYukteshwar` | Sri Yukteswar's calculation |
| …42 more | Autocomplete `vedastro.Ayanamsa` in your editor |

**Three ways to set it:**

```go
// 1. Per call, via context. Affects this call and its children only,
//    so concurrent requests stay independent.
raman := vedastro.WithAyanamsa(context.Background(), vedastro.AyanamsaRaman)
degree, err := client.AyanamsaDegree(raman, birth)

// 2. Per client, as the default for every request from this client.
client := vedastro.NewClient(apiKey,
	vedastro.WithDefaultAyanamsa(vedastro.AyanamsaKrishnamurti))

// 3. Unset: the server applies its own default, which measures as Lahiri.
```

A context override **wins over** the client default, and an empty override value explicitly asks for the server default even when the client has one.

**Recommendation:**
- 🇮🇳 **Indian astrology** → `vedastro.AyanamsaLahiri`
- 🌏 **Unsure** → omit it (server default measures as Lahiri)
- 📐 **KP system** → `vedastro.AyanamsaKrishnamurti`
- 🌍 **Western sidereal** → `vedastro.AyanamsaFaganBradley`

---

### Tutorial 3: Interpreting Match Reports

**Understanding Kuta Score:**

| Score Range | Compatibility | Recommendation |
|-------------|--------------|----------------|
| 33-36 points | Excellent | Highly compatible |
| 25-32 points | Good | Compatible, proceed with confidence |
| 18-24 points | Average | Requires careful consideration |
| Below 18 | Poor | Not recommended without other factors |

**The 16 Kutas Explained:**

1. **Graha Maitram (5 pts)** - Mental compatibility, happiness
2. **Gana (6 pts)** - Temperament match (Deva/Manushya/Rakshasa)
3. **Yoni (4 pts)** - Sexual compatibility
4. **Nadi (8 pts)** - Health & progeny (most important!)
5. **Varna (1 pt)** - Spiritual/ego compatibility
6. …and 11 more factors

**Reading the Report:**

```go
report, err := client.MatchReport(ctx, person1, person2)
if err != nil {
	log.Fatal(err)
}
match := report.(map[string]any)

fmt.Println(match["KutaScore"])                              // normalised to 100
fmt.Println(match["Summary"].(map[string]any)["ScoreSummary"])

for _, item := range match["PredictionList"].([]any) {
	kuta := item.(map[string]any)
	fmt.Printf("%v: %v\n", kuta["Name"], kuta["Nature"])     // Good / Bad / Neutral
	fmt.Printf("  Info: %v\n", kuta["Info"])
}
```

> Type assertions panic if the shape differs. For untrusted input, use the two-value form (`value, ok := x.(T)`) or a type switch instead.

---

## 🔧 Common Use Cases & Demo Files

21 runnable programs live in [`examples/`](examples/README.md). Run any of them from the repository root:

```sh
go run ./examples/quick_start
go run ./examples/daily_panchanga
```

| Example | Use Case | Skill Level |
|---------|----------|-------------|
| [`quick_start`](examples/quick_start/main.go) | First calculation and decoded result | Beginner |
| [`birth_chart_basics`](examples/birth_chart_basics/main.go) | Sun, Moon, ascendant, birth constellation | Beginner |
| [`marriage_compatibility`](examples/marriage_compatibility/main.go) | Compatibility calculation | Beginner |
| [`match_checker`](examples/match_checker/main.go) | Full match report | Beginner |
| [`current_planets`](examples/current_planets/main.go) | Positions at the current time | Beginner |
| [`daily_panchanga`](examples/daily_panchanga/main.go) | Tithi, Nakshatra, Yoga, Karana | Beginner |
| [`numerology_calculator`](examples/numerology_calculator/main.go) | Name and birth numerology | Beginner |
| [`horoscope_prediction_names`](examples/horoscope_prediction_names/main.go) | Prediction names | Beginner |
| [`custom_ayanamsa`](examples/custom_ayanamsa/main.go) | Lahiri, Raman, Krishnamurti context overrides | Intermediate |
| [`divisional_charts`](examples/divisional_charts/main.go) | Divisional (varga) signs | Intermediate |
| [`transit_analysis`](examples/transit_analysis/main.go) | Planetary transits | Intermediate |
| [`vimshottari_dasa`](examples/vimshottari_dasa/main.go) | Dasa over a date range | Intermediate |
| [`all_planet_data`](examples/all_planet_data/main.go) | Complete data for one planet | Intermediate |
| [`all_astro_data`](examples/all_astro_data/main.go) | Aggregate planet, house, zodiac data | Intermediate |
| [`bhava_chart_data`](examples/bhava_chart_data/main.go) | House chart | Intermediate |
| [`code_from_api_builder`](examples/code_from_api_builder/main.go) | API Builder-style call | Intermediate |
| [`rag_vedic_books`](examples/rag_vedic_books/main.go) | Available books and source-text search | Intermediate |
| [`error_handling`](examples/error_handling/main.go) | Cancellation and ordinary error handling | Intermediate |
| [`all_astro_data_json_output`](examples/all_astro_data_json_output/main.go) | Export JSON to `astro_data.json` | Intermediate |
| [`all_astro_data_csv`](examples/all_astro_data_csv/main.go) | Export planet results to CSV | Intermediate |
| [`batch_processing`](examples/batch_processing/main.go) | Multiple birth charts | Advanced |

Examples call the live API. Except for the single-call quick start, the shared runner reads `VEDASTRO_API_KEY`, handles Ctrl+C, and spaces free-tier requests 13 seconds apart. Export examples create their files exclusively and fail if the file already exists — rename old exports before re-running.

---

## 🔢 Decoded results

The API's `Payload` is decoded into standard Go JSON values:

| JSON | Go |
| --- | --- |
| Object | `map[string]any` |
| Array | `[]any` |
| Number | `float64` |
| String / Boolean / null | `string` / `bool` / `nil` |
| SVG response | `string` |

As in Python, a payload object with exactly one property is unwrapped to that property's value. Zero, `false`, empty strings, empty collections, and `nil` are all valid results, so never treat a zero value as "missing" — check the error first, then assert the type.

```go
result, err := client.PlanetNirayanaLongitude(ctx, vedastro.PlanetSun, birth)
if err != nil {
	log.Fatal(err)
}
longitude, ok := result.(float64) // numbers decode as float64
if !ok {
	log.Fatalf("unexpected type %T", result)
}
fmt.Println(longitude)
```

---

## ⚙️ Optional arguments

Methods with optional parameters accept **zero or one** method-specific options struct. Nil fields are omitted, letting the server apply its defaults. Use `vedastro.Ptr` to send explicit zero, `false`, or empty values:

```go
result, err := client.SearchSourceText(ctx, "Saturn",
	vedastro.SearchSourceTextOptions{
		TopK:       vedastro.Ptr(3),
		SourceName: vedastro.Ptr("Brihat-Parashara-Hora-Shastra"),
	})
```

Passing more than one options value returns an error **without sending a request**. Defaults are documented on generated fields and in [api_manifest.json](api_manifest.json).

---

## 🔑 Keys, contexts, and configuration

Use your own key via the environment:

```go
client := vedastro.NewClient(os.Getenv("VEDASTRO_API_KEY"))
```

An **empty key** leaves authentication to the service's free tier. Never embed paid keys in distributed source — Go binaries are trivially string-dumped.

Each client owns its settings and is **safe for concurrent requests**. Options are applied at construction:

```go
client := vedastro.NewClient(apiKey,
	vedastro.WithDefaultAyanamsa(vedastro.AyanamsaLahiri),
	vedastro.WithBaseURL("https://vedastro.zaishi.net/api/Calculate"),
	vedastro.WithHTTPClient(customClient),
)
```

| Option | Effect |
|--------|--------|
| `WithDefaultAyanamsa(v)` | Ayanamsa sent with this client's requests |
| `WithBaseURL(url)` | Point at a different Calculate service or a test server |
| `WithHTTPClient(c)` | Custom `*http.Client`; a shallow copy is retained |
| `WithTimeout(d)` | Optional calculation deadline; **none is applied by default** (see below) |

Invalid configuration is **reported by the first calculator call**, not by `NewClient` — so a bad base URL surfaces as an error from your request rather than a panic at startup.

### ⏳ There is no request deadline by default

Deliberately. A VedAstro calculation can take milliseconds or minutes depending on the endpoint and
the load on the service, and a client library has no way to know what is acceptable for your
workload. A built-in deadline would be exactly the kind of brittle logic that silently truncates a
valid answer, so **no deadline is applied** unless you ask for one. If a call is still running, it
is still working.

`go doc`'s `http.Client` has no timeout either, which is why the client only wraps your context
when you supply one:

```go
// A deadline, because *you* decided this call has run too long.
client := vedastro.NewClient(apiKey, vedastro.WithTimeout(30*time.Second))
```

Cancellation still works exactly as Go expects — pass a context with a deadline, or cancel it, and
the call stops:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
```

The difference is who decides: `WithTimeout` is the client's own deadline, while a context deadline
is yours. Both are honoured if set, and whichever expires first wins.

---

## 🧯 Errors and cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := client.PlanetRasiD1Sign(ctx, vedastro.PlanetSun, birth)
switch {
case err == nil:
	// use result
case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
	// caller cancelled, or the deadline expired
default:
	// remote HTTP / API failure arrives as *vedastro.APIError
	var apiErr *vedastro.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API error on %s: HTTP %d\n", apiErr.Endpoint, apiErr.StatusCode)
	} else {
		fmt.Printf("transport failure: %v\n", err)
	}
}
```

Remote HTTP/API and malformed-response failures use `*vedastro.APIError` (with `Endpoint`, `StatusCode`, and `Message`), retrievable via `errors.As`. Transport errors omit request credentials, and API messages redact the configured key. **The SDK never retries automatically** — add your own backoff if you need it.

---

## ⏱️ Pacing requests and rate limits

The free tier allows **five requests per minute**, counted across all clients in all processes on your key. Space requests at least **13 seconds** apart. The SDK does not throttle or retry for you.

```go
// A minimal pacer: at most one call every 13 seconds.
type pacer struct {
	mu   sync.Mutex
	last time.Time
	gap  time.Duration
}

func (p *pacer) wait(ctx context.Context) error {
	p.mu.Lock()
	wait := p.gap - time.Since(p.last)
	p.mu.Unlock()
	if wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	p.mu.Lock()
	p.last = time.Now()
	p.mu.Unlock()
	return nil
}
```

See [`examples/batch_processing`](examples/batch_processing/main.go) for a complete multi-chart version. Paid keys use the service's own (higher) limits.

---

## 🔄 Updates

`NewClient` prints the familiar one-line banner **once per process**. There is no network update check at startup. To check explicitly:

```go
info, err := client.CheckForUpdate(ctx)
if err == nil && info.Available {
	fmt.Println(info.UpdateCommand)          // e.g. go get ...@v1.0.1
	fmt.Println(info.RebuildInstructions)
}
```

This queries Go's public module proxy, never sends your API key, and never modifies files. Update your dependency, rebuild, and restart — a running binary keeps its compiled version. To silence the banner entirely, see [FAQ](FAQ.md).

---

## 🐛 Troubleshooting

### Q: "`interface conversion: interface {} is map[string]interface {}, not string`"

**A:** You asserted the wrong type. Result shapes vary by calculator.

```go
// ❌ Object payloads are not strings — this panics at runtime
// name := result.(string)

// ✅ Assert the object, then read the field
object, ok := result.(map[string]any)
if !ok {
	log.Fatalf("unexpected type %T", result)
}
name, _ := object["Name"].(string)
```

`HouseSignName` returns an unwrapped **string**, while `PlanetRasiD1Sign` returns a **`map[string]any`** — a common source of this panic. Use `go doc` to check a signature, or a two-value assertion while exploring.

### Q: "vedastro: time must use HH:MM DD/MM/YYYY +/-HH:MM"

**A:** The time string failed validation. Use 24-hour time, DD/MM/YYYY order, and an explicit numeric UTC offset:

```go
vedastro.NewTime("14:30 25/10/1992 +05:30", mumbai) // ✅
vedastro.NewTime("2:30 PM 25/10/1992 +05:30", mumbai) // ❌ AM/PM
vedastro.NewTime("14:30 1992-10-25 +05:30", mumbai)   // ❌ ISO date
vedastro.NewTime("14:30 25/10/1992", mumbai)          // ❌ no offset
```

### Q: "vedastro: a location name and valid longitude/latitude are required"

**A:** Longitude must be within ±180 and latitude within ±90, and the name must be non-empty. Check the argument order — `NewGeoLocation(name, longitude, latitude)`.

### Q: "I'm hitting rate limits"

**A:** The free tier allows five requests per minute across all clients. Pace calls 13 seconds apart (see [Pacing requests](#-pacing-requests-and-rate-limits)) or use a paid key. The error arrives as an `*APIError` whose `Message` contains the server's explanation.

### Q: "Which ayanamsa should I use?"

**A:** `vedastro.AyanamsaLahiri` for Indian astrology, `vedastro.AyanamsaKrishnamurti` for KP, `vedastro.AyanamsaFaganBradley` for Western sidereal. Omit it entirely to accept the server default, which measures as Lahiri.

### Q: "Can I use this in a server handling many requests?"

**A:** Yes. A `*Client` is safe for concurrent use and owns its settings, so one shared client is enough. Use `vedastro.WithAyanamsa(ctx, …)` for per-request ayanamsa overrides — it does not mutate shared state. Free-tier limits are per key, not per goroutine.

### Q: "Can I use this commercially?"

**A:** **Yes.** MIT licensed, and both the free and paid tiers permit commercial use.

### Q: "What's included in free vs premium?"

**A:** **All 684 calculations and all 47 ayanamsa systems in both.** Only throughput differs:

| Feature | Free | Premium |
|---------|------|---------|
| All calculations | ✅ | ✅ |
| All 47 ayanamsa | ✅ | ✅ |
| Commercial use | ✅ | ✅ |
| Swiss Ephemeris | ✅ | ✅ |
| Rate limit | 5 req/min | Higher (service-set) |
| Cost | $0/month | $1/month |

> Premium is quoted as **200 calls/minute** by the API's own rate-limit error. Pricing and limits are set by the service — confirm current numbers at [vedastro.org/API.html](https://vedastro.org/API.html).

### Q: "How do I get a premium API key?"

**A:**

1. Go to [vedastro.org/API.html](https://vedastro.org/API.html)
2. Choose a plan: $1/month or ₹758/year (India)
3. Pay via card, UPI, Google Pay, or PayPal
4. Get your key instantly from [vedastro.org/Account.html](https://vedastro.org/Account.html)

```go
client := vedastro.NewClient(os.Getenv("VEDASTRO_API_KEY"))
```

---

## 🚀 Why VedAstro? The Simplest & Most Affordable Vedic Astrology API

### ✨ Unbeatable Value

| What You Get | VedAstro | Competitors |
|--------------|----------|-------------|
| **Monthly Cost** | **$1/month** | $50-$200/month |
| **Free Tier** | ✅ 5 req/min | ❌ None or very limited |
| **Calculations** | **684 methods** | 50-200 methods |
| **Ayanamsa Systems** | **47 systems** | 3-10 systems |
| **Setup Complexity** | **Zero setup** | Complex (DLLs, ephemeris files) |
| **Commercial Use** | ✅ Both tiers | ❌ Enterprise only |

**The Bottom Line:** Get 10x more features at 1/50th the price. No credit card needed to start.

### 💰 Pricing That Makes Sense

| Tier | Price | Rate Limit | Best For |
|------|-------|------------|----------|
| **Free** | $0/month | 5 req/min | Learning, testing, personal projects |
| **Premium** | **$1/month** | Higher (service-set) | Production apps, commercial use |

**Indian Developers:** ₹79/month or ₹758/year (₹63/month, most popular)

> **All 684 calculations included in both tiers.** The only difference is throughput.

---

## 📖 What Can You Calculate? (684 Methods)

**Full API reference:** [vedastro.org/API.html](https://vedastro.org/API.html) · [pkg.go.dev](https://pkg.go.dev/github.com/VedAstro/VedAstro.Go) · `go doc`

Every calculator is a method on `*Client`, takes `context.Context` first, and returns `(any, error)`.

### Core chart analysis (225 methods)
`PlanetRasiD1Sign`, `PlanetNirayanaLongitude`, `PlanetConstellation`, `PlanetsInSign`, `PlanetsInConjunction`, `IsPlanetRetrograde`, `IsPlanetExalted`, `IsPlanetDebilitated`, `HouseSignName`, `HouseRasiSign`, `LordOfHouse`, `AllPlanetData`, `AllHouseData`, `AllZodiacSignData`, +211 more

### Divisional charts — vargas (102 methods)
`AllHouseNavamshaSign` (D9), `AllHouseDrekkanaSign` (D3), `AllHouseChaturthamsaSign` (D4), `PlanetShashtyamshaD60Sign` (D60), `PlanetDivisionalLongitude`, D1-D60 vargas

### Panchanga, muhurta & time (65 methods)
`LunarDay`, `TithiNumber`, `NakshatraPada`, `NithyaYoga`, `Karana`, `RahuKala`, `GulikaKala`, `Durmuhurta`, `SunriseTime`, `SunsetTime`, `HoraTable`, `AyanamsaDegree`, `PlanetEphemerisLongitude`, +51 more

### Strength, dignity & Shadbala (62 methods)
`PlanetShadbalaPinda`, `PlanetStrength`, `HouseStrength`, `PlanetIshtaScore`, `PlanetKashtaScore`, `AllPlanetOrderedByStrength`, `PickOutStrongestPlanet`, `PlanetDignity`, `IsPlanetVargottama`, +53 more

### Charts, aspects & events (55 methods)
`NorthIndianChart`, `SouthIndianChart`, `SkyChart`, `PlanetAspectDegree`, `PlanetsAspectingPlanet`, `IsPlanetAspectedByPlanet`, `EventsAtTime`, `EventsAtRange`, `EventStartTime`, `HoroscopePredictions`, `SwissEphemeris`, +44 more

### Prashna, chakra & Pancha Pakshi (28 methods)
`Chapter5PrashnaMargaPredictions` … `Chapter30PrashnaMargaPredictions`, `SarvatobhadraChakra`, `KotaChakra`, `SudarsanaChakra`, `BirthYamaPanchaPakshi`, `CalculateAshtamangalaNumberFromShells`, +few more

### Ashtakvarga (20 methods)
`SarvashtakavargaChart`, `BhinnashtakavargaChart`, `PlanetAshtakvargaBindu`, `GocharaKakshas`, `AshtakavargaLongevity`, `PrastaraAshtakavarga`, `SodyaAshtakavarga`, +13 more

### Dasa — planetary periods (20 methods)
`DasaAtRange`, `DasaAtTime`, `DasaForNow`, `DasaForLife`, `MoolaDasa`, `NarayanaDasa`, `KalachakraDasa`, `TithiAshtottariDasa`, +12 more

### Earthquake research (20 methods)
`CalculateEarthquakeRiskScore`, `IsEarthquakeNearEclipse`, `IsEarthquakeJupiterSaturnConjunction`, `IsEarthquakePlanetsClusteredInNarrowArc`, +16 more

### Longevity & life events (18 methods)
`DetailedAshtakavargaLongevity`, `MarriageByJupiter`, `ChildBirthByJupiter1`, `NativeDeathBySaturn`, `HasBalarishtaExceptions`, `MarakaPlanetList`, +12 more

### Upagraha & special points (15 methods)
`GulikaLongitude`, `MaandiLongitude`, `DhumaLongitude`, `UpaketuLongitude`, `KaalaLongitude`, `FortunaPoint`, `DestinyPoint`, `IsUpagraha`, +7 more

### Transit & timing — gochara (13 methods)
`PlanetSignTransit`, `TransitHouseFromLagna`, `TransitHouseFromMoon`, `IsGocharaOccurring`, `GetConstellationTransitStartTime`, `IsPlanetRetrograde`, +7 more

### Yogas, doshas & kartari (12 methods)
`KalaSarpaYoga`, `JHoraYogaList`, `IsPlanetInGandanta`, `KujaDosaScore`, `ClassifyForKartari`, `ShubKartariPlanets`, `PaapaKartariPlanets`, +5 more

### Jaimini & Tajika (8 methods)
`JaiminiRasiDrishti`, `JaiminiRasiStrength`, `TajakaVarshaphala`, `TajakaYogaList`, `TrueSiderealSolarReturn`, `ArudhaLagnaSign`, +2 more

### AI, ML & text search (7 methods)
`FindBirthTimeByMachineLearning`, `FindBirthTimeByMachineLearningTopK`, `FindBirthTimeByAnimal`, `SearchSourceText` (RAG over classical texts), `GetAvailableSourceTexts`, `HoroscopePredictionsForLargeAstrologyModelTrainingData`, +1 more

### Numerology (6 methods)
`BirthNumber`, `DestinyNumber`, `NameNumber`, `NameNumberPrediction`, `RootNumberFriendship`, `MainActivity`

### Compatibility & matching (4 methods)
`MatchReport`, `Tarabala`, `Chandrabala`, `YoniKutaAnimal`

### Health & nature scores (4 methods)
`PredictMedicalHealthConditions`, `HouseNatureScore`, `PlanetNatureScore`, `GetActiveNccBodyRulesAtTime`

---

## 🎯 Next Steps

1. **Install**: `go get github.com/VedAstro/VedAstro.Go@v1.0.0` (10 seconds)
2. **Try the examples above**: copy, paste, run! (5 minutes)
3. **Explore demos**: `go run ./examples/birth_chart_basics` (30 minutes)
4. **Read the docs**: [pkg.go.dev](https://pkg.go.dev/github.com/VedAstro/VedAstro.Go) or `go doc github.com/VedAstro/VedAstro.Go`
5. **Build something**: your first horoscope service! (1-2 hours)
6. **Upgrade when ready**: [vedastro.org/API.html](https://vedastro.org/API.html)

---

## 💡 Why Developers Love VedAstro

> "I was paying $150/month for a competing API. VedAstro is $1/month with more features and better docs. Absolute no-brainer." — Rahul, India

> "Setup took 2 minutes. First calculation worked immediately. No configuration hell. This is how all APIs should be." — Sarah, USA

> "684 calculations, 47 ayanamsas, Swiss Ephemeris accuracy, $1/month. I thought there was a catch. There isn't." — Yuki, Japan

> "The free tier is generous enough for my personal app with 50 users. When I scale up, $1/month won't break the bank." — Carlos, Brazil

---

## 🏗️ How It Works (Architecture)

```
Your Go Code
      v
vedastro package (this module, stdlib only)
      v
REST API (vedastro.zaishi.net)
      v
VedAstro Engine (Azure Cloud)
      v
Swiss Ephemeris (NASA JPL data)
```

The SDK is a thin, typed client: it builds the request body, POSTs it, unwraps the standard `{ Status, Payload }` envelope, and returns the inner payload. Calculations themselves run on VedAstro's servers.

**Why this design?**
- ✅ Zero third-party dependencies — standard library only
- ✅ Instant updates (684 calculations, always latest, no SDK release needed)
- ✅ Single static binary, no `cgo`, no ephemeris downloads
- ✅ Cross-platform (Windows/Mac/Linux)

---

## 🛠️ Development and generation

`calculate_generated.go`, `enums_generated.go`, and `api_manifest.json` come from **StaticTableGenerator Task 11** in the [main VedAstro repository](https://github.com/VedAstro/VedAstro). **Do not edit generated files.** Coverage is derived from metadata, with first-overload selection matching the other SDKs.

```sh
# From the main repository
dotnet run --project StaticTableGenerator -- --task 11

# From this repository
gofmt -w .
go vet ./...
go build ./...
go test -timeout 90s ./...
go test -race -timeout 120s ./...
python scripts/check_python_parity.py ../VedAstro.Python-PUBLIC/vedastro/calculate.py
```

Task 11 output defaults to the sibling `VedAstro.Go-PUBLIC` folder; set `VEDASTRO_GO_OUT` to override it, and optionally `VEDASTRO_GOFMT` to the gofmt executable. It uses local metadata only (no API, Cosmos, or LLM calls) and writes deterministic, formatted output.

Race checks require a supported platform and a C toolchain. CI tests Go 1.26 and 1.27 on Windows and Linux, including Linux race checks. Tests use local HTTP servers; **examples call the live API only when explicitly run**.

Release with a new semantic version tag after validation. Never move or replace an existing tag — Go's proxy caches releases permanently. See [Go's publishing guide](https://go.dev/doc/modules/publishing).

---

## 🤝 Contributing

Contributions are welcome — issues and PRs at [VedAstro.Go](https://github.com/VedAstro/VedAstro.Go).

> **Note:** everything named `*_generated.go` is auto-generated. Do not edit it directly; regenerate with Task 11 and re-run the checks above. `scripts/check_python_parity.py` guards name and argument-order parity with the Python SDK.

---

## 📄 License

MIT License - Use freely in commercial and personal projects.

---

## 🙏 Support the Project

VedAstro is non-profit and user-funded. If it saves you time and money:

- ⭐ **Star on GitHub** (helps others discover us)
- 💰 **Subscribe $1/month** at [vedastro.org/API.html](https://vedastro.org/API.html)
- 🎁 **Donate** at [vedastro.org/Donate](https://vedastro.org/Donate)
- 📢 **Share** with other developers

Every subscription helps keep VedAstro free and open-source! 🙏

---

## 📚 Additional Resources

- 📖 **Full API Docs**: [vedastro.org/API.html](https://vedastro.org/API.html)
- 📘 **Go package reference**: [pkg.go.dev/github.com/VedAstro/VedAstro.Go](https://pkg.go.dev/github.com/VedAstro/VedAstro.Go)
- 🚀 **Quick Start**: [QUICKSTART.md](QUICKSTART.md)
- ❓ **FAQ**: [FAQ.md](FAQ.md)
- 🧪 **Runnable examples**: [examples/README.md](examples/README.md)
- 🐍 **Python SDK**: [VedAstro.Python](https://github.com/VedAstro/VedAstro.Python)
- 🟨 **Node.js SDK**: [VedAstro.NodeJS](https://github.com/VedAstro/VedAstro.NodeJS)
- 💬 **Telegram**: [t.me/vedastro_org](https://t.me/vedastro_org)
- 🐛 **Issues**: [GitHub Issues](https://github.com/VedAstro/VedAstro.Go/issues)
- 🌐 **Website**: [vedastro.org](https://vedastro.org)

---

<p align="center">
  <strong>Made with ❤️ by users, for users</strong><br>
  <a href="https://vedastro.org">Website</a> •
  <a href="https://vedastro.org/API.html">API Docs</a> •
  <a href="https://github.com/VedAstro/VedAstro">GitHub</a> •
  <a href="https://t.me/vedastro_org">Telegram</a> •
  <a href="https://vedastro.org/Donate">Donate</a>
</p>

<p align="center">
  <em>🪐 Empowering developers to build amazing astrology apps since 2020</em>
</p>
