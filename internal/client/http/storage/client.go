package storage

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// requestTimeout ограничение времени на скачивание файла из хранилища фотографий
const requestTimeout = 15 * time.Second

// HttpGetFileByUrl скачивает файл по публичной ссылке, которую API вернул для фотографии
func HttpGetFileByUrl(t *testing.T, fileUrl string) *http.Response {
	t.Helper()

	client := &http.Client{Timeout: requestTimeout}

	raw, err := client.Get(fileUrl)
	require.NoErrorf(t, err, "Не удалось скачать файл по ссылке %s", fileUrl)

	return raw
}
