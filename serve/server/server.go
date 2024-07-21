package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"serve/server/storage"
)

type Item storage.Item

func (it Item) getInfo() string {
	return fmt.Sprintf("ItemLoginz: %v , ItemMoneu: %v, ItemScore: %v", it.Login, it.Money, it.Score)
}

func GETMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Ожидался GET запрос")
		return
	}
	if r.Method == http.MethodGet {
		fmt.Println("Был получен GET")
		fmt.Fprintln(w, "Был получен GET!")
	}
}

func POSTMethod(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "SecPage Working")
		if r.Method != http.MethodPost {
			fmt.Println("Ожидался POST запрос")
			return
		}
		if r.Method == http.MethodPost {
			item := storage.Item{}
			err := json.NewDecoder(r.Body).Decode(&item)
			if err != nil {
				fmt.Println("Проблема с декодированием")
				return
			}
			log.Printf("Полученный JSON: %+v", item)

			err = s.PostStorageFunc(&item)
			if err != nil {
				fmt.Println("Проблема с записью в БД")
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}

func GetCheck(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			fmt.Println("Ожидался GET запрос")
			return
		}
		if r.Method == http.MethodGet {
			fmt.Println("прошло")
			items, err := s.GetAllItems()
			if err != nil {
				http.Error(w, "Error fetching items", http.StatusInternalServerError)
				return
			}

			fmt.Println(items)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(items); err != nil {
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
			}
		}
	}
}

func MainHandle(store *storage.Storage) {
	http.HandleFunc("/", GETMethod)
	http.HandleFunc("/post", POSTMethod(store))
	http.HandleFunc("/check", GetCheck(store))
	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println("And on http://localhost:8080/second")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main() {
	store, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("Error starting DB:", err)
	}
	store.CreateTables()

	MainHandle(store)
	var excaper string
	fmt.Println("Please enter C to break server:")
	for true {
		fmt.Scan(&excaper)
		if excaper == "C" {
			break
		}
	}

	defer store.Db.Close()

}
