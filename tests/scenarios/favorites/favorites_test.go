package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/constants/messages"
	"api-tests-template/internal/constants/token"
	"api-tests-template/internal/helpers/photos"
	"api-tests-template/internal/helpers/response"
	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	authModels "api-tests-template/internal/managers/auth/models"
	favoritesManager "api-tests-template/internal/managers/favorites"
	favoritesModels "api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	// advertisementDescription описание объявлений, которые создаются для подготовки тестов
	advertisementDescription = "Кресло офисное с регулировкой высоты"
	// advertisementPrice цена объявлений, которые создаются для подготовки тестов
	advertisementPrice = 50000
	// advertisementQuantity количество товара в объявлениях, которые создаются для подготовки тестов
	advertisementQuantity = 8
)

type TestSuite struct {
	suite.Suite
	login       authModels.LoginOkResponse
	sellerLogin authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация тестового пользователя (работает с избранным)")
	s.login = base.LoginTestUser(s.T())

	base.Precondition("Авторизация пользователя-продавца (владеет объявлениями, которые добавляются в избранное)")
	s.sellerLogin = base.LoginTestSeller(s.T())
}

// Positive: полный цикл работы с избранным - добавление объявления, проверка в списке, удаление из избранного и проверка, что объявление пропало из списка
func (s *TestSuite) TestFavoritesFullCycle() {
	t := s.T()

	base.Precondition("Продавец создаёт объявление, которое тестовый пользователь добавит в избранное")
	advertisement := s.createAdvertisement(t, s.sellerLogin, base.UniqueAdvertisementTitle("Кресло офисное"))

	t.Cleanup(func() {
		base.Step(fmt.Sprintf("Убираем объявление %s из избранного после теста", advertisement.Id))
		favoritesManager.RemoveFavoriteIfExists(t, s.login.Token, advertisement.Id)
	})

	var favorite favoritesModels.FavoriteItem

	s.Run("Добавляем объявление в избранное", func() {
		base.Step(fmt.Sprintf("POST /favorites: advertisement_id=%s", advertisement.Id))

		favorite = favoritesManager.AddFavorite(s.T(), s.login.Token, advertisement.Id)

		base.Check("Запись избранного создана и содержит объявление продавца")
		require.NotEmpty(s.T(), favorite.FavoriteId, "В ответе нет идентификатора записи избранного")
		require.NotEmpty(s.T(), favorite.AddedAt, "В ответе нет даты добавления в избранное")
		require.Equal(s.T(), advertisement.Id, favorite.Advertisement.Id, "В избранное добавлено другое объявление")
		require.Equal(s.T(), advertisement.Title, favorite.Advertisement.Title, "Название объявления в избранном не совпадает")
	})

	s.Run("Объявление появляется в списке избранного", func() {
		base.Step(fmt.Sprintf("GET /favorites?search=%s", advertisement.Title))

		favoritesList := favoritesManager.GetFavorites(s.T(), s.login.Token, advertisement.Title)
		favoriteItem := assertions.ContainsFavoriteAdvertisement(s.T(), favoritesList.Items, advertisement.Id)

		base.Check(fmt.Sprintf("Объявление найдено в избранном, всего записей: %d", favoritesList.Total))
		require.Equal(s.T(), favorite.FavoriteId, favoriteItem.FavoriteId, "Идентификатор записи избранного не совпадает")
	})

	s.Run("Удаляем объявление из избранного", func() {
		base.Step(fmt.Sprintf("DELETE /favorites/%s", advertisement.Id))

		favoritesManager.DeleteFavorite(s.T(), s.login.Token, advertisement.Id)

		base.Check("Сервер подтвердил удаление объявления из избранного")
	})

	s.Run("Объявление пропадает из списка избранного", func() {
		base.Step(fmt.Sprintf("GET /favorites?search=%s", advertisement.Title))

		favoritesList := favoritesManager.GetFavorites(s.T(), s.login.Token, advertisement.Title)

		base.Check("Удалённого объявления нет в списке избранного")
		assertions.NotContainsFavoriteAdvertisement(s.T(), favoritesList.Items, advertisement.Id)
	})
}

// Positive: если продавец удаляет само объявление, оно исчезает из избранного пользователя
func (s *TestSuite) TestFavoriteDeletedWithAdvertisement() {
	t := s.T()

	base.Precondition("Продавец создаёт объявление, тестовый пользователь добавляет его в избранное")
	advertisement := s.createAdvertisement(t, s.sellerLogin, base.UniqueAdvertisementTitle("Велосипед детский"))

	t.Cleanup(func() {
		base.Step(fmt.Sprintf("Убираем объявление %s из избранного после теста", advertisement.Id))
		favoritesManager.RemoveFavoriteIfExists(t, s.login.Token, advertisement.Id)
	})

	s.Run("Добавляем объявление в избранное", func() {
		base.Step(fmt.Sprintf("POST /favorites: advertisement_id=%s", advertisement.Id))

		favoritesManager.AddFavorite(s.T(), s.login.Token, advertisement.Id)

		base.Check("Объявление добавлено в избранное")
	})

	s.Run("Проверяем, что объявление есть в избранном", func() {
		base.Step(fmt.Sprintf("GET /favorites?search=%s", advertisement.Title))

		favoritesList := favoritesManager.GetFavorites(s.T(), s.login.Token, advertisement.Title)
		assertions.ContainsFavoriteAdvertisement(s.T(), favoritesList.Items, advertisement.Id)

		base.Check("Объявление есть в списке избранного")
	})

	s.Run("Продавец удаляет объявление", func() {
		base.Step(fmt.Sprintf("DELETE /advertisement?id=%s от имени продавца", advertisement.Id))

		advertisements.DeleteAdvertisement(s.T(), s.sellerLogin.Token, advertisement.Id)

		base.Check("Сервер подтвердил удаление объявления")
	})

	s.Run("Удалённое объявление пропадает из поиска", func() {
		base.Step(fmt.Sprintf("GET /advertisements?search=%s", advertisement.Title))

		found := advertisements.GetAdvertisementsList(s.T(), s.login.Token, advertisement.Title)

		base.Check("В поиске больше нет удалённого объявления")
		assertions.NotContainsAdvertisementByTitle(s.T(), found.Items, advertisement.Title)
	})

	s.Run("Объявление исчезает из избранного пользователя", func() {
		base.Step(fmt.Sprintf("GET /favorites?search=%s", advertisement.Title))

		favoritesList := favoritesManager.GetFavorites(s.T(), s.login.Token, advertisement.Title)

		base.Check("Удалённого объявления нет в списке избранного")
		assertions.NotContainsFavoriteAdvertisement(s.T(), favoritesList.Items, advertisement.Id)
	})
}

// Negative: бизнес-ограничения ручек избранного и запросы без авторизации
func (s *TestSuite) TestFavoritesNegative() {
	t := s.T()

	base.Precondition("Создаём объявление тестового пользователя - его нельзя добавить в собственное избранное")
	ownAdvertisement := s.createAdvertisement(t, s.login, base.UniqueAdvertisementTitle("Стол письменный чёрный"))

	base.Precondition("Создаём объявление продавца и добавляем его в избранное - на нём проверяем повторное добавление")
	sellerAdvertisement := s.createAdvertisement(t, s.sellerLogin, base.UniqueAdvertisementTitle("Тумба на ножках серая"))

	favoritesManager.AddFavorite(t, s.login.Token, sellerAdvertisement.Id)
	t.Cleanup(func() {
		base.Step(fmt.Sprintf("Убираем объявление %s из избранного после теста", sellerAdvertisement.Id))
		favoritesManager.RemoveFavoriteIfExists(t, s.login.Token, sellerAdvertisement.Id)
	})

	// идентификатор, которого гарантированно нет на сервере
	nonExistentAdvertisementId := utils.RandomUuid()

	cases := []struct {
		name      string
		request   func(t *testing.T) *http.Response
		wantError response.ExpectedError
	}{
		{
			name: "Добавление в избранное несуществующего объявления",
			request: func(t *testing.T) *http.Response {
				return favoritesManager.AddFavoriteRaw(t, s.login.Token, nonExistentAdvertisementId)
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusNotFound,
				Code:       messages.CodeNotFound,
				Message:    messages.AdvertisementNotFound,
			},
		},
		{
			name: "Добавление в избранное собственного объявления",
			request: func(t *testing.T) *http.Response {
				return favoritesManager.AddFavoriteRaw(t, s.login.Token, ownAdvertisement.Id)
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusForbidden,
				Code:       messages.CodeForbidden,
				Message:    messages.CannotAddOwnAdvertisement,
			},
		},
		{
			name: "Повторное добавление объявления в избранное",
			request: func(t *testing.T) *http.Response {
				return favoritesManager.AddFavoriteRaw(t, s.login.Token, sellerAdvertisement.Id)
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusConflict,
				Code:       messages.CodeAlreadyExists,
				Message:    messages.AdvertisementAlreadyInFavorites,
			},
		},
		{
			name: "Добавление в избранное без токена авторизации",
			request: func(t *testing.T) *http.Response {
				return favoritesManager.AddFavoriteRaw(t, "", sellerAdvertisement.Id)
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusUnauthorized,
				Code:       messages.CodeUnauthorized,
				Message:    messages.AuthorizationHeaderRequired,
			},
		},
		{
			name: "Добавление в избранное с некорректным токеном",
			request: func(t *testing.T) *http.Response {
				return favoritesManager.AddFavoriteRaw(t, token.InvalidToken, sellerAdvertisement.Id)
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusUnauthorized,
				Code:       messages.CodeUnauthorized,
				Message:    messages.InvalidSession,
			},
		},
		{
			name: "Получение избранного без токена авторизации",
			request: func(t *testing.T) *http.Response {
				return favoritesManager.GetFavoritesRaw(t, "", "")
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusUnauthorized,
				Code:       messages.CodeUnauthorized,
				Message:    messages.AuthorizationHeaderRequired,
			},
		},
		{
			name: "Удаление из избранного несуществующего объявления",
			request: func(t *testing.T) *http.Response {
				return favoritesManager.DeleteFavoriteRaw(t, s.login.Token, nonExistentAdvertisementId)
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusNotFound,
				Code:       messages.CodeNotFound,
				Message:    messages.FavoriteNotFound,
			},
		},
		{
			name: "Удаление из избранного без токена авторизации",
			request: func(t *testing.T) *http.Response {
				return favoritesManager.DeleteFavoriteRaw(t, "", sellerAdvertisement.Id)
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusUnauthorized,
				Code:       messages.CodeUnauthorized,
				Message:    messages.AuthorizationHeaderRequired,
			},
		},
	}

	base.Precondition(fmt.Sprintf("Проверяем %d негативных вариантов работы с избранным", len(cases)))

	for _, testCase := range cases {
		s.Run(testCase.name, func() {
			base.Step(fmt.Sprintf("Запрос к ручкам избранного. %s", testCase.name))

			raw := testCase.request(s.T())
			testCase.wantError.Assert(s.T(), raw)

			base.Check(fmt.Sprintf("API вернул %d: %s", testCase.wantError.StatusCode, testCase.wantError.Message))
		})
	}
}

// createAdvertisement создаёт объявление от имени переданного пользователя и удаляет его после теста, чтобы прогоны оставались идемпотентными
func (s *TestSuite) createAdvertisement(t *testing.T, login authModels.LoginOkResponse, title string) advModels.Advertisement {
	t.Helper()

	base.Step(fmt.Sprintf("POST /advertisement: создаём объявление %q от имени %s", title, login.User.Email))

	created := advertisements.CreateAdvertisement(t, login.Token, advModels.CreateAdvertisementRequest{
		Title:       title,
		Description: advertisementDescription,
		Price:       advertisementPrice,
		Quantity:    advertisementQuantity,
		PhotoPaths:  []string{photos.Path(photos.FirstPhoto)},
	})

	t.Cleanup(func() {
		base.Step(fmt.Sprintf("Удаляем объявление %s, созданное для теста", created.Id))
		advertisements.DeleteAdvertisementIfExists(t, login.Token, created.Id)
	})

	return created
}
