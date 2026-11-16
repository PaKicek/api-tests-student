package favorites

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// favoritesPageSize количество записей избранного на странице (максимальное значение из спецификации)
const favoritesPageSize = 50

// HttpGetFavorites выполняет GET-запрос для получения избранного пользователя с фильтрацией по названию объявления
func HttpGetFavorites(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("search", textForSearch).
		Query("page_size", strconv.Itoa(favoritesPageSize)).
		Expect(t).
		End().Response
}

// HttpPostFavorite выполняет POST-запрос для добавления объявления в избранное
func HttpPostFavorite(t *testing.T, token, advertisementId string) *http.Response {
	request := models.AddFavoriteRequest{AdvertisementId: advertisementId}

	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite выполняет DELETE-запрос для удаления объявления из избранного
func HttpDeleteFavorite(t *testing.T, token, advertisementId string) *http.Response {
	favoritePath := fmt.Sprintf(path.FavoritePath, advertisementId)

	return apiRunner.GetRunner().Auth(token).Create().Delete(favoritePath).
		Expect(t).
		End().Response
}
