package vedastro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// UpdateInfo describes a module release. A new dependency takes effect only after rebuilding.
type UpdateInfo struct {
	CurrentVersion      string
	LatestVersion       string
	Available           bool
	UpdateCommand       string
	RebuildInstructions string
}

// CheckForUpdate explicitly checks the public Go module proxy with a five-second deadline.
// It never sends the API key, changes module files, or runs a package manager.
func (c *Client) CheckForUpdate(ctx context.Context) (UpdateInfo, error) {
	if c == nil || ctx == nil {
		return UpdateInfo{}, errors.New("vedastro: client and context are required")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.updateURL, nil)
	if err != nil {
		return UpdateInfo{}, errors.New("vedastro: cannot create update request")
	}
	response, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return UpdateInfo{}, fmt.Errorf("vedastro: update check interrupted: %w", ctx.Err())
		}
		return UpdateInfo{}, errors.New("vedastro: update check failed")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return UpdateInfo{}, fmt.Errorf("vedastro: module proxy returned HTTP %d", response.StatusCode)
	}
	var release struct{ Version string }
	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		return UpdateInfo{}, errors.New("vedastro: invalid module proxy response")
	}
	latest, err := releaseNumbers(release.Version)
	if err != nil {
		return UpdateInfo{}, err
	}
	current, err := releaseNumbers(Version)
	if err != nil {
		return UpdateInfo{}, err
	}
	available := false
	for index := range latest {
		if latest[index] != current[index] {
			available = latest[index] > current[index]
			break
		}
	}
	return UpdateInfo{
		CurrentVersion: Version, LatestVersion: release.Version, Available: available,
		UpdateCommand:       "go get github.com/VedAstro/VedAstro.Go@" + release.Version,
		RebuildInstructions: "Rebuild your application with go build and restart it to use the updated SDK.",
	}, nil
}

func releaseNumbers(version string) ([3]uint64, error) {
	var result [3]uint64
	invalid := errors.New("vedastro: module proxy did not return a stable semantic version")
	if !strings.HasPrefix(version, "v") {
		return result, invalid
	}
	parts := strings.Split(version[1:], ".")
	if len(parts) != 3 {
		return result, invalid
	}
	for i, part := range parts {
		if part == "" || (len(part) > 1 && part[0] == '0') {
			return result, invalid
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return result, invalid
			}
		}
		value, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return result, invalid
		}
		result[i] = value
	}
	return result, nil
}
