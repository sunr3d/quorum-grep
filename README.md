# quorum-grep

Распределенная CLI-утилита для поиска текста в файлах с поддержкой кворума и параллельной обработки.

## Архитектура

- **3 сервера** в Docker контейнерах (порты 50051, 50052, 50053)
- **gRPC** для сетевого взаимодействия между клиентом и серверами
- **Кворум N/2+1** для обеспечения отказоустойчивости
- **Параллельная обработка** данных с использованием goroutines
- **Clean Architecture** с разделением на слои
- **Дедупликация результатов** по номерам строк

## Функциональность

- ✅ Все основные флаги `grep`: `-n`, `-A`, `-B`, `-C`, `-c`, `-i`, `-v`, `-F`
- ✅ Поиск в файлах и stdin
- ✅ Контекстные флаги (`-A`, `-B`, `-C`) с перекрывающимися чанками
- ✅ Отказоустойчивость через кворум
- ✅ Параллельная обработка данных
- ✅ Graceful shutdown серверов

## Установка и запуск

### 1. Запуск серверов

```bash
# Запуск всех серверов в Docker контейнерах
make up
```

### 2. Сборка клиента

```bash
# Сборка исполняемого файла
make build-client

# Результат: исполняемый файл `mygrep`
```

### 3. Использование

```bash
# Базовый поиск
./mygrep pattern file.txt

# С номерами строк
./mygrep -n pattern file.txt

# С контекстом
./mygrep -A 2 -B 1 pattern file.txt

# Подсчет совпадений
./mygrep -c pattern file.txt

# Игнорирование регистра
./mygrep -i PATTERN file.txt

# Инвертированный поиск
./mygrep -v pattern file.txt

# Фиксированная строка
./mygrep -F "exact.pattern" file.txt
```

## Примеры использования

### Базовые команды

```bash
# Поиск в файле
./mygrep "error" /var/log/app.log

# Поиск с контекстом (2 строки после, 1 строка до)
./mygrep -A 2 -B 1 "ERROR" /var/log/app.log

# Подсчет количества совпадений
./mygrep -c "warning" /var/log/app.log

# Поиск в stdin
echo "test line\npattern found\nanother line" | ./mygrep pattern
```

## Конфигурация

Настройки клиента находятся в `config.yaml`:

```yaml
CLIENT:
  SERVER_LIST:
    - "localhost:50051"
    - "localhost:50052"
    - "localhost:50053"
  TIMEOUT: 30s
  CHUNK_SIZE: 1024
```

## Структура проекта

```
├── cmd/
│   ├── client/          # CLI клиент
│   └── server/          # gRPC сервер
├── internal/
│   ├── client/          # Клиентская логика
│   ├── server/          # Серверная логика
│   ├── services/        # Бизнес-логика
│   ├── handlers/        # gRPC обработчики
│   ├── config/          # Конфигурация
│   └── entrypoint/      # Точки входа
├── models/              # Доменные модели
├── proto/               # gRPC протоколы (Proto stub)
└── api/                 # API определения (protobuf)
```

### Команды разработки

```bash
# Запуск тестов
make test

# Линтер
make lint

# Форматирование кода
make fmt

# Генерация protobuf
make protogen

# Сборка клиента
make build-client

# Сравнительное тестирование
make test-comparison

# Логи серверов
make logs
```

### Тестирование

```bash
# Все тесты
go test -v ./...

# Тесты с покрытием
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Сравнительное тестирование с оригинальным grep
make test-comparison
```

### Мониторинг

```bash
# Статус контейнеров
docker compose ps

# Логи всех серверов
docker compose logs -f

# Перезапуск
make restart
```

## Производительность

- **Параллельная обработка**: Каждый чанк обрабатывается в отдельной горутине
- **Кворум**: Система работает при отказе до N/2 серверов
- **Оптимизация памяти**: Использование `[]byte` вместо `string` для минимизации аллокаций
- **Контекстные флаги**: Перекрывающиеся чанки для корректной обработки `-A`, `-B`, `-C`
