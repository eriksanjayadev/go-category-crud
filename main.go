package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori Makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori Minuman"},
}

func JSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"message": message,
	})
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func getCategories(w http.ResponseWriter) {
	JSON(w, http.StatusOK, categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var newCategory Category

	err := json.NewDecoder(r.Body).Decode(&newCategory)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	JSON(w, http.StatusCreated, newCategory)
}

func getCategoryById(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// konversi id ke int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	// loop categories and get data by id
	for _, data := range categories {
		if data.ID == id {
			JSON(w, http.StatusAccepted, data)
			return
		}
	}

	JSONError(w, http.StatusNotFound, "Category Not Found")
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// konversi id ke integer
	id, err := strconv.Atoi(idStr)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid Category ID")
	}

	// get data dari request
	var updateCategory Category

	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid Request")
	}

	fmt.Println(err)

	// loop categories, get by id
	for i := range categories {
		if categories[i].ID == id {
			updateCategory.ID = id
			categories[i] = updateCategory

			JSON(w, http.StatusOK, updateCategory)
			return
		}
	}

	JSONError(w, http.StatusNotFound, "Category Not Found")
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")

	// konversi id ke integer
	id, err := strconv.Atoi(idStr)

	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid Category ID")
	}

	var deleteCategory Category

	// loop category, get category by id
	for i, data := range categories {
		if data.ID == id {

			// ambil data yang didelete
			deleteCategory = categories[i]

			// buat slice baru, dengan data sebelum dan sesudah index
			categories = append(categories[:i], categories[i+1:]...)

			JSON(w, http.StatusOK, deleteCategory)
			return
		}
	}

	JSONError(w, http.StatusNotFound, "Category Not Found")
}

func main() {

	// GET http://localhost:8080/api/categories/{id}
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryById(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}
	})

	// FindAll	:GET http://localhost:8080/api/categories
	// Create	  :POST http://localhost:8080/api/categories
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategories(w)
		} else if r.Method == "POST" {
			createCategory(w, r)
		}
	})

	// GET http://localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("Halo"))

		w.Header().Set("Content-Type", "Application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "Api Running",
		})
	})

	fmt.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server failed to start")
	}
}
