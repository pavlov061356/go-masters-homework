# Как использовать

## Скомпилируйте программу

```bash
go build -o loadtest
```

## Запустите нагрузочное тестирование

``` bash
./loadtest -std http://localhost:8080 -fast http://localhost:8081 -requests 10000 -concurrent 100
```

Параметры командной строки:

- `-std`: URL стандартного HTTP сервера (по умолчанию <http://localhost:8080>)

- `-fast`: URL fasthttp сервера (по умолчанию <http://localhost:8081>)

- `-requests`: Общее количество запросов (по умолчанию 1000)

- `-concurrent`: Количество одновременных запросов (по умолчанию 50)
