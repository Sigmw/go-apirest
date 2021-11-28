package habbo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	UniqueID string
	Name     string
	Motto    string
}

func GetUser(username string) (User, error) {
	req, err := http.Get(fmt.Sprintf("https://www.habbo.com.br/api/public/users?name=%v", username))
	if err != nil {
		return User{}, err
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return User{}, err
	}
	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		return User{}, err
	}
	if user.UniqueID == "" {
		return User{}, errors.New("nickname invalido")
	}
	return user, nil
}
