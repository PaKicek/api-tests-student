package favorites

import (
	"fmt"
	"net/http"
	"testing"
	"time"

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
	// favoritesLimit максимальное количество записей избранного на пользователя
	favoritesLimit = 50
	// rateLimitDelay пауза перед запросами, которые создают, добавляют и удаляют данные (сервер ограничивает частоту таких запросов и выдаёт ошибку 503)
	rateLimitDelay = time.Second
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

// waitForRateLimit выдерживает паузу перед запросом, чтобы не превысить ограничение частоты запросов на сервере
func waitForRateLimit() {
	time.Sleep(rateLimitDelay)
}

// Positive: полный цикл работы с избранным - добавление объявления, проверка в списке, удаление из избранного и проверка, что объявление пропало из списка
func (s *TestSuite) TestFavoritesFullCycle() {
	t := s.T()

	base.Precondition("Продавец создаёт объявление, которое тестовый пользователь добавит в избранное")
	advertisement := s.createAdvertisement(t, s.sellerLogin, base.UniqueAdvertisementTitle("Кресло офисное"))

	t.Cleanup(func() {
		base.Step(fmt.Sprintf("Убираем объявление %s из избранного после теста", advertisement.Id))
		s.removeFavoriteIfExists(t, s.login.Token, advertisement.Id)
	})

	var favorite favoritesModels.FavoriteItem

	s.Run("Добавляем объявление в избранное", func() {
		base.Step(fmt.Sprintf("POST /favorites: advertisement_id=%s", advertisement.Id))

		favorite = s.addFavorite(s.T(), s.login.Token, advertisement.Id)

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

		s.deleteFavorite(s.T(), s.login.Token, advertisement.Id)

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
		s.removeFavoriteIfExists(t, s.login.Token, advertisement.Id)
	})

	s.Run("Добавляем объявление в избранное", func() {
		base.Step(fmt.Sprintf("POST /favorites: advertisement_id=%s", advertisement.Id))

		s.addFavorite(s.T(), s.login.Token, advertisement.Id)

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

	s.addFavorite(t, s.login.Token, sellerAdvertisement.Id)
	t.Cleanup(func() {
		base.Step(fmt.Sprintf("Убираем объявление %s из избранного после теста", sellerAdvertisement.Id))
		s.removeFavoriteIfExists(t, s.login.Token, sellerAdvertisement.Id)
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

// Negative: пользователь не может добавить в избранное больше 50 объявлений - на 51 объявлении сервер возвращает ошибку лимита.
// Тест долгий, так как для достижения лимита нужно создать 50 объявлений и 50 записей избранного
func (s *TestSuite) TestFavoritesLimitReached() {
	t := s.T()

	if testing.Short() {
		s.T().Skip("в режиме -short тест пропускается: он создаёт 50 объявлений и 50 записей избранного")
	}

	base.Precondition(fmt.Sprintf("Между запросами на создание объявления, добавление в избранное и удаление "+
		"выдерживаем паузу %s, чтобы не превысить ограничение частоты запросов на сервере", rateLimitDelay))

	base.Precondition("Очищаем избранное тестового пользователя, чтобы тест начинался с пустого списка")
	s.clearFavorites(t)

	base.Precondition(fmt.Sprintf("Создаём %d объявления продавца: %d добавляем в избранное, на последнем проверяем лимит",
		favoritesLimit+1, favoritesLimit))

	favoriteAdvertisements := make([]advModels.Advertisement, 0, favoritesLimit+1)
	for i := 1; i <= favoritesLimit+1; i++ {
		title := base.UniqueAdvertisementTitle(fmt.Sprintf("Кресло офисное %d", i))
		favoriteAdvertisements = append(favoriteAdvertisements, s.createAdvertisement(t, s.sellerLogin, title))
	}

	base.Step(fmt.Sprintf("POST /favorites для %d объявлений продавца", favoritesLimit))
	for _, advertisement := range favoriteAdvertisements[:favoritesLimit] {
		s.addFavorite(t, s.login.Token, advertisement.Id)
	}

	base.Check(fmt.Sprintf("В избранном ровно %d записей - лимит достигнут", favoritesLimit))
	favoritesList := favoritesManager.GetFavorites(t, s.login.Token, "")
	require.Equal(t, favoritesLimit, favoritesList.Total, "Количество записей в избранном не совпадает с лимитом")

	advertisementAboveLimit := favoriteAdvertisements[favoritesLimit]

	s.Run("Добавление объявления сверх лимита возвращает ошибку", func() {
		base.Step(fmt.Sprintf("POST /favorites: advertisement_id=%s сверх лимита", advertisementAboveLimit.Id))

		wantError := response.ExpectedError{
			StatusCode: http.StatusBadRequest,
			Code:       messages.CodeLimitReached,
			Message:    messages.FavoritesLimitReached,
		}

		raw := favoritesManager.AddFavoriteRaw(s.T(), s.login.Token, advertisementAboveLimit.Id)
		wantError.Assert(s.T(), raw)

		base.Check(fmt.Sprintf("API вернул %d: %s", wantError.StatusCode, wantError.Message))
	})

	s.Run("Объявление сверх лимита не попало в избранное", func() {
		base.Step("GET /favorites")

		favoritesList := favoritesManager.GetFavorites(s.T(), s.login.Token, "")

		base.Check(fmt.Sprintf("В избранном по-прежнему %d записей", favoritesLimit))
		require.Equal(s.T(), favoritesLimit, favoritesList.Total, "Количество записей в избранном изменилось после отказа")
		assertions.NotContainsFavoriteAdvertisement(s.T(), favoritesList.Items, advertisementAboveLimit.Id)
	})
}

// createAdvertisement создаёт объявление от имени переданного пользователя и удаляет его после теста, чтобы прогоны оставались идемпотентными
func (s *TestSuite) createAdvertisement(t *testing.T, login authModels.LoginOkResponse, title string) advModels.Advertisement {
	t.Helper()

	base.Step(fmt.Sprintf("POST /advertisement: создаём объявление %q от имени %s", title, login.User.Email))
	waitForRateLimit()

	created := advertisements.CreateAdvertisement(t, login.Token, advModels.CreateAdvertisementRequest{
		Title:       title,
		Description: advertisementDescription,
		Price:       advertisementPrice,
		Quantity:    advertisementQuantity,
		PhotoPaths:  []string{photos.Path(photos.FirstPhoto)},
	})

	t.Cleanup(func() {
		base.Step(fmt.Sprintf("Удаляем объявление %s, созданное для теста", created.Id))
		waitForRateLimit()
		advertisements.DeleteAdvertisementIfExists(t, login.Token, created.Id)
	})

	return created
}

// addFavorite добавляет объявление в избранное с паузой перед запросом
func (s *TestSuite) addFavorite(t *testing.T, token, advertisementId string) favoritesModels.FavoriteItem {
	t.Helper()

	waitForRateLimit()

	return favoritesManager.AddFavorite(t, token, advertisementId)
}

// deleteFavorite удаляет объявление из избранного с паузой перед запросом
func (s *TestSuite) deleteFavorite(t *testing.T, token, advertisementId string) {
	t.Helper()

	waitForRateLimit()
	favoritesManager.DeleteFavorite(t, token, advertisementId)
}

// removeFavoriteIfExists удаляет объявление из избранного с паузой перед запросом, используется при очистке данных после теста
func (s *TestSuite) removeFavoriteIfExists(t *testing.T, token, advertisementId string) {
	t.Helper()

	waitForRateLimit()
	favoritesManager.RemoveFavoriteIfExists(t, token, advertisementId)
}

// clearFavorites удаляет все записи избранного у тестового пользователя, чтобы прогоны оставались идемпотентными (т.е. тест на лимит начинался с пустого списка)
func (s *TestSuite) clearFavorites(t *testing.T) {
	t.Helper()

	favoritesList := favoritesManager.GetFavorites(t, s.login.Token, "")
	for _, item := range favoritesList.Items {
		favoritesManager.DeleteFavorite(t, s.login.Token, item.Advertisement.Id)
	}

	base.Check(fmt.Sprintf("Избранное тестового пользователя пустое, удалено записей: %d", len(favoritesList.Items)))
}
