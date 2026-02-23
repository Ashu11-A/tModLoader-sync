package api

import (
  "fmt"
  "io"
  "net/http"
  "os"
  "tml-sync/client/internal/models"

  "github.com/go-resty/resty/v2"
)

type Client struct {
  restClient *resty.Client
}

func New(host string, port int) *Client {
  baseURL := fmt.Sprintf("http://%s:%d", host, port)
  return &Client{
    restClient: resty.New().SetBaseURL(baseURL),
  }
}

// fetchData centraliza a lógica de requisição GET e decodificação JSON
func fetchData[T any](client *Client, endpoint string) (*T, error) {
  var data T
  response, err := client.restClient.R().SetResult(&data).Get(endpoint)
  if err != nil {
    return nil, err
  }
  
  if response.StatusCode() != http.StatusOK {
    return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode())
  }
  
  return &data, nil
}

func (client *Client) GetLanguage() (string, error) {
  response, err := client.restClient.R().Get("/v1/language")
  if err != nil {
    return "", err
  }
  return response.String(), nil
}

func (client *Client) GetVersion() (*models.ServerVersionResponse, error) {
  return fetchData[models.ServerVersionResponse](client, "/version")
}

func (client *Client) GetSyncStatus() (*models.SyncData, error) {
  return fetchData[models.SyncData](client, "/v1/sync")
}

func (client *Client) TriggerServerUpdate(version string) error {
	response, err := client.restClient.R().
		SetQueryParam("version", version).
		Get("/v1/update")
	if err != nil {
		return err
	}
	if response.StatusCode() != http.StatusOK {
		return fmt.Errorf("server error triggering update: %s", response.String())
	}
	return nil
}

type progressReader struct {
  io.Reader
  totalBytes int64
  bytesSent  int64
  onProgress func(int64, int64)
}

func (reader *progressReader) Read(buffer []byte) (int, error) {
  bytesRead, err := reader.Reader.Read(buffer)
  if bytesRead > 0 {
    reader.bytesSent += int64(bytesRead)
    if reader.onProgress != nil {
      reader.onProgress(reader.totalBytes, reader.bytesSent)
    }
  }
  return bytesRead, err
}

func (client *Client) UploadMod(filePath, name, version, hash string, onProgress func(int64, int64)) error {
  file, err := os.Open(filePath)
  if err != nil {
    return err
  }
  defer file.Close()

  fileInfo, err := file.Stat()
  if err != nil {
    return err
  }

  reader := &progressReader{
    Reader:     file,
    totalBytes: fileInfo.Size(),
    onProgress: onProgress,
  }

  response, err := client.restClient.R().
    SetFileReader("mod", fileInfo.Name(), reader).
    SetFormData(map[string]string{
      "name":    name,
      "version": version,
      "hash":    hash,
    }).
    Post("/v1/upload")

  if err != nil {
    return err
  }

  if response.StatusCode() != http.StatusOK {
    return fmt.Errorf("server error (%d): %s", response.StatusCode(), response.String())
  }

  return nil
}

// UploadEnabledJSON uploads the enabled.json file to the server.
func (client *Client) UploadEnabledJSON(filePath, hash string, onProgress func(int64, int64)) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open enabled.json: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for enabled.json: %w", err)
	}

	reader := &progressReader{
		Reader:     file,
		totalBytes: fileInfo.Size(),
		onProgress: onProgress,
	}

	response, err := client.restClient.R().
		SetFileReader("enabled", fileInfo.Name(), reader).
		SetFormData(map[string]string{
			"hash": hash,
		}).
		Post("/v1/enabled")

	if err != nil {
		return err
	}

		if response.StatusCode() != http.StatusOK {
			return fmt.Errorf("server error uploading enabled.json (%d): %s", response.StatusCode(), response.String())
		}
	
		return nil
	}
	
	func (client *Client) Stop() error {
		response, err := client.restClient.R().Post("/v1/stop")
		if err != nil {
			return err
		}
		if response.StatusCode() != http.StatusOK {
			return fmt.Errorf("server error stopping (%d): %s", response.StatusCode(), response.String())
		}
		return nil
	}
	