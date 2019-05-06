package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type userRepository interface {
	GetOne(userID string) (user *User, err error)
	Create(user User) (err error)
	Update(user User) (err error)
	Delete(userID string) (err error)
}

type UserService struct {
	userRepository userRepository
}

func NewUserService(userRepo userRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

func (s *UserService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleRead(w, r)

	case http.MethodPost:
		s.handleCreate(w, r)

	case http.MethodPut:
		s.handleUpdate(w, r)

	case http.MethodDelete:
		s.handleDelete(w, r)

	default:
		http.Error(
			w,
			"Invalid method",
			http.StatusMethodNotAllowed,
		)
	}
}

func (s *UserService) handleRead(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("id")
	if userID == "" {
		//Implement Get all if needed
	} else {
		user, err := s.userRepository.GetOne(userID)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusBadRequest,
			)
			return
		}

		userJSONBytes, err := json.Marshal(user)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.Write(userJSONBytes)
		w.WriteHeader(http.StatusOK)
		return
	}
}

func (s *UserService) handleCreate(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err = s.userRepository.Create(user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *UserService) handleUpdate(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	userID := r.FormValue("id")
	if userID == "" {
		http.Error(
			w,
			"User ID not present as a query param",
			http.StatusBadRequest,
		)
		return
	}

	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err = s.userRepository.Update(user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusNoContent,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UserService) handleDelete(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("id")
	if userID == "" {
		http.Error(
			w,
			"User ID not present as a query param",
			http.StatusBadRequest,
		)
		return
	}

	err := s.userRepository.Delete(userID)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusNoContent,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
