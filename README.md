## Подготовка к запуску

Создайте файл `.env` в папке `/config` с необходимыми данными. Пример содержания файла:
```env
GRPC_SERVER_HOST=0.0.0.0
GRPC_SERVER_PORT=50051
CDN_HOST=cdn.ru
```

## Запуск

### Запуск тестов и приложения
```bash
make full-run
```
### Запуск тестов
```bash
make test
```
### Запуск нагрузочного тестирования
```bash
make stress-test
```





