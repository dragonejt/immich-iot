package immich

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ImmichClient interface {
	GetServerStatistics() (ServerStatistics, error)
	GetStorage() error
}

type ImmichAPI struct {
	ImmichClient
	URL         string
	API_KEY     string
	HTTP_CLIENT *http.Client
}

type ServerStatistics struct {
	Photos      int   `json:"photos"`
	Videos      int   `json:"videos"`
	Usage       int64 `json:"usage"`
	UsagePhotos int64 `json:"usagePhotos"`
	UsageVideos int64 `json:"usageVideos"`
	UsageByUser []UsageByUser
}

type UsageByUser struct {
	UserID           string `json:"userId"`
	Username         string `json:"userName"`
	Photos           int    `json:"photos"`
	Videos           int    `json:"videos"`
	QuotaSizeInBytes int64  `json:"quotaSizeInBytes"`
	Usage            int64  `json:"usage"`
	UsagePhotos      int64  `json:"usagePhotos"`
	UsageVideos      int64  `json:"usageVideos"`
}

type Storage struct {
	DiskUse             string  `json:"diskUse"`
	DiskUseRaw          int64   `json:"diskUseRaw"`
	DiskUsagePercentage float64 `json:"diskUsagePercentage"`
	DiskAvailable       string  `json:"diskAvailable"`
	DiskAvailableRaw    int64   `json:"diskAvailableRaw"`
	DiskSize            string  `json:"diskSize"`
	DiskSizeRaw         int64   `json:"diskSizeRaw"`
}

func (immichAPI *ImmichAPI) GetServerStatistics() (ServerStatistics, error) {
	serverStatistics := ServerStatistics{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/server/statistics", immichAPI.URL), nil)
	if err != nil {
		return serverStatistics, fmt.Errorf("failed to create immich API request: %w", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", immichAPI.API_KEY))

	response, err := immichAPI.HTTP_CLIENT.Do(request)
	if err != nil {
		return serverStatistics, fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return serverStatistics, fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return serverStatistics, fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &serverStatistics)
	if err != nil {
		return serverStatistics, fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return serverStatistics, nil
}

func (immichAPI *ImmichAPI) GetStorage() (Storage, error) {
	storage := Storage{}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/server/storage", immichAPI.URL), nil)
	if err != nil {
		return storage, fmt.Errorf("failed to create immich API request: %w", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", immichAPI.API_KEY))

	response, err := immichAPI.HTTP_CLIENT.Do(request)
	if err != nil {
		return storage, fmt.Errorf("failed to execute API request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return storage, fmt.Errorf("failed status code API response: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return storage, fmt.Errorf("failed reading API response body: %d", response.StatusCode)
	}

	err = json.Unmarshal(body, &storage)
	if err != nil {
		return storage, fmt.Errorf("failed to unmarshal API response: %s", string(body))
	}

	return storage, nil
}
