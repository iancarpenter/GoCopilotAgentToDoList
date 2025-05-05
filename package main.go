package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Todo struct {
	ID   int
	Task string
	Done bool
}

var (
	todos  []Todo
	nextID = 1
	mu     sync.Mutex
	dbFile = "todos.json"
)

// loadTodos loads the todo list from the JSON file and sets the next available ID.
func loadTodos() {
	file, err := os.Open(dbFile)
	if err != nil {
		return // No file yet
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	json.Unmarshal(data, &todos)
	// Set nextID
	for _, t := range todos {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
}

// saveTodos saves the current todo list to the JSON file.
func saveTodos() {
	data, _ := json.MarshalIndent(todos, "", "  ")
	ioutil.WriteFile(dbFile, data, 0644)
}

// main is the entry point of the application. It sets up HTTP handlers and starts the server.
func main() {
	loadTodos()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.HandleFunc("/list", listHandler)
	http.ListenAndServe(":8080", nil)
}

// indexHandler serves the main HTML page for the todo app.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Go Todo App</title>
		<script>
		async function addTodo() {
			const task = document.getElementById('task').value;
			await fetch('/add?task=' + encodeURIComponent(task));
			document.getElementById('task').value = '';
			loadTodos();
		}
		async function deleteTodo(id) {
			await fetch('/delete?id=' + id);
			loadTodos();
		}
		async function toggleTodo(id) {
			await fetch('/toggle?id=' + id);
			loadTodos();
		}
		async function loadTodos() {
			const res = await fetch('/list');
			const todos = await res.json();
			let html = '';
			for (const t of todos) {
				html += '<li>' +
					'<input type="checkbox" ' + (t.Done ? 'checked' : '') + ' onclick="toggleTodo(' + t.ID + ')">' +
					(t.Done ? '<s>' + t.Task + '</s>' : t.Task) +
					' <button onclick="deleteTodo(' + t.ID + ')">Delete</button></li>';
			}
			document.getElementById('todos').innerHTML = html;
		}
		window.onload = loadTodos;
		</script>
	</head>
	<body>
		<h1>Go Todo App</h1>
		<input id="task" type="text" placeholder="New todo">
		<button onclick="addTodo()">Add</button>
		<ul id="todos"></ul>
	</body>
	</html>
	`
	template.Must(template.New("index").Parse(tmpl)).Execute(w, nil)
}

// addHandler handles adding a new todo item via the /add endpoint.
func addHandler(w http.ResponseWriter, r *http.Request) {
	task := r.URL.Query().Get("task")
	if task == "" {
		http.Error(w, "Task required", http.StatusBadRequest)
		return
	}
	mu.Lock()
	todos = append(todos, Todo{ID: nextID, Task: task})
	nextID++
	saveTodos()
	mu.Unlock()
}

// deleteHandler handles deleting a todo item by ID via the /delete endpoint.
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	mu.Lock()
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			saveTodos()
			break
		}
	}
	mu.Unlock()
}

// toggleHandler toggles the completion status of a todo item by ID via the /toggle endpoint.
func toggleHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	mu.Lock()
	for i, t := range todos {
		if t.ID == id {
			todos[i].Done = !todos[i].Done
			saveTodos()
			break
		}
	}
	mu.Unlock()
}

// listHandler returns the list of todos as a JSON response via the /list endpoint.
func listHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
