package vedastro

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimeAndLocationHelpers(t *testing.T) {
	geo := NewGeoLocation("New York", -74.006, 40.7128)
	birth, err := NewTime("09:00 01/01/1990 -05:00", geo)
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(birth)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"StdTime":"09:00 01/01/1990 -05:00","Location":{"Name":"New York","Longitude":-74.006,"Latitude":40.7128}}`
	if string(data) != want {
		t.Fatalf("wire: %s", data)
	}
	fromDate, err := NewTimeFromDate(time.Date(1990, 1, 1, 9, 0, 0, 0, time.FixedZone("NY", -5*3600)), geo)
	if err != nil || fromDate != birth {
		t.Fatalf("date conversion: %v, %v", fromDate, err)
	}
	for _, input := range []string{"1990-01-01", "09:00 30/02/1990 -05:00", "09:00 01/01/1990", "25:00 01/01/1990 +00:00"} {
		if _, err := NewTime(input, geo); err == nil {
			t.Fatalf("accepted %q", input)
		}
	}
	if _, err := json.Marshal(NewGeoLocation("", 0, 0)); err == nil {
		t.Fatal("accepted unnamed location")
	}
	if _, err := json.Marshal(NewGeoLocation("invalid", 181, 0)); err == nil {
		t.Fatal("accepted invalid coordinate")
	}
}
