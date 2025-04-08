package bucket

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

type BucketService interface {
	GetFile(ctx context.Context, path string) ([]byte, error)
}

type HTTPBucket struct {
	BaseURL string
	Client  *http.Client
}

func NewHTTPBucket() *HTTPBucket {
	baseURL := os.Getenv("BUCKET_BASE_URL")
	return &HTTPBucket{
		BaseURL: baseURL,
		Client:  http.DefaultClient,
	}
}

func (b *HTTPBucket) GetFile(ctx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", b.BaseURL, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição GET para o bucket: %w", err)
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição ao bucket: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arquivo não encontrado no bucket: %s (status %d)", path, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o corpo da resposta do bucket: %w", err)
	}

	return body, nil
}