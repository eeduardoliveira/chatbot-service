package bucket

import (
	"context"
	"encoding/json"
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

func (b *HTTPBucket) GetFile(ctx context.Context, clienteID string) ([]byte, error) {
	bucketSignerURL := os.Getenv("BUCKET_SIGNER_URL")
	bucketName := os.Getenv("BUCKET_NAME")

	req, err := buildSignedURLRequest(ctx, bucketSignerURL, bucketName, clienteID)
	if err != nil {
		return nil, fmt.Errorf("erro ao montar requisição para bucket-signer: %w", err)
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição ao bucket-signer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro do bucket-signer: status %d", resp.StatusCode)
	}

	var result struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta do bucket-signer: %w", err)
	}

	// Requisição para buscar o arquivo via presigned URL
	fileReq, err := http.NewRequestWithContext(ctx, http.MethodGet, result.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição GET do arquivo: %w", err)
	}

	fileResp, err := b.Client.Do(fileReq)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar arquivo com URL assinada: %w", err)
	}
	defer fileResp.Body.Close()

	if fileResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("arquivo não encontrado no bucket: %s (status %d)", result.URL, fileResp.StatusCode)
	}

	body, err := io.ReadAll(fileResp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o corpo do arquivo: %w", err)
	}

	return body, nil
}

func buildSignedURLRequest(ctx context.Context, signerBaseURL, bucketName, clienteID string) (*http.Request, error) {
	reqURL := fmt.Sprintf("%s/signed-url?bucket=%s&clienteID=%s&upload=false", signerBaseURL, bucketName, clienteID)
	return http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
}