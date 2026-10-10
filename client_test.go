package vedastro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func testBirth(t *testing.T) Time {
	t.Helper()
	value, err := NewTime("14:30 25/10/1992 +05:30", NewGeoLocation("Mumbai", 72.8777, 19.076))
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestPythonWireFixture(t *testing.T) {
	var received map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/PlanetRasiD1Sign" {
			t.Errorf("request: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing JSON content type")
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Error(err)
		}
		fmt.Fprint(w, `{"Status":"Pass","Payload":{"PlanetRasiD1Sign":{"Name":"Libra","DegreesIn":8.5}}}`)
	}))
	defer server.Close()
	result, err := NewClient("test-secret", WithBaseURL(server.URL), WithDefaultAyanamsa(AyanamsaLahiri)).
		PlanetRasiD1Sign(context.Background(), PlanetSun, testBirth(t))
	if err != nil {
		t.Fatal(err)
	}
	if result.(map[string]any)["Name"] != "Libra" {
		t.Fatalf("result: %#v", result)
	}
	// Captured independently from Python's GeoLocation.to_json / Time.to_json
	// and the generated PlanetRasiD1Sign POST contract.
	var expected map[string]any
	if err := json.Unmarshal([]byte(`{"planetName":"Sun","time":{"StdTime":"14:30 25/10/1992 +05:30","Location":{"Name":"Mumbai","Longitude":72.8777,"Latitude":19.076}},"APIKey":"test-secret","Ayanamsa":"LAHIRI"}`), &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(received, expected) {
		t.Fatalf("body: %#v", received)
	}
}

func TestResponseShapes(t *testing.T) {
	cases := []struct {
		name, body, contentType string
		status                  int
		want                    any
		failure                 string
	}{
		{"zero", `{"Payload":0}`, "", 200, float64(0), ""},
		{"false", `{"Payload":false}`, "", 200, false, ""},
		{"empty string", `{"Payload":""}`, "", 200, "", ""},
		{"null", `{"Payload":null}`, "", 200, nil, ""},
		{"empty list", `{"Payload":[]}`, "", 200, []any{}, ""},
		{"empty object", `{"Payload":{}}`, "", 200, map[string]any{}, ""},
		{"list", `{"Payload":[1,"Sun"]}`, "", 200, []any{float64(1), "Sun"}, ""},
		{"single key", `{"Payload":{"Method":false}}`, "", 200, false, ""},
		{"several keys", `{"Payload":{"a":0,"b":false}}`, "", 200, map[string]any{"a": float64(0), "b": false}, ""},
		{"SVG", `<?xml version="1.0"?><svg/>`, "image/svg+xml; charset=utf-8", 200, `<?xml version="1.0"?><svg/>`, ""},
		{"missing", `{"Status":"Pass"}`, "", 200, nil, "missing"},
		{"invalid JSON", "<html>upstream error</html>", "", 200, nil, "invalid JSON"},
		{"invalid envelope", "[]", "", 200, nil, "invalid JSON"},
		{"failure", `{"Status":"Fail","Payload":"bad key test-secret"}`, "", 200, nil, "bad key [redacted]"},
		{"structured failure", `{"Status":"Fail","Payload":{"APIKey":"test-secret"}}`, "", 200, nil, "API reported failure"},
		{"rate limit", "test-secret", "", 429, nil, "HTTP request failed"},
		{"server failure", "test-secret", "", 503, nil, "HTTP request failed"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.contentType != "" {
					w.Header().Set("Content-Type", tc.contentType)
				}
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.body)
			}))
			defer server.Close()
			result, err := NewClient("test-secret", WithBaseURL(server.URL)).GetAvailableSourceTexts(context.Background())
			if tc.failure != "" {
				if err == nil || !strings.Contains(err.Error(), tc.failure) {
					t.Fatalf("error: %v", err)
				}
				if strings.Contains(err.Error(), "test-secret") {
					t.Fatal("credential leaked")
				}
				var apiError *APIError
				if !errors.As(err, &apiError) || apiError.StatusCode != tc.status {
					t.Fatalf("typed error: %v", err)
				}
			} else if err != nil || !reflect.DeepEqual(result, tc.want) {
				t.Fatalf("got %#v, %v; want %#v", result, err, tc.want)
			}
		})
	}
}

func TestOptionalArguments(t *testing.T) {
	var bodies []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		bodies = append(bodies, body)
		fmt.Fprint(w, `{"Payload":true}`)
	}))
	defer server.Close()
	c := NewClient("", WithBaseURL(server.URL))
	ctx := context.Background()
	if _, err := c.SearchSourceText(ctx, "Saturn"); err != nil {
		t.Fatal(err)
	}
	if _, exists := bodies[0]["topK"]; exists {
		t.Fatal("omitted default was sent")
	}
	if _, exists := bodies[0]["APIKey"]; exists {
		t.Fatal("empty key was sent")
	}
	if _, err := c.SearchSourceText(ctx, "Saturn", SearchSourceTextOptions{TopK: Ptr(0), SourceName: Ptr("")}); err != nil {
		t.Fatal(err)
	}
	if bodies[1]["topK"] != float64(0) || bodies[1]["sourceName"] != "" {
		t.Fatal(bodies[1])
	}
	if _, err := c.LunarMonth(ctx, testBirth(t), LunarMonthOptions{IgnoreLeapMonth: Ptr(false)}); err != nil {
		t.Fatal(err)
	}
	if value, ok := bodies[2]["ignoreLeapMonth"]; !ok || value != false {
		t.Fatal(bodies[2])
	}
	if _, err := c.SearchSourceText(ctx, "Saturn", SearchSourceTextOptions{}, SearchSourceTextOptions{}); err == nil {
		t.Fatal("accepted several options")
	}
	if len(bodies) != 3 {
		t.Fatal("invalid options sent a request")
	}
}

func TestConcurrentClientsAndContexts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		switch body["APIKey"] {
		case "client-a":
			if body["Ayanamsa"] != "LAHIRI" {
				t.Errorf("client A scope leaked: %#v", body)
			}
		case "client-b":
			if _, exists := body["Ayanamsa"]; exists {
				t.Errorf("explicit server default lost: %#v", body)
			}
		default:
			t.Errorf("unknown key: %#v", body)
		}
		fmt.Fprint(w, `{"Payload":true}`)
	}))
	defer server.Close()
	a := NewClient("client-a", WithBaseURL(server.URL), WithDefaultAyanamsa(AyanamsaRaman))
	b := NewClient("client-b", WithBaseURL(server.URL), WithDefaultAyanamsa(AyanamsaKrishnamurti))
	parent := WithAyanamsa(context.Background(), AyanamsaRaman)
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			if _, err := a.GetAvailableSourceTexts(WithAyanamsa(parent, AyanamsaLahiri)); err != nil {
				t.Error(err)
			}
		}()
		go func() {
			defer wg.Done()
			if _, err := b.GetAvailableSourceTexts(WithAyanamsa(parent, "")); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if parent.Value(ayanamsaKey{}) != AyanamsaRaman {
		t.Fatal("parent context changed")
	}
	if a.ayanamsa != AyanamsaRaman || b.ayanamsa != AyanamsaKrishnamurti {
		t.Fatal("client defaults changed")
	}
}

func TestCancellationAndTimeout(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	defer server.Close()
	defer close(release)
	c := NewClient("", WithBaseURL(server.URL), WithTimeout(30*time.Millisecond))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.GetAvailableSourceTexts(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel: %v", err)
	}
	if _, err := c.GetAvailableSourceTexts(context.Background()); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("timeout: %v", err)
	}
}

func TestInvalidConfiguration(t *testing.T) {
	for _, option := range []ClientOption{WithBaseURL("not a URL"), WithBaseURL("https://user:secret@example.com"),
		WithTimeout(0), WithHTTPClient(nil), WithDefaultAyanamsa("unknown")} {
		if _, err := NewClient("secret", option).GetAvailableSourceTexts(context.Background()); err == nil {
			t.Fatal("invalid configuration accepted")
		}
	}
	c := NewClient("")
	if _, err := c.GetAvailableSourceTexts(WithAyanamsa(context.Background(), "invalid")); err == nil {
		t.Fatal("invalid scope accepted")
	}
	if _, err := c.GetAvailableSourceTexts(nil); err == nil {
		t.Fatal("nil context accepted")
	}
}

func TestNoDeadlineByDefault(t *testing.T) {
	// A calculation can take milliseconds or minutes, so the client must not impose a deadline
	// of its own. A server that answers slowly must therefore still succeed with no WithTimeout.
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Status":"Pass","Payload":"ok"}`))
	}))
	defer server.Close()

	slowButFine := make(chan error, 1)
	go func() {
		// No WithTimeout: the call should wait for the server rather than give up.
		_, err := NewClient("", WithBaseURL(server.URL)).GetAvailableSourceTexts(context.Background())
		slowButFine <- err
	}()

	select {
	case err := <-slowButFine:
		t.Fatalf("call returned before the server answered: %v", err)
	case <-time.After(250 * time.Millisecond):
		// Still waiting, as it should be: nothing cut it short.
	}

	close(release)
	if err := <-slowButFine; err != nil {
		t.Fatalf("slow call failed without a deadline: %v", err)
	}
}

func TestGeneratedMethodSurface(t *testing.T) {
	data, err := os.ReadFile("api_manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		MethodCount int
		Methods     []struct {
			Name       string
			Parameters []struct{ Optional bool }
		}
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.MethodCount != len(manifest.Methods) {
		t.Fatal("manifest count mismatch")
	}
	clientType := reflect.TypeOf((*Client)(nil))
	seen := map[string]bool{}
	for _, entry := range manifest.Methods {
		if seen[entry.Name] {
			t.Fatalf("duplicate: %s", entry.Name)
		}
		seen[entry.Name] = true
		method, exists := clientType.MethodByName(entry.Name)
		if !exists {
			t.Fatalf("missing method %s", entry.Name)
		}
		wantInputs, optional := 2, false // receiver and context
		for _, p := range entry.Parameters {
			if p.Optional {
				optional = true
			} else {
				wantInputs++
			}
		}
		if optional {
			wantInputs++
		}
		if method.Type.NumIn() != wantInputs || method.Type.IsVariadic() != optional ||
			method.Type.NumOut() != 2 || method.Type.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			t.Errorf("signature mismatch: %s", entry.Name)
		}
	}
}
