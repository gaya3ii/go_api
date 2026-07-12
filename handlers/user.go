package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/gayathriad/go_api/models"
	"github.com/gayathriad/go_api/repository"
)

// max concurrent DB writes for a bulk request

const bulkWorkerLimit = 10

type bulkUserResult struct {
	Index int         `json:"index"`
	User  models.User `json:"user,omitempty"`
	Error string      `json:"error,omitempty"`
}

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// writeJSON marshals data and writes it to the response, along with the
// given status code. If marshaling fails, it falls back to a 500 error.
func writeJSON(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to encode response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

// readJSON reads and unmarshals the request body into v. It writes a 400
// error response and returns false if reading or parsing fails.
func readJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"failed to read body"}`))
		return false
	}
	if err := json.Unmarshal(body, v); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid body"}`))
		return false
	}
	return true
}

// get all users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// create user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if !readJSON(w, r, &user) {
		return
	}
	created, err := h.repo.Create(r.Context(), user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// get user by id
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}
	user, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// update user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}
	// check user exists
	existing, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	// decode new values
	if !readJSON(w, r, &existing) {
		return
	}
	updated, err := h.repo.Update(r.Context(), existing)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// create many users concurrently via a bounded worker pool
func (h *UserHandler) CreateUsersBulk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var users []models.User
	if !readJSON(w, r, &users) {
		return
	}
	if len(users) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "users array is required"})
		return
	}

	type job struct {
		index int
		user  models.User
	}

	jobs := make(chan job, len(users))
	results := make(chan bulkUserResult, len(users))

	workers := bulkWorkerLimit
	if len(users) < workers {
		workers = len(users)
	}

	ctx := r.Context()

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				created, err := h.repo.Create(ctx, j.user)
				res := bulkUserResult{Index: j.index, User: created}
				if err != nil {
					res.Error = err.Error()
				}
				results <- res
			}
		}()
	}

	for i, u := range users {
		jobs <- job{index: i, user: u}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	ordered := make([]bulkUserResult, len(users))
	for res := range results {
		ordered[res.Index] = res
	}

	writeJSON(w, http.StatusMultiStatus, ordered)
}

// delete user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}
	if err := h.repo.Delete(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
