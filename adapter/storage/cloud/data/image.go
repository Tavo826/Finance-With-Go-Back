package data

import "time"

type ErrorResponse struct {
	Message string `json:"message"`
}

type ImageResponse struct {
	ApiKey           string    `json:"api_key"`
	AssetFolder      string    `json:"asset_folder"`
	AssetId          string    `json:"asset_id"`
	Bytes            float64   `json:"bytes"`
	CreatedAt        time.Time `json:"created_at"`
	DisplayName      string    `json:"display_name"`
	Etag             string    `json:"etag"`
	Format           string    `json:"format"`
	Height           int64     `json:"height"`
	OriginalFilename string    `json:"original_filename"`
	Placeholder      bool      `json:"placeholder"`
	PublicId         string    `json:"public_id"`
	ResourceType     string    `json:"resource_type"`
	SecureUrl        string    `json:"secure_url"`
	Signature        string    `json:"signature"`
	Tags             []string  `json:"tags"`
	Type             string    `json:"type"`
	Url              string    `json:"url"`
	Version          float64   `json:"version"`
	VersionId        string    `json:"version_id"`
	Width            int64     `json:"width"`
}

type ImageUploadResponse struct {
	AssetId               string        `json:"asset_id"`
	PublicId              string        `json:"public_id"`
	AssetFolder           string        `json:"asset_folder"`
	DisplayName           string        `json:"display_name"`
	Version               float64       `json:"version"`
	VersionId             string        `json:"version_id"`
	Signature             string        `json:"signature"`
	Width                 int64         `json:"width"`
	Height                int64         `json:"height"`
	Format                string        `json:"format"`
	ResourceType          string        `json:"resource_type"`
	CreatedAt             time.Time     `json:"created_at"`
	Bytes                 float64       `json:"bytes"`
	Type                  string        `json:"type"`
	Etag                  string        `json:"etag"`
	Url                   string        `json:"url"`
	SecureUrl             string        `json:"secure_url"`
	AccessMode            string        `json:"access_mode"`
	Overwritten           bool          `json:"overwritten"`
	OriginalFilename      string        `json:"original_filename"`
	Eager                 string        `json:"eager"`
	ResponsiveBreakpoints string        `json:"responsive_breakpoints"`
	HookExecution         string        `json:"hook_execution"`
	Error                 ErrorResponse `json:"error"`
	Response              ImageResponse
}
