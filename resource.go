package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/errors"
	"github.com/go-chi/chi"
)

type ctx int

const (
	bodyKey ctx = 0
	loadKey ctx = 1
)

type Resource struct {
	Repo grimoire.Repo
}

func (resource Resource) Index(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	resource.Repo.From(TodoTable).Order(c.Asc("order")).MustAll(&todos)

	resource.send(w, todos, 200)
}

func (resource Resource) Create(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(bodyKey).(map[string]interface{})

	ch := CreateTodo(params)
	if ch.Error() != nil {
		resource.send(w, ch.Errors(), 422)
		return
	}

	var todo Todo
	err := resource.Repo.From(TodoTable).Insert(&todo, ch)
	if err != nil {
		if e, ok := err.(errors.Error); ok && e.Kind() != errors.Unexpected {
			resource.send(w, e, 422)
			return
		}

		panic(err)
	}

	w.Header().Set("Location", fmt.Sprint("/", todo.ID))
	resource.send(w, todo, 201)
}

func (resource Resource) Show(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value(loadKey).(Todo)
	resource.send(w, todo, 201)
}

func (resource Resource) Update(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value(loadKey).(Todo)
	params := r.Context().Value(bodyKey).(map[string]interface{})

	ch := ChangeTodo(todo, params)
	if ch.Error() != nil {
		resource.send(w, ch.Errors(), 422)
		return
	}

	err := resource.Repo.From(TodoTable).Find(todo.ID).Update(&todo, ch)
	if err != nil {
		if e, ok := err.(errors.Error); ok && e.Kind() != errors.Unexpected {
			resource.send(w, e, 422)
			return
		}

		panic(err)
	}

	resource.send(w, todo, 200)
}

func (resource Resource) Delete(w http.ResponseWriter, r *http.Request) {
	todo := r.Context().Value(loadKey).(Todo)

	resource.Repo.From(TodoTable).Find(todo.ID).Delete()
	resource.send(w, nil, 204)
}

func (resource Resource) Clear(w http.ResponseWriter, r *http.Request) {
	resource.Repo.From(TodoTable).Delete()
	resource.send(w, nil, 204)
}

func (resource Resource) BodyParser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&body)

		if err != nil {
			resource.send(w, nil, 400)
			return
		}

		ctx := context.WithValue(r.Context(), bodyKey, body)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (resource Resource) Load(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "ID"))

		var todo Todo
		err := resource.Repo.From(TodoTable).Find(id).One(&todo)
		if err != nil {
			if e, ok := err.(errors.Error); ok && e.Kind() == errors.NotFound {
				resource.send(w, e, 404)
				return
			}
			panic(err)
		}

		ctx := context.WithValue(r.Context(), loadKey, todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (resource Resource) send(w http.ResponseWriter, body interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
