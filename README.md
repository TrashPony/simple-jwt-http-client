# simple-jwt-http-client

Для работы необходимо знать логин и пароль для авторизации в системе с которой вы работаете и система должна иметь jwt
авторизацию. Когда у нас есть логин и пароль мы можем создать клиент для создания запросов.

В каждом методе помимо аргументов запроса можно передать callback для логера: loggerCallBack func(string), если вместо
него передан nil то ошибки будут выводиться в консоль

## Пример использование

Создание клиента:

	package backend

	import (
		client "github.com/TrashPony/simple-jwt-http-client"
	)

	var backend *client.Backend

	func createClient() *client.Backend {

		newBackend, err := client.New("API_URL", "API_LOGIN", "API_PASS")
		if err != nil {
			panic(err)
		}

		return newBackend
	}

	func Client() *client.Backend {

		if backend == nil {
			backend = createClient()
		}

		return backend
	}

Использование клиента:

	randomData, err := backend.Client().RandomMethod(payload, nil)
