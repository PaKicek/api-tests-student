package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// uniqueTitleSuffixLength длина случайного суффикса в названии объявления
const uniqueTitleSuffixLength = 6

// SetupSuite действия перед запуском сьюта с тестами
func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

// TearDownSuite действия после запуска сьюта с тестами
func TearDownSuite() {
	log.Println("Tear down suite")
}

// Precondition выводит текст в консоль с меткой [Precondition]
func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}

// Step выводит текст шага теста в консоль с меткой [Step]
func Step(text string) {
	utils.LogWithLabelAndTimeStamp("Step", text)
}

// Check выводит результат проверки в консоль с меткой [Check]
func Check(text string) {
	utils.LogWithLabelAndTimeStamp("Check", text)
}

// LoginTestUser авторизует основного тестового пользователя по кредам из переменных окружения
func LoginTestUser(t *testing.T) authModels.LoginOkResponse {
	t.Helper()

	return auth.Login(t, os.Getenv("TEST_USER_LOGIN"), os.Getenv("TEST_USER_PASSWORD"))
}

// LoginTestSeller авторизует пользователя-владельца объявлений, которые добавляются в избранное
func LoginTestSeller(t *testing.T) authModels.LoginOkResponse {
	t.Helper()

	return auth.Login(t, os.Getenv("TEST_SELLER_LOGIN"), os.Getenv("TEST_SELLER_PASSWORD"))
}

// UniqueAdvertisementTitle формирует название объявления, уникальное для каждого запуска теста, чтобы тесты не зависели от объявлений, созданных в предыдущих прогонах
func UniqueAdvertisementTitle(prefix string) string {
	return fmt.Sprintf("%s %s", prefix, utils.RandomString(uniqueTitleSuffixLength))
}
