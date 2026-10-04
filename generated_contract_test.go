package vedastro

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

type contractTransport func(*http.Request) (*http.Response, error)

func (f contractTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// Exercise every generated wrapper against its independently emitted metadata,
// catching swapped same-typed arguments, spelling mistakes, and lost options.
func TestAllGeneratedRequestContracts(t *testing.T) {
	data, err := os.ReadFile("api_manifest.json")
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Methods []struct {
			Name       string
			Parameters []struct {
				Name     string
				Optional bool
			}
		}
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	var received map[string]any
	var path string
	transport := contractTransport(func(r *http.Request) (*http.Response, error) {
		path = r.URL.Path
		received = nil
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			return nil, err
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{"Payload":true}`)), Request: r}, nil
	})
	client := NewClient("fixture-key", WithBaseURL("https://example.test"),
		WithHTTPClient(&http.Client{Transport: transport}))
	for _, entry := range manifest.Methods {
		t.Run(entry.Name, func(t *testing.T) {
			method := reflect.ValueOf(client).MethodByName(entry.Name)
			args := []reflect.Value{reflect.ValueOf(context.Background())}
			expected := map[string]any{"APIKey": "fixture-key"}
			requiredIndex := 1
			var options reflect.Value
			if method.Type().IsVariadic() {
				options = reflect.New(method.Type().In(method.Type().NumIn() - 1).Elem()).Elem()
			}
			optionalIndex := 0
			for index, parameter := range entry.Parameters {
				var value reflect.Value
				if parameter.Optional {
					field := options.Field(optionalIndex)
					optionalIndex++
					value = contractValue(t, field.Type().Elem(), index+1)
					pointer := reflect.New(field.Type().Elem())
					pointer.Elem().Set(value)
					field.Set(pointer)
				} else {
					value = contractValue(t, method.Type().In(requiredIndex), index+1)
					requiredIndex++
					args = append(args, value)
				}
				expected[parameter.Name] = value.Interface()
			}
			if options.IsValid() {
				args = append(args, options)
			}
			result := method.Call(args)
			if !result[1].IsNil() {
				t.Fatal(result[1].Interface())
			}
			if path != "/"+entry.Name {
				t.Fatalf("endpoint %s", path)
			}
			encoded, err := json.Marshal(expected)
			if err != nil {
				t.Fatal(err)
			}
			var normalized map[string]any
			if err := json.Unmarshal(encoded, &normalized); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(received, normalized) {
				t.Fatalf("request mismatch\nreceived: %#v\nexpected: %#v", received, normalized)
			}
		})
	}
}

func contractValue(t *testing.T, typ reflect.Type, seed int) reflect.Value {
	t.Helper()
	if typ == reflect.TypeOf(Time{}) {
		return reflect.ValueOf(testBirth(t))
	}
	if typ == reflect.TypeOf(GeoLocation{}) {
		return reflect.ValueOf(NewGeoLocation("fixture", float64(seed), 10))
	}
	value := reflect.New(typ).Elem()
	switch typ.Kind() {
	case reflect.String:
		value.SetString(fmt.Sprintf("argument-%d", seed))
	case reflect.Bool:
		value.SetBool(seed%2 == 0)
	case reflect.Int, reflect.Int64:
		value.SetInt(int64(seed))
	case reflect.Uint, reflect.Uint64:
		value.SetUint(uint64(seed))
	case reflect.Float64:
		value.SetFloat(float64(seed) + 0.25)
	case reflect.Interface:
		value.Set(reflect.ValueOf(map[string]any{"fixture": seed}))
	case reflect.Pointer:
		value.Set(reflect.New(typ.Elem()))
		value.Elem().Set(contractValue(t, typ.Elem(), seed))
	case reflect.Slice:
		value.Set(reflect.MakeSlice(typ, 1, 1))
		value.Index(0).Set(contractValue(t, typ.Elem(), seed))
	case reflect.Map:
		value.Set(reflect.MakeMap(typ))
		key := reflect.New(typ.Key()).Elem()
		key.SetString("fixture")
		value.SetMapIndex(key, contractValue(t, typ.Elem(), seed))
	default:
		t.Fatalf("unhandled fixture type %v", typ)
	}
	return value
}
