package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Item struct {
	Login string  `json:"login"`
	Money int     `json:"money"`
	Score float64 `json:"score"`
}

var pathMainGET = "http://localhost:8080"
var pathMainPOST = "http://localhost:8080/post"
var pathcheckGET = "http://localhost:8080/check"

func GetFunc() error {
	response, err := http.Get(pathMainGET)
	if err != nil {
		log.Fatalf("Could not make GET request: %s\n", err.Error())
	}
	defer response.Body.Close()
	return nil
}

func PostFunc() error {
	fmt.Println("Введите данные о предмете {login,money,score}")
	item := Item{}
	_, err := fmt.Scanln(&item.Login, &item.Money, &item.Score)
	if err != nil {
		return err
	}

	fmt.Println(item.Login, item.Money, item.Score)
	itemJSON, err := json.Marshal(item)
	if err != nil {
		fmt.Println("Ошибка сериализации JSON:", err)
		return err
	}
	response, err := http.Post(pathMainPOST, "application/json", bytes.NewBuffer(itemJSON))
	if err != nil {
		log.Fatalf("Could not make POST request: %s\n", err.Error())
	}
	defer response.Body.Close()
	return nil
}

func GetInfo() error {
	response, err := http.Get(pathcheckGET)
	if err != nil {
		log.Fatalf("Could not make checkGetReq request: %s\n", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("Получен ненормальный статус ответа:", response.Status)
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return err
	}

	fmt.Println(response.Body)
	var items []Item
	err = json.Unmarshal(body, &items)
	if err != nil {
		fmt.Println("Ошибка при расшифровке JSON:", err)
		return err
	}

	for _, item := range items {
		fmt.Printf("Login: %s, Money: %d, Score: %.2f\n", item.Login, item.Money, item.Score)
	}
	defer response.Body.Close()
	return nil
}

func main() {

	fmt.Println("Клиент начал работу")
loop:
	for true {
		fmt.Println("Введите 1, чтобы выйти из программы и завершить работу")
		fmt.Println("Введите 2, чтобы выполнить GET запрос")
		fmt.Println("Введите 3, чтобы выполнить POST запрос")
		fmt.Println("Введите 4, чтобы получить всю информацию из таблицы запрос")
		var a int
		_, err := fmt.Scan(&a)
		if err != nil {
			fmt.Println("Неверный ввод")
			continue
		}
		fmt.Println("Вы ввели:", a)
		switch a {
		case 1:
			break loop
		case 2:
			GetFunc()
		case 3:
			PostFunc()
		case 4:
			GetInfo()
		default:
			fmt.Println("нет такого варианта")
		}
	}
	fmt.Println("Клиент завершил работу")
}
