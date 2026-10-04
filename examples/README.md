# Runnable examples

Run these from the repository root with Go 1.26 or later:

```sh
go run ./examples/quick_start
go run ./examples/daily_panchanga
```

Examples call the live API. Except for the single-call quick start, the shared runner reads `VEDASTRO_API_KEY`, handles Ctrl+C, and spaces free-tier requests 13 seconds apart. Pacing applies within one process; leave a gap between separate runs. Paid keys use the service's own limits.

| Directory | Python demo topic |
| --- | --- |
| quick_start | First calculation and decoded result |
| custom_ayanamsa | Lahiri, Raman, and Krishnamurti context overrides |
| batch_processing | Multiple birth charts |
| all_astro_data | Aggregate planet, house, and zodiac data |
| all_astro_data_json_output | Export JSON to astro_data.json |
| all_astro_data_csv | Export planet results to astro_data.csv |
| all_planet_data | Complete data for one planet |
| bhava_chart_data | House chart |
| birth_chart_basics | Signs, ascendant, and birth chart basics |
| code_from_api_builder | API Builder-style call |
| current_planets | Positions at the current time |
| daily_panchanga | Daily calendar calculations |
| divisional_charts | Divisional signs |
| error_handling | Cancellation and ordinary error handling |
| horoscope_prediction_names | Prediction names |
| marriage_compatibility | Compatibility calculation |
| match_checker | Match report |
| numerology_calculator | Name and birth numerology |
| rag_vedic_books | Available books and source-text search |
| transit_analysis | Planetary transits |
| vimshottari_dasa | Dasa over a date range |

Export examples create their output files exclusively and return an error if the file already exists. Rename existing exports before running again. Large aggregate and dasa calculations can take longer than simple planet calls.

All sample times are explicit examples; change the birth date, time, UTC offset, and coordinates for your own inputs. The current-planets example uses the current clock time with Mumbai's UTC offset.
