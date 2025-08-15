package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"solid-octo-umbrella/models"
)

func GetNOAAActiveAlerts(region string) (*models.NOAAActiveAlertsResponse, error) {
	base, _ := url.Parse("https://api.weather.gov/alerts/active")
	q := base.Query()
	q.Set("region", region)
	base.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("User-Agent", "solid-octo-umbrella (courtney@elsner.dev)")
	req.Header.Set("Accept", "application/geo+json")

	client := &http.Client{Timeout: 15 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nws non-200: %s", res.Status)
	}

	var out models.NOAAActiveAlertsResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &out, nil
}
