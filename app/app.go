package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"

	controller "spender/v1/app/controller"
	models "spender/v1/app/models"
	configs "spender/v1/config"
	auth "spender/v1/app/auth"

)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// App initialize with predefined configuration
func (a *App) Initialize(config *configs.Config) {
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable Timezone=Asia/Rangoon",
		config.DB.Host,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Port)

	fmt.Println(dbURI)

	//db, err := gorm.Open(config.DB.Dialect, dbURI)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: configs.DBURL , //staging
		//cDSN: dbURI, //local
		PreferSimpleProtocol: true, // disables implicit prepared statement usage. By default pgx automatically uses the extended protocol
	  }), &gorm.Config{})

	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = models.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

// Set all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/v1/employees", a.GetAllEmployees)
	a.Post("/v1/employees", a.CreateEmployee)
	a.Get("/employees/{title}", a.GetEmployee)
	a.Put("/employees/{title}", a.UpdateEmployee)
	a.Delete("/employees/{title}", a.DeleteEmployee)
	a.Put("/employees/{title}/disable", a.DisableEmployee)
	a.Put("/employees/{title}/enable", a.EnableEmployee)

	//user auth
	a.Post("/v1/auth/login", a.AuthLogin)
	a.Post("/v1/auth/signup", a.AuthSignUp)
	a.Post("/v1/auth/logout", a.Logout)

	//transactions
	a.Get("/v1/transactions", a.GetAllTransactions)
	a.Post("/v1/transaction", a.CreateTransaction)
}

// Wrap the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	//a.Router.HandleFunc(path, f).Methods("GET")
	a.Router.HandleFunc(path, auth.CheckAuth(a.DB, f)).Methods("GET")

}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	
	a.Router.HandleFunc(path, auth.CheckAuth(a.DB, f)).Methods("POST")
	
}

// Wrap the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, auth.CheckAuth(a.DB, f)).Methods("PUT")
}

// Wrap the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, auth.CheckAuth(a.DB, f)).Methods("DELETE")
}

//--------------------------------------------
//Auth Login
func (a *App) AuthLogin(w http.ResponseWriter, r *http.Request) {
	controller.Login(a.DB, w, r)
}

func (a *App) AuthSignUp(w http.ResponseWriter, r *http.Request) {
	controller.SignUp(a.DB, w, r)
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	controller.Logout(a.DB, w, r)
}

//-------------------------------------------
// Handlers to manage Employee Data
func (a *App) GetAllEmployees(w http.ResponseWriter, r *http.Request) {
	controller.GetAllEmployees(a.DB, w, r)
}

func (a *App) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	controller.CreateEmployee(a.DB, w, r)
}

func (a *App) GetEmployee(w http.ResponseWriter, r *http.Request) {
	controller.GetEmployee(a.DB, w, r)
}

func (a *App) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	controller.UpdateEmployee(a.DB, w, r)
}

func (a *App) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	controller.DeleteEmployee(a.DB, w, r)
}

func (a *App) DisableEmployee(w http.ResponseWriter, r *http.Request) {
	controller.DisableEmployee(a.DB, w, r)
}

func (a *App) EnableEmployee(w http.ResponseWriter, r *http.Request) {
	controller.EnableEmployee(a.DB, w, r)
}


//--------------------------------------------------
//Transactions
func (a *App) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	controller.GetAllTransactions(a.DB, w, r)
}

func (a *App) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	controller.CreateTransaction(a.DB, w, r)
}


//------------------------------------------------
// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
