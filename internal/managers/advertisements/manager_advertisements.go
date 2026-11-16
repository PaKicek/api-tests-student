package advertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements/models"
)

// GetAdvertisementsList возвращает список объявлений, найденных по строке поиска, ожидает 200
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	t.Helper()

	body := response.Read(t, advertisements.HttpGetAdvertisements(t, token, text))
	body.AssertStatus(t, http.StatusOK)

	return response.Decode[models.AdvertisementsList](t, body)
}

// GetAdvertisementById возвращает объявление по идентификатору, ожидает 200
func GetAdvertisementById(t *testing.T, token, advertisementId string) models.Advertisement {
	t.Helper()

	body := response.Read(t, advertisements.HttpGetAdvertisementById(t, token, advertisementId))
	body.AssertStatus(t, http.StatusOK)

	return response.Decode[models.Advertisement](t, body)
}

// CreateAdvertisement создаёт объявление с фотографиями и возвращает его параметры, ожидает 201
func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) models.Advertisement {
	t.Helper()

	body := response.Read(t, advertisements.HttpPostAdvertisement(t, token, request))
	body.AssertStatus(t, http.StatusCreated)

	return response.Decode[models.Advertisement](t, body)
}

// CreateAdvertisementRaw отправляет запрос на создание объявления без проверки ответа (используется в негативных сценариях, где статус-код и текст ошибки проверяет сам тест)
func CreateAdvertisementRaw(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	t.Helper()

	return advertisements.HttpPostAdvertisement(t, token, request)
}

// GetAdvertisementPhotos возвращает фотографии объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, advertisementId string) []models.Photo {
	t.Helper()

	body := response.Read(t, advertisements.HttpGetAdvertisementPhotos(t, token, advertisementId))
	body.AssertStatus(t, http.StatusOK)

	return response.Decode[[]models.Photo](t, body)
}

// DeleteAdvertisement удаляет объявление пользователя, ожидает 204
func DeleteAdvertisement(t *testing.T, token, advertisementId string) {
	t.Helper()

	response.Read(t, advertisements.HttpDeleteAdvertisement(t, token, advertisementId)).
		AssertStatus(t, http.StatusNoContent)
}

// DeleteAdvertisementIfExists удаляет объявление и не считает ошибкой случаи, когда объявления уже нет: на удаление чужого или несуществующего объявления сервер отвечает 403.
// Используется при очистке данных после теста, чтобы повторные прогоны оставались идемпотентными
func DeleteAdvertisementIfExists(t *testing.T, token, advertisementId string) {
	t.Helper()

	response.Read(t, advertisements.HttpDeleteAdvertisement(t, token, advertisementId)).
		AssertStatus(t, http.StatusNoContent, http.StatusForbidden, http.StatusNotFound)
}
