package storage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	storageClient "api-tests-template/internal/client/http/storage"
	"api-tests-template/internal/helpers/response"
)

// contentTypePrefixImage префикс content-type, который должен вернуться для файла изображения
const contentTypePrefixImage = "image/"

// IsPhotoAvailable проверяет, что фотография действительно загружена в хранилище (по публичной ссылке возвращается 200 и файл изображения)
func IsPhotoAvailable(t *testing.T, photoUrl string) {
	t.Helper()

	raw := storageClient.HttpGetFileByUrl(t, photoUrl)
	contentType := raw.Header.Get("Content-Type")

	body := response.Read(t, raw)
	body.AssertStatus(t, http.StatusOK)

	require.Containsf(t, contentType, contentTypePrefixImage,
		"По ссылке %s вернулся файл с content-type %q вместо изображения", photoUrl, contentType)
	require.NotEmptyf(t, body.Content, "По ссылке %s вернулся пустой файл", photoUrl)
}
