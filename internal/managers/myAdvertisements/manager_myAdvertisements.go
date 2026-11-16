package myAdvertisements

import (
	"testing"

	"api-tests-template/internal/client/http/myAdvertisements"
	"api-tests-template/internal/helpers/response"
)

// GetMyAdvertisements возвращает список объявлений авторизованного пользователя и проверяет ожидаемый статус-код
func GetMyAdvertisements(t *testing.T, token string, expectedStatusCode int, limit int) string {
	t.Helper()

	body := response.Read(t, myAdvertisements.HttpGetMyAdvertisements(t, token, limit))
	body.AssertStatus(t, expectedStatusCode)

	return body.Content
}
