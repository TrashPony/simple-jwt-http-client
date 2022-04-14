package backend

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {

	backend, err := New("http://localhost", "admin@gmail.com", "12345678")
	if err != nil {
		t.Error(err)
	}

	t.Run("get_token", func(t2 *testing.T) {
		err = backend.getToken()
		if err != nil {
			t2.Error(err)
		}
	})

	t.Run("refresh_token", func(t2 *testing.T) {
		err := backend.refreshTokenCall()
		if err != nil {
			t2.Error(err)
		}
	})

	t.Run("example_post", func(t2 *testing.T) {
		data, err := backend.ExamplePostRequest("1", nil)
		if err != nil {
			t2.Error(err)
		}

		fmt.Println(data)
	})

	t.Run("example_get", func(t2 *testing.T) {
		data, err := backend.ExampleGetRequest("1", nil)
		if err != nil {
			t2.Error(err)
		}

		fmt.Println(data)
	})
}
