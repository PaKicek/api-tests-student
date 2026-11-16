# API-тесты для проекта https://testboard.avito.com/

## Полезные ссылки

1. [Swagger-спецификация](https://testboard.avito.com/swagger/index.html)
2. [Техническое задание](https://docs.google.com/document/d/1IZ9SxEcePoQT6lyEmukycp1qNok4HZtok73UQnPwKSA/edit?usp=sharing)
3. [Библиотека для написания API-тестов](https://github.com/steinfletcher/apitest)
4. [Мощная библиотека анализа JSON для Go](https://github.com/tidwall/gjson)

## Запуск тестов

1. Установите любое IDE для разработки: [Jetbrains GoLand](https://www.jetbrains.com/ru-ru/go/) (Платное или по студенческой лицензии), [VS Code](https://code.visualstudio.com/) (бесплатное).
2. Установите [Git](https://git-scm.com/install/) на компьютер
3. [Установите](https://go.dev/doc/install) Go версии 1.25.1
4. [Настройте](https://docs.github.com/ru/authentication/connecting-to-github-with-ssh) доступ к своему Github-аккаунту
   по SSH-ключу.
5. Склонируйте себе репозиторий: `git clone <путь к вашему репозиторию>`
6. Подтяните необходимые зависимости через команду в корне проекта:

```bash
   go mod tidy
```

7. В файле `.env.override` (создайте его в корне, если нет) выставьте правильные значения для переменных окружения:

```
TEST_USER_LOGIN=your_user_login
TEST_USER_PASSWORD=your_user_password
TEST_SELLER_LOGIN=your_seller_login
TEST_SELLER_PASSWORD=your_seller_password
```

8. Для запуска всех тестов (с отключенным кешированием) достаточно выполнить команду в корне проекта:

```bash
  go test -v -count=1 ./...
```

9. Отдельный пакет с тестами (с отключенным кешированием) можно запустить командой:

```bash
  go test -v -count=1 ./tests/scenarios/myAdvertisement/
```

## Текущие тестовые сценарии

| Пакет                                    | Что проверяем                                                                                          |
|------------------------------------------|--------------------------------------------------------------------------------------------------------|
| `tests/scenarios/advertisementSearch`    | Поиск объявлений по названию (по подстроке, и регистру)                                                       |
| `tests/scenarios/myAdvertisement`        | Список собственных объявлений и авторизация приватной ручки                                             |
| `tests/scenarios/advertisementCreate`    | Создание объявления со всеми параметрами, доступность по `GET /advertisement`, поиск и загрузка фотографий |
| `tests/scenarios/favorites`              | Полный цикл работы с избранным, удаление объявления владельцем и негативные случаи для ручек избранного       |

Сценарии создания объявлений и избранного создают собственные тестовые данные и удаляют их после
теста, поэтому их можно запускать несколько раз подряд - результат не зависит от предыдущих прогонов.

Фотографии для загрузки объявлений лежат в каталоге `testdata/photos`.

## Логирование

Для включения логирования request/response в консоль: выставьте переменную окружения `DEBUG=true` в
`.env.override`

## Запуск линтера

1. Установить приложение [golangci-lint](https://golangci-lint.run/docs/welcome/install/local/)
2. Запуск линтера (из корня репозитория):

```bash
  golangci-lint run
```