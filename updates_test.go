package vedastro

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateCheck(t *testing.T) {
	for _, tc := range []struct {
		version        string
		newer, invalid bool
	}{
		{"v1.0.0", false, false}, {"v1.0.1", true, false}, {"v1.10.0", true, false},
		{"v0.9.0", false, false}, {"bad", false, true}, {"v1.0.1;command", false, true},
	} {
		t.Run(tc.version, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet || r.Header.Get("Authorization") != "" || r.URL.RawQuery != "" {
					t.Error("update request contains authentication or is not GET")
				}
				fmt.Fprintf(w, `{"Version":%q}`, tc.version)
			}))
			defer server.Close()
			client := NewClient("private-key")
			client.updateURL = server.URL
			info, err := client.CheckForUpdate(context.Background())
			if tc.invalid {
				if err == nil {
					t.Fatal("invalid version accepted")
				}
				return
			}
			if err != nil || info.Available != tc.newer || info.CurrentVersion != Version ||
				!strings.Contains(info.UpdateCommand, tc.version) || info.RebuildInstructions == "" {
				t.Fatalf("update info: %#v %v", info, err)
			}
		})
	}
}
