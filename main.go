package main

import (
	//"TO-DO-LIST/db"
	"encoding/json"
	"log"
	"net/http"
	"todo/db"

	"github.com/google/uuid"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	}
}
func main() {

	db.Init()

	// create a new todo
	http.HandleFunc("/create", enableCORS(handleCreateTodo))

	//get
	http.HandleFunc("/getall", enableCORS(handleGetAllTodos))

	//update
	http.HandleFunc("/update", enableCORS(handleUpdateTodos))

	//delete
	http.HandleFunc("/delete", enableCORS(handleDeleteTodos))

	log.Printf("Starting server on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	//1.读取前端传来的参数
	params := map[string]string{}
	//2.解析参数
	err := json.NewDecoder((r.Body)).Decode(&params)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//3.处理参数
	name := params["name"]
	description := params["description"]

	//4.生成ID
	id := uuid.New().String()
	var newTodo db.Todo = db.Todo{
		ID:          id,
		Name:        name,
		Description: description,
		Completed:   false,
	}
	//5. 将新建的待办事项添加到待办事项列表中
	err = db.CreateTodo(newTodo)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//6. 返回响应
	w.WriteHeader(http.StatusOK)
}

func handleGetAllTodos(w http.ResponseWriter, r *http.Request) {
	//1. 返回所有待办事项
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//log.Println("handleGetAllTodos", db.Todos)

	todos, err := db.GetAllTodos()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func handleUpdateTodos(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{}

	err := json.NewDecoder((r.Body)).Decode(&params)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := params["id"]
	name := params["name"]
	description := params["description"]
	completed := params["completed"]

	err = db.UpdateTodo(db.Todo{
		ID:          id,
		Name:        name,
		Description: description,
		Completed:   completed == "true",
	})

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDeleteTodos(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{}
	err := json.NewDecoder((r.Body)).Decode(&params)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := params["id"]

	err = db.DeleteTodo(id)

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
