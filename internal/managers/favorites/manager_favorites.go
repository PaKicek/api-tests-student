package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/favorites/models"
)

// GetFavorites возвращает избранное пользователя с фильтрацией по названию объявления, ожидает 200
func GetFavorites(t *testing.T, token, textForSearch string) models.FavoritesList {
	t.Helper()

	body := response.Read(t, favorites.HttpGetFavorites(t, token, textForSearch))
	body.AssertStatus(t, http.StatusOK)

	return response.Decode[models.FavoritesList](t, body)
}

// GetFavoritesRaw отправляет запрос на получение избранного без проверки ответа (используется в негативных сценариях)
func GetFavoritesRaw(t *testing.T, token, textForSearch string) *http.Response {
	t.Helper()

	return favorites.HttpGetFavorites(t, token, textForSearch)
}

// AddFavorite добавляет объявление в избранное и возвращает созданную запись, ожидает 201
func AddFavorite(t *testing.T, token, advertisementId string) models.FavoriteItem {
	t.Helper()

	body := response.Read(t, favorites.HttpPostFavorite(t, token, advertisementId))
	body.AssertStatus(t, http.StatusCreated)

	return response.Decode[models.FavoriteItem](t, body)
}

// AddFavoriteRaw отправляет запрос на добавление объявления в избранное без проверки ответа (используется в негативных сценариях)
func AddFavoriteRaw(t *testing.T, token, advertisementId string) *http.Response {
	t.Helper()

	return favorites.HttpPostFavorite(t, token, advertisementId)
}

// DeleteFavorite удаляет объявление из избранного, ожидает 204
func DeleteFavorite(t *testing.T, token, advertisementId string) {
	t.Helper()

	response.Read(t, favorites.HttpDeleteFavorite(t, token, advertisementId)).
		AssertStatus(t, http.StatusNoContent)
}

// DeleteFavoriteRaw отправляет запрос на удаление объявления из избранного без проверки ответа (используется в негативных сценариях)
func DeleteFavoriteRaw(t *testing.T, token, advertisementId string) *http.Response {
	t.Helper()

	return favorites.HttpDeleteFavorite(t, token, advertisementId)
}

// RemoveFavoriteIfExists удаляет объявление из избранного. Используется при очистке данных после теста
func RemoveFavoriteIfExists(t *testing.T, token, advertisementId string) {
	t.Helper()

	response.Read(t, favorites.HttpDeleteFavorite(t, token, advertisementId)).
		AssertStatus(t, http.StatusNoContent, http.StatusNotFound)
}
