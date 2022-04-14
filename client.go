package backend

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const tokenName = "Authorization"
const prefixToken = "Bearer "

type (
	Backend struct {
		httpClient   *http.Client
		token        string
		refreshToken string
		jwt          jwt.MapClaims
		config       config
	}

	config struct {
		BaseUrl  string
		Login    string
		Password string
	}
)

// New создает новый http клиент
func New(baseUrl string, login string, password string) (*Backend, error) {
	newClient := &Backend{
		httpClient: http.DefaultClient,
		config: config{
			BaseUrl:  baseUrl,
			Login:    login,
			Password: password,
		},
	}

	return newClient, newClient.getToken()
}

// call метод принимает запрос который необходимо отправить к системе, подставляет заголовок и токен
func (api *Backend) call(method, url string, bodyRequest io.Reader) (json.RawMessage, *http.Response, error) {

	err := api.checkAuth()
	if err != nil {
		return nil, nil, err
	}

	return api.call2(method, url, bodyRequest, false)
}

func (api *Backend) call2(method, url string, bodyRequest io.Reader, forceExit bool) (json.RawMessage, *http.Response, error) {

	err := api.checkAuth()
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(method, url, bodyRequest)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set(tokenName, prefixToken+api.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	response, err := api.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode > 299 {

		if response.StatusCode == 401 && !forceExit {

			err := api.checkAuth()
			if err != nil {
				return body, nil, err
			}

			return api.call2(method, url, bodyRequest, true)
		}

		return body, response, errors.New("API: request failed\n" + string(body))
	}

	return body, response, nil
}

// checkAuth проверяем валидность токена
func (api *Backend) checkAuth() error {

	if api.token == "" {
		if err := api.getToken(); err != nil {
			return err
		}
	}

	// если токен не валиден то пытаемся обновить его, если и обновление не удалось пытаемся взять новый, если и это не удалось то выдаем ошибку
	if !api.jwt.VerifyExpiresAt(time.Now().Unix(), true) {
		if err := api.refreshTokenCall(); err != nil {
			if err := api.getToken(); err != nil {
				return err
			}
		}
	}

	return nil
}

// getToken запрашивает токен по пользователю и паролю
func (api *Backend) getToken() error {

	jsonValue, err := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		api.config.Login,
		api.config.Password,
	})

	request, err := http.NewRequest("POST", api.config.BaseUrl+"/api/login", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	return api.handleNewTokenMsg(request)
}

// refreshTokenCall обновляет токен по refresh_token
func (api *Backend) refreshTokenCall() error {

	jsonValue, err := json.Marshal(struct {
		RefreshToken string `json:"refresh_token"`
	}{
		api.refreshToken,
	})

	request, err := http.NewRequest("POST", api.config.BaseUrl+"/api/token/refresh", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	return api.handleNewTokenMsg(request)
}

// handleNewTokenMsg отправляет запрос на обновление или создание токена, парсит ответ и подставляет их в клиент
func (api *Backend) handleNewTokenMsg(request *http.Request) error {
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return errors.New("authorization error")
	}

	responseData := struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return err
	}

	api.token = responseData.Token
	api.refreshToken = responseData.RefreshToken
	api.jwt = jwt.MapClaims{}

	_, _, err = new(jwt.Parser).ParseUnverified(api.token, api.jwt)
	if err != nil {
		return err
	}

	return nil
}
