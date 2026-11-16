package advertisementCreate

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"unicode/utf8"

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
	storageManager "api-tests-template/internal/managers/storage"

	base "api-tests-template/tests"
)

const (
	// titleMaxLength максимальная длина названия объявления, которую принимает API
	titleMaxLength = 50
	// photosMaxCount максимальное количество фотографий, которое можно приложить к объявлению
	photosMaxCount = 3
	// priceMaxValue максимальная цена объявления
	priceMaxValue = 1000000
	// priceBelowMinValue отрицательная цена объявления
	priceBelowMinValue = -1000
	// quantityMinValue минимальное количество товара в объявлении
	quantityMinValue = 1
	// quantityMaxValue максимальное количество товара в объявлении
	quantityMaxValue = 100
	// quantityAboveMaxValue количество товара больше максимально допустимого
	quantityAboveMaxValue = quantityMaxValue + 1

	// advertisementDescriptionAllParameters описание объявления в позитивном сценарии
	advertisementDescriptionAllParameters = "Горный велосипед, рама M, пробег небольшой"
	// advertisementDescriptionBoundary описание объявления в сценарии с граничными значениями
	advertisementDescriptionBoundary = "Горный велосипед, рама XL, состояние новое"
	// advertisementDescriptionForNegative описание корректного запроса в негативных сценариях
	advertisementDescriptionForNegative = "Настольная лампа"
)

type TestSuite struct {
	suite.Suite
	login authModels.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация тестового пользователя с кредами из переменных окружения")
	s.login = base.LoginTestUser(s.T())
}

// Positive: объявление со всеми параметрами создаётся, доступно по идентификатору, находится в поиске, а загруженные фотографии доступны по ручке фотографий и в хранилище
func (s *TestSuite) TestCreateAdvertisementPositive() {
	t := s.T()

	title := base.UniqueAdvertisementTitle("Горный велосипед")
	request := advModels.CreateAdvertisementRequest{
		Title:       title,
		Description: advertisementDescriptionAllParameters,
		Price:       45000,
		Quantity:    2,
		PhotoPaths:  []string{photos.Path(photos.FirstPhoto), photos.Path(photos.SecondPhoto)},
	}

	base.Precondition(fmt.Sprintf("Подготовлены %d фотографии из каталога testdata и уникальное название %q",
		len(request.PhotoPaths), title))

	var created advModels.Advertisement

	s.Run("Создаём объявление со всеми параметрами", func() {
		base.Step(fmt.Sprintf("POST /advertisement: title=%q, price=%v, quantity=%d, фотографий=%d",
			request.Title, request.Price, request.Quantity, len(request.PhotoPaths)))

		created = advertisements.CreateAdvertisement(s.T(), s.login.Token, request)
		t.Cleanup(func() { s.deleteAdvertisement(t, created.Id) })

		base.Check("Объявление создано: в ответе есть идентификатор, владелец и фотографии")
		require.NotEmpty(s.T(), created.Id, "В ответе на создание объявления отсутствует идентификатор")
		require.Equal(s.T(), s.login.User.Id, created.UserId, "Объявление создано от имени другого пользователя")
		require.Len(s.T(), created.Photos, len(request.PhotoPaths),
			"Количество фотографий в ответе не совпадает с количеством загруженных")
	})

	s.Run("Объявление доступно по ручке GET /advertisement", func() {
		base.Step(fmt.Sprintf("GET /advertisement?id=%s", created.Id))

		fetched := advertisements.GetAdvertisementById(s.T(), s.login.Token, created.Id)

		base.Check("Параметры объявления из ручки совпадают с переданными при создании")
		require.Equal(s.T(), created.Id, fetched.Id, "Идентификатор объявления не совпадает")
		assertions.AdvertisementEquals(s.T(), advModels.Advertisement{
			Title:       request.Title,
			Description: request.Description,
			Price:       request.Price,
			Quantity:    request.Quantity,
			UserId:      s.login.User.Id,
		}, fetched)
		require.Len(s.T(), fetched.Photos, len(request.PhotoPaths),
			"Количество фотографий в объявлении не совпадает с количеством загруженных")
	})

	s.Run("Объявление находится в поиске по полному названию", func() {
		base.Step(fmt.Sprintf("GET /advertisements?search=%s", url.QueryEscape(title)))

		found := advertisements.GetAdvertisementsList(s.T(), s.login.Token, title)
		advertisement := assertions.FindAdvertisementByTitle(s.T(), found.Items, title)

		base.Check(fmt.Sprintf("Среди %d найденных объявлений есть созданное объявление %s", len(found.Items), created.Id))
		require.Equal(s.T(), created.Id, advertisement.Id, "В поиске найдено объявление с другим идентификатором")
	})

	s.Run("Фотографии доступны по ручке GET /advertisements/{id}/photos", func() {
		base.Step(fmt.Sprintf("GET /advertisements/%s/photos", created.Id))

		advertisementPhotos := advertisements.GetAdvertisementPhotos(s.T(), s.login.Token, created.Id)

		require.Len(s.T(), advertisementPhotos, len(request.PhotoPaths),
			"Количество фотографий не совпадает с количеством загруженных")
		require.ElementsMatch(s.T(), assertions.PhotoIds(created.Photos), assertions.PhotoIds(advertisementPhotos),
			"Список фотографий не совпадает с фотографиями из ответа на создание объявления")

		for _, photo := range advertisementPhotos {
			base.Step(fmt.Sprintf("Проверяем доступность фотографии по ссылке %s", photo.Url))
			require.NotEmpty(s.T(), photo.Url, "В ответе ручки фотографий нет ссылки на файл")
			storageManager.IsPhotoAvailable(s.T(), photo.Url)
		}

		base.Check("Все фотографии загружены в хранилище и отдаются как изображения")
	})
}

// Positive: объявление создаётся с граничными допустимыми значениями параметров
func (s *TestSuite) TestCreateAdvertisementWithBoundaryValues() {
	t := s.T()

	request := advModels.CreateAdvertisementRequest{
		Title:       maxLengthTitle(),
		Description: advertisementDescriptionBoundary,
		Price:       priceMaxValue,
		Quantity:    quantityMaxValue,
		PhotoPaths: []string{
			photos.Path(photos.FirstPhoto),
			photos.Path(photos.SecondPhoto),
			photos.Path(photos.ThirdPhoto),
		},
	}

	base.Precondition(fmt.Sprintf("Готовим объявление с названием из %d символов, ценой %d, количеством %d и %d фотографиями",
		titleMaxLength, priceMaxValue, quantityMaxValue, photosMaxCount))

	var created advModels.Advertisement

	s.Run("Создаём объявление с граничными значениями параметров", func() {
		base.Step(fmt.Sprintf("POST /advertisement: title=%q, price=%d, quantity=%d, фотографий=%d",
			request.Title, priceMaxValue, request.Quantity, len(request.PhotoPaths)))

		created = advertisements.CreateAdvertisement(s.T(), s.login.Token, request)
		t.Cleanup(func() { s.deleteAdvertisement(t, created.Id) })

		base.Check("Объявление создано с максимальными допустимыми значениями параметров")
		require.NotEmpty(s.T(), created.Id, "В ответе на создание объявления отсутствует идентификатор")
		require.Equal(s.T(), float64(priceMaxValue), created.Price, "Цена объявления не совпадает с переданной")
		require.Equal(s.T(), quantityMaxValue, created.Quantity, "Количество объявления не совпадает с переданным")
		require.Len(s.T(), created.Photos, photosMaxCount, "Количество фотографий не совпадает с загруженным")
	})

	s.Run("Объявление доступно по ручке GET /advertisement", func() {
		base.Step(fmt.Sprintf("GET /advertisement?id=%s", created.Id))

		fetched := advertisements.GetAdvertisementById(s.T(), s.login.Token, created.Id)

		base.Check("Название, цена и количество объявления сохранились без изменений")
		require.Equal(s.T(), request.Title, fetched.Title, "Название объявления не совпадает с переданным")
		require.Equal(s.T(), float64(priceMaxValue), fetched.Price, "Цена объявления не совпадает с переданной")
		require.Equal(s.T(), quantityMaxValue, fetched.Quantity, "Количество объявления не совпадает с переданным")
	})
}

// Negative: объявление с некорректными параметрами не создаётся, API возвращает ошибку валидации
func (s *TestSuite) TestCreateAdvertisementWithInvalidParameters() {
	tooLongTitle := strings.Repeat("я", titleMaxLength+1)

	cases := []struct {
		name      string
		modify    func(request *advModels.CreateAdvertisementRequest)
		wantError response.ExpectedError
	}{
		{
			name: "Пустой title",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Title = ""
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    messages.TitleRequired,
			},
		},
		{
			name: fmt.Sprintf("Title длиннее %d символов", titleMaxLength),
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Title = tooLongTitle
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    fmt.Sprintf(messages.TitleTooLongTemplate, titleMaxLength, utf8.RuneCountInString(tooLongTitle)),
			},
		},
		{
			name: "Пустой description",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Description = ""
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    messages.DescriptionRequired,
			},
		},
		{
			name: "Отрицательная цена",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Price = priceBelowMinValue
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    messages.PriceOutOfRange,
			},
		},
		{
			name: fmt.Sprintf("Цена больше %d", priceMaxValue),
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Price = priceMaxValue + 1
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    messages.PriceOutOfRange,
			},
		},
		{
			name: "Количество меньше минимального",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Quantity = quantityMinValue - 1
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    messages.QuantityOutOfRange,
			},
		},
		{
			name: fmt.Sprintf("Количество больше %d", quantityMaxValue),
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.Quantity = quantityAboveMaxValue
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    messages.QuantityOutOfRange,
			},
		},
		{
			name: "Без фотографий",
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.PhotoPaths = nil
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    messages.AtLeastOnePhoto,
			},
		},
		{
			name: fmt.Sprintf("Больше %d фотографий", photosMaxCount),
			modify: func(request *advModels.CreateAdvertisementRequest) {
				request.PhotoPaths = []string{
					photos.Path(photos.FirstPhoto),
					photos.Path(photos.SecondPhoto),
					photos.Path(photos.ThirdPhoto),
					photos.Path(photos.ThirdPhoto),
				}
			},
			wantError: response.ExpectedError{
				StatusCode: http.StatusBadRequest,
				Code:       messages.CodeInvalidRequest,
				Message:    messages.TooManyPhotos,
			},
		},
	}

	base.Precondition(fmt.Sprintf("Проверяем %d вариантов некорректных параметров объявления", len(cases)))

	for _, testCase := range cases {
		s.Run(testCase.name, func() {
			request := s.validCreateRequest()
			testCase.modify(&request)

			base.Step(fmt.Sprintf("POST /advertisement с некорректными параметрами. %s", testCase.name))

			raw := advertisements.CreateAdvertisementRaw(s.T(), s.login.Token, request)
			testCase.wantError.Assert(s.T(), raw)

			base.Check(fmt.Sprintf("Объявление не создано, API вернул %d: %s",
				testCase.wantError.StatusCode, testCase.wantError.Message))
		})
	}
}

// Negative: создание объявления без авторизации и с некорректным токеном возвращает 401
func (s *TestSuite) TestCreateAdvertisementWithoutAuth() {
	cases := []struct {
		name      string
		authToken string
		wantError response.ExpectedError
	}{
		{
			name:      "Без токена авторизации",
			authToken: "",
			wantError: response.ExpectedError{
				StatusCode: http.StatusUnauthorized,
				Code:       messages.CodeUnauthorized,
				Message:    messages.AuthorizationHeaderRequired,
			},
		},
		{
			name:      "С некорректным токеном",
			authToken: token.InvalidToken,
			wantError: response.ExpectedError{
				StatusCode: http.StatusUnauthorized,
				Code:       messages.CodeUnauthorized,
				Message:    messages.InvalidSession,
			},
		},
	}

	base.Precondition("Корректный запрос на создание объявления передаём без авторизации")

	for _, testCase := range cases {
		s.Run(testCase.name, func() {
			request := s.validCreateRequest()

			base.Step(fmt.Sprintf("POST /advertisement без авторизации. %s", testCase.name))

			raw := advertisements.CreateAdvertisementRaw(s.T(), testCase.authToken, request)
			testCase.wantError.Assert(s.T(), raw)

			base.Check(fmt.Sprintf("Объявление не создано, API вернул %d: %s",
				testCase.wantError.StatusCode, testCase.wantError.Message))
		})
	}
}

// validCreateRequest формирует корректный запрос на создание объявления, от которого отталкиваются негативные проверки
func (s *TestSuite) validCreateRequest() advModels.CreateAdvertisementRequest {
	return advModels.CreateAdvertisementRequest{
		Title:       base.UniqueAdvertisementTitle("Настольная лампа"),
		Description: advertisementDescriptionForNegative,
		Price:       2500,
		Quantity:    quantityMinValue,
		PhotoPaths:  []string{photos.Path(photos.FirstPhoto)},
	}
}

// deleteAdvertisement удаляет объявление, созданное в рамках теста, чтобы прогоны были идемпотентными
func (s *TestSuite) deleteAdvertisement(t *testing.T, advertisementId string) {
	t.Helper()

	if advertisementId == "" {
		return
	}

	base.Step(fmt.Sprintf("Удаляем объявление %s, созданное в рамках теста", advertisementId))
	advertisements.DeleteAdvertisementIfExists(t, s.login.Token, advertisementId)
}

// maxLengthTitle формирует название максимальной допустимой длины (уникальный префикс добивается символами до лимита API)
func maxLengthTitle() string {
	prefix := base.UniqueAdvertisementTitle("Горный велосипед ")

	return prefix + strings.Repeat("X", titleMaxLength-utf8.RuneCountInString(prefix))
}
