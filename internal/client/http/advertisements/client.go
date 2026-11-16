package advertisements

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// HttpGetAdvertisements выполняет GET-запрос для получения списка объявлений по строке поиска
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementById выполняет GET-запрос для получения объявления по идентификатору
func HttpGetAdvertisementById(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", advertisementId).
		Expect(t).
		End().Response
}

// HttpPostAdvertisement выполняет POST-запрос для создания объявления вместе с фотографиями
func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		MultipartFormData("title", request.Title).
		MultipartFormData("description", request.Description).
		MultipartFormData("price", formatPrice(request.Price)).
		MultipartFormData("quantity", strconv.Itoa(request.Quantity)).
		MultipartFile("photos", request.PhotoPaths...).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement выполняет DELETE-запрос для удаления объявления по идентификатору
func HttpDeleteAdvertisement(t *testing.T, token, advertisementId string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", advertisementId).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos выполняет GET-запрос для получения фотографий объявления
func HttpGetAdvertisementPhotos(t *testing.T, token, advertisementId string) *http.Response {
	photosPath := fmt.Sprintf(path.AdvertisementPhotosPath, advertisementId)

	return apiRunner.GetRunner().Auth(token).Create().Get(photosPath).
		Expect(t).
		End().Response
}

// formatPrice приводит цену к виду, который ожидает сервер в поле формы price
func formatPrice(price float64) string {
	return strconv.FormatFloat(price, 'f', -1, 64)
}
