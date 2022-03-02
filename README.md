# Тестовое задание на позицию стажера-бекендера

## Задание

Описание задания вынесено в [TASK.md](docs/TASK.md)

## Решение

### Ответы на возникшие вопросы

> В нашей компании есть много различных микросервисов. Многие из них так или иначе хотят взаимодействовать
> с балансом пользователя. На архитектурном комитете приняли решение централизовать работу с балансом пользователя
> в отдельный сервис.

Исходим из того, что у пользователя один баланс в системе. Тогда можно создать абстракцию счёта для пользователя.
Счёт пользователя и пользователь имеют связь один к одному, так что можно работать с идентификатором пользователя
как с индентификатором и его счёта

### API
Спецификация API составлена при помощи OpenAPI: [openapi.yaml](docs/api/openapi.yaml) 

### Схема БД

![img of schema](docs/images/db_schema.png)

SQL скрипты представлены в директории [scripts](scripts)

### Стэк технологий

* OpenAPI
* Go
* GORM
* Gin
* PostgreSQL
* Docker & Docker Compose

### Запуск

```
git clone https://github.com/naumov-andrey/avito-intern-assignment.git
cd avito-intern-assignment/
```

Далее необходимо скопировать содержимое `.env.example` в `.env` и внести API ключ. В качестве источника информации 
о текущем курсе взят сервис [freecurrencyapi.net](https://freecurrencyapi.net)

Запуск Docker приложения:
```
docker compose up -d
```

Остановка Docker приложения:
```
docker compose down
```


### Примеры

#### GET /account/{userId}/balance

![img.png](docs/images/balance/get_balance_200.png)

---

![img.png](docs/images/balance/get_converted_balance_200.png)

---

![img.png](docs/images/balance/get_balance_400.png)

#### PUT /account/{userId}/balance

![img.png](docs/images/balance/update_balance_200_1.png)

---

![img.png](docs/images/balance/update_balance_200_2.png)

---

![img.png](docs/images/balance/update_balance_400_2.png)

---

![img.png](docs/images/balance/update_balance_400_1.png)

#### GET /account/{userId}/history

![img.png](docs/images/history/history_200_1.png)

---

![img.png](docs/images/history/history_200_2.png)

---

![img.png](docs/images/history/history_200_3.png)

---

![img.png](docs/images/history/history_200_4.png)

---

![img.png](docs/images/history/history_200_5.png)

---

![img.png](docs/images/history/history_400_1.png)

---

![img.png](docs/images/history/history_400_2.png)

---

![img.png](docs/images/history/history_400_3.png)

---

![img.png](docs/images/history/history_400_4.png)

#### POST /account/transfer

![img.png](docs/images/transfer/transfer_200_1.png)

---

![img.png](docs/images/transfer/transfer_200_2.png)

---

![img.png](docs/images/transfer/transfer_400_1.png)

---

![img.png](docs/images/transfer/transfer_400_2.png)

---

![img.png](docs/images/transfer/transfer_400_3.png)

---

