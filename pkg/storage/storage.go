package storage

import (
	"api/internal/common"
	"api/internal/configuration"
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"mime"
	"net/http"
	"slices"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type storage struct {
	client     *minio.Client
	bucketName string
}

func New(cfg *configuration.S3) (*storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	return &storage{client: client, bucketName: cfg.BucketName}, nil
}

func (s *storage) PutImage(c context.Context, object []byte, resizeWidth int) (string, error) {
	mimeType := http.DetectContentType(object)

	supportedTypes := []string{"image/jpeg", "image/png", "image/webp"}
	if !slices.Contains(supportedTypes, mimeType) {
		return "", common.ErrUnsupportedFormat
	}

	if resizeWidth > 0 {
		resizedImage, err := thumbnail(object, resizeWidth)
		if err != nil {
			return "", err
		}
		object = resizedImage
	}

	return s.put(c, object)
}

func (s *storage) Put(c context.Context, object []byte) (string, error) {
	return s.put(c, object)
}

func (s *storage) put(c context.Context, object []byte) (string, error) {
	mimeType := http.DetectContentType(object)
	opts := minio.PutObjectOptions{
		ContentType: mimeType,
	}

	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", err
	}

	randomKey, err := generateRandomHash()
	if err != nil {
		return "", err
	}

	filename := randomKey + exts[0]

	print(s.bucketName)

	_, err = s.client.PutObject(c, s.bucketName, filename, bytes.NewReader(object), int64(len(object)), opts)
	return filename, err
}

func (s *storage) Remove(c context.Context, objectName string) error {
	return s.client.RemoveObject(c, s.bucketName, objectName, minio.RemoveObjectOptions{ForceDelete: true})
}

func generateRandomHash() (string, error) {
	randomData := make([]byte, 256)
	_, err := rand.Read(randomData)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(randomData)
	return hex.EncodeToString(hash[:]), nil
}
