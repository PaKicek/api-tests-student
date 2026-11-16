package response

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Body прочитанный ответ API: статус-код и тело в виде строки
type Body struct {
	StatusCode int
	Content    string
}

// ErrorResponse тело ответа API с описанием ошибки
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ExpectedError ожидаемая ошибка API: статус-код, код и текст сообщения
type ExpectedError struct {
	StatusCode int
	Code       string
	Message    string
}

// Read читает тело ответа и закрывает его, чтобы вызывающий код не забывал про Close
func Read(t *testing.T, raw *http.Response) Body {
	t.Helper()

	defer func() {
		require.NoError(t, raw.Body.Close(), "Не удалось закрыть тело ответа")
	}()

	content, err := io.ReadAll(raw.Body)
	require.NoErrorf(t, err, "Не удалось прочитать тело ответа со статус-кодом %d", raw.StatusCode)

	return Body{StatusCode: raw.StatusCode, Content: string(content)}
}

// Decode разбирает тело ответа в переданную структуру
func Decode[T any](t *testing.T, body Body) T {
	t.Helper()

	var decoded T
	require.NoErrorf(t, json.Unmarshal([]byte(body.Content), &decoded),
		"Не удалось разобрать тело ответа в структуру. Тело ответа: %s", body.Content)

	return decoded
}

// AssertStatus проверяет, что ответ пришёл с одним из ожидаемых статус-кодов
func (b Body) AssertStatus(t *testing.T, expectedStatusCodes ...int) {
	t.Helper()

	for _, expectedStatusCode := range expectedStatusCodes {
		if b.StatusCode == expectedStatusCode {
			return
		}
	}

	require.Failf(t, "неожиданный статус-код ответа",
		"Ожидался статус-код %v, получен %d. Тело ответа: %s", expectedStatusCodes, b.StatusCode, b.Content)
}

// Assert проверяет, что ответ API содержит ожидаемую ошибку: статус-код, код и текст сообщения
func (e ExpectedError) Assert(t *testing.T, raw *http.Response) {
	t.Helper()

	body := Read(t, raw)
	errorResponse := Decode[ErrorResponse](t, body)

	body.AssertStatus(t, e.StatusCode)
	require.Equalf(t, e.Code, errorResponse.Error,
		"Ожидался код ошибки %q, получен %q. Тело ответа: %s", e.Code, errorResponse.Error, body.Content)
	require.Equalf(t, e.Message, errorResponse.Message,
		"Ожидалось сообщение об ошибке %q, получено %q", e.Message, errorResponse.Message)
}
