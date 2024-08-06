package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"notify/cmd/config"
	"notify/internal/domain"
	"notify/internal/repo/employer"
	"notify/internal/storage"
	"notify/token_generator"
	"strconv"
	"strings"
)

func (a *Api) Init(ctx context.Context, config *config.HTTPConfig, db *storage.DB) {
	a.db = db
	router := chi.NewRouter()
	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов
	router.Post("/sign-in", a.TokenAuth)
	router.Post("/register", a.register)
	router.Post("/subscribe/{publisherID}", a.subscribe)
	router.Post("/unsubscribe/{publisherID}", a.unsubscribe)
	router.Get("/list", a.getListSubscribe)
	router.Get("/empinfo", a.getAllEmp)
	router.Get("/empID/{userID}", a.getEmpID)
	fmt.Printf("Запуск сервера на http://%s:", config.Address)
	err := http.ListenAndServe(fmt.Sprintf("%s", config.Address), router)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-ctx.Done():
			fmt.Println("exit from api server")
			return
		default:

		}
	}
}

func (a *Api) getListSubscribe(w http.ResponseWriter, r *http.Request) {
	var getFromDB *domain.ResponseEmployee
	var subscribeID int
	accessToken, err := headerStr(w, r)
	if err != nil {
		log.Printf("Error get header string: %v", err)
		return
	}
	token := token_generator.New([]byte("secret-key"), a.db)
	subscribeID, err = token.ParseToken(accessToken)
	if err != nil {
		log.Printf("Error parse token: %v", err)
		return
	}
	userRepo := employer.New(a.db)
	getList, err := userRepo.List(context.Background(), subscribeID)
	if err != nil {
		log.Printf("Error getting employee data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, id := range getList {
		getFromDB, err = userRepo.GetID(context.Background(), id)
		if err != nil {
			log.Printf("Error getting employee data: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonEmpInfo, err := json.Marshal(*getFromDB)
		if err != nil {
			log.Printf("Error marshalling employee data to JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonEmpInfo)
	}
}

func (a *Api) unsubscribe(w http.ResponseWriter, r *http.Request) {
	emp, err := a.getValues(w, r)
	if err != nil {
		log.Printf("Error getting values for unsibscribe: %v", err)
		return
	}
	userRepo := employer.New(a.db)
	err = userRepo.Unsub(context.Background(), emp.subscribeID, emp.publisherID)
	if err != nil {
		log.Printf("Error Unsubscribe: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Successfully unsubscribe!")
}

func (a *Api) subscribe(w http.ResponseWriter, r *http.Request) {
	emp, err := a.getValues(w, r)
	if err != nil {
		log.Printf("Error getting values subscribe: %v", err)
		return
	}
	userRepo := employer.New(a.db)
	err = userRepo.Sub(context.Background(), emp.subscribeID, emp.publisherID)
	if err != nil {
		log.Printf("Error subscribe: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Successfully subscribe!")
}

func (a *Api) TokenAuth(w http.ResponseWriter, r *http.Request) {
	var user domain.Employee
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	token := token_generator.New([]byte("secret-key"), a.db)
	tokenGenerate, err := token.Generate(user.Username)
	//tokenGenerate, err := a.auth.Generate(user.Username)
	if err != nil {
		log.Printf("Error generate token from api: %v", err)
		http.Error(w, "Error generate token from api", http.StatusBadRequest)
		return
	}
	tokenString, err := tokenGenerate.String()
	if err != nil {
		log.Printf("Error tokenString from api: %v", err)
		http.Error(w, "Error tokenString from api", http.StatusBadRequest)
		return
	}
	w.Header().Set("Authorization", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (a *Api) register(w http.ResponseWriter, r *http.Request) {
	var emp domain.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	userRepo := employer.New(a.db)
	err = userRepo.Create(context.Background(), emp)
	//err = a.userRepo.Create(context.Background(), emp) //todo
	if err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

func (a *Api) getAllEmp(w http.ResponseWriter, r *http.Request) {
	userRepo := employer.New(a.db)
	getFromDB, err := userRepo.Get(context.Background())
	//getFromDB, err := a.userRepo.Get(context.Background()) // todo
	if err != nil {
		log.Printf("Error getting employee data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if getFromDB == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonEmpInfo, err := json.Marshal(*getFromDB)
	if err != nil {
		log.Printf("Error marshalling employee data to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonEmpInfo)
}

func (a *Api) getEmpID(w http.ResponseWriter, r *http.Request) {
	empIDstr := chi.URLParam(r, "userID")
	empIDstr = strings.Trim(empIDstr, "{}")
	id, err := strconv.Atoi(empIDstr)
	if err != nil {
		log.Printf("strconv.Atoi: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	userRepo := employer.New(a.db)
	getFromDB, err := userRepo.GetID(context.Background(), id)
	//getFromDB, err := a.userRepo.GetID(context.Background(), id) //todo
	if err != nil {
		log.Printf("Error getting employee data: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if getFromDB == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonEmpInfo, err := json.Marshal(*getFromDB)
	if err != nil {
		log.Printf("Error marshalling employee data to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonEmpInfo)
}

func NewApi(ctx context.Context, httpConfig *config.HTTPConfig, db *storage.DB) *Api {
	api := Api{}
	newCtx := ctx
	api.Init(newCtx, httpConfig, db)
	return &api
}

type Api struct {
	db        *storage.DB
	auth      tokenGenerator
	generator *token_generator.TokenGenerator

	userRepo userRepository
	user     User
}

type tokenGenerator interface {
	Generate(username string) (*token_generator.Token, error)
	ParseToken(username string) (string, error)
}

type userRepository interface {
	List(ctx context.Context, id int) (*[]domain.ResponseEmployee, error)
	Create(ctx context.Context, employee domain.Employee) error
	Unsub(ctx context.Context, sub, pub int) error
	Sub(ctx context.Context, sub, pub int) error
	Get(ctx context.Context) (*[]domain.ResponseEmployee, error)
	GetID(ctx context.Context, id int) (*domain.ResponseEmployee, error)
}
