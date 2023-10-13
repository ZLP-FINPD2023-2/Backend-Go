# Бэкенд ЛФП

## Эндпоинты

- `/api/v1`

  * `/auth`

    - `/login` - Вход
    - `/register` - Регистрация

  * `/schemes/swagger-ui` - Документация

## Правила взаимодействия с репозиторием

### Взаимодействие с ветками

В репозитори присутствуют 3 основные ветки:

- `master` - основная, стабильные релизы.
- `dev` - выходящая из `master`, тестовые релизы.
- `draft` - выходящая из `dev`, разработка.

Создавайте ветки из `draft` для разработки новых
функциональностей или исправления ошибок.

Каждый коммит должен быть сделан в рамках соответствующей
ветки, а не в `master`, `dev` или `draft`.
После завершения задачи сливайте ветку в основную ветку.

### Коммиты

#### Перед коммитом

- Прогоняйте стандартный форматер:

`
go fmt ./...
`

- Приводите в порядок пакеты:

`
go mod tidy
`

#### Оформление

Оформляйте коммиты в соответствии с
[Conventional Commits](https://www.conventionalcommits.org/ru/v1.0.0/),
который предоставляет структурированный формат для
описания изменений.

### Оформления кода

#### Именование

##### Переменные и функции

Используйте CamelCase для именования переменных и функций,
Например: CalculateIncome, totalBalance.

##### Пакеты

Именуйте пакеты короткими и понятными именами, состоящими
из одного слова или сокращений, если это общепринято
(например, db, api, utils).

##### Константы

Используйте верхний регистр для именования констант и
разделяйте слова символом подчеркивания.
Например: DEFAULT_PORT, MAX_RETRIES.

##### Именование ошибок

Именуйте переменные, представляющие ошибки,
с суффиксом Err. Например: err, dbErr.

#### Форматирование

##### Отступы

Используйте 1 символ табуляции для каждого уровня вложенности.

##### Максимальная длина строки

Старайтесь придерживаться максимальной длины строки
в 80-100 символов. Это улучшает читаемость кода.

##### Форматирование комментариев

Используйте комментарии в стиле
[GoDoc](https://go.dev/blog/godoc) для документирования
вашего кода.

Добавляйте комментарии к функциям, методам и структурам,
описывая их назначение и возвращаемые значения.

Используйте комментарии в конце строк для пояснения
сложных участков кода.

#### Импорты

Группируйте импорты:

- Стандартные библиотеки Go.
- Внешние зависимости.
- Локальные пакеты.

Пример:

```go
import (
    "fmt"
    "os"

    "github.com/external/package"
    "github.com/your/project/localpackage"
)
```

#### Обработка ошибок

Всегда проверяйте ошибки и обрабатывайте их адекватно.

Не игнорируйте ошибки с помощью _, если это не обосновано.

Возвращайте ошибки из функций и методов, если они могут
вызвать проблемы.
