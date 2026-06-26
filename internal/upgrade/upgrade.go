package upgrade

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type releaseResponse struct {
	TagName string `json:"tag_name"`
}

// CheckLatest fetches the latest release tag from GitHub and reports whether it's newer.
func CheckLatest(client *http.Client, currentVersion string) (latest string, isNewer bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.github.com/repos/AestheticAutonomy/justctx/releases/latest", nil)
	if err != nil {
		return "", false, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "jctx/"+currentVersion)

	resp, err := client.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("checking for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", false, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	var rel releaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", false, fmt.Errorf("parsing release response: %w", err)
	}

	latest = rel.TagName
	isNewer = IsNewer(currentVersion, latest)
	return latest, isNewer, nil
}

// IsNewer reports whether candidate is a newer version than current.
// Strips leading 'v' and compares semver components numerically.
func IsNewer(current, candidate string) bool {
	cv := parseVersion(current)
	nv := parseVersion(candidate)
	for i := 0; i < 3; i++ {
		if nv[i] > cv[i] {
			return true
		}
		if nv[i] < cv[i] {
			return false
		}
	}
	return false
}

func parseVersion(v string) [3]int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	var out [3]int
	for i, p := range parts {
		if i >= 3 {
			break
		}
		n := 0
		for _, c := range p {
			if c >= '0' && c <= '9' {
				n = n*10 + int(c-'0')
			} else {
				break
			}
		}
		out[i] = n
	}
	return out
}
