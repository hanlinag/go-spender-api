package app

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	auth "spender/v1/app/auth"
	controller "spender/v1/app/controller"
	models "spender/v1/app/models"
	configs "spender/v1/config"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// App initialize with predefined configuration
func (a *App) Initialize(config *configs.Config) {

	var dbURI = ""

	if configs.ISLOCAL {
		dbURI = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable Timezone=Asia/Rangoon",
			config.DB.Host,
			config.DB.Username,
			config.DB.Password,
			config.DB.Name,
			config.DB.Port)
	} else {
		dbURI = configs.DBURL
	}

	fmt.Println(dbURI)

	//db, err := gorm.Open(config.DB.Dialect, dbURI)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbURI,
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
	// Routing for handling the projec
	a.Get("/api/v1/employees", a.GetAllEmployees)
	a.Post("/api/v1/employees", a.CreateEmployee)
	a.Get("/api/v1/employees/{title}", a.GetEmployee)
	a.Put("/api/v1/employees/{title}", a.UpdateEmployee)
	a.Delete("/api/v1/employees/{title}", a.DeleteEmployee)
	a.Put("/api/v1/employees/{title}/disable", a.DisableEmployee)
	a.Put("/api/v1/employees/{title}/enable", a.EnableEmployee)

	//app config
	a.Get("/api/v1/app-config", a.GetAppConfig)
	a.Post("/api/v1/app-config", a.UpdateAppConfig)

	//user auth
	a.Post("/api/v1/auth/login", a.AuthLogin)
	a.Post("/api/v1/auth/signup", a.AuthSignUp)
	a.Post("/api/v1/auth/logout", a.Logout)

	//transactions
	a.Get("/api/v1/transactions", a.GetAllTransactions)
	a.Post("/api/v1/transaction", a.CreateTransaction)
	a.Get("/api/v1/transaction/{uuid}", a.GetSingleTransaction)
	a.Post("/api/v1/transaction/{uuid}", a.UpdateTransaction)
	a.Delete("/api/v1/transaction/{uuid}", a.DeleteTransaction)

	//wallet
	a.Get("/api/v1/wallets", a.GetAllWallets)
	a.Post("/api/v1/wallet", a.CreateWallet)
	a.Get("/api/v1/wallet/{uuid}", a.GetSingleWallet)
	a.Post("/api/v1/wallet/{uuid}", a.UpdateWallet)
	a.Delete("/api/v1/wallet/{uuid}", a.DeleteWallet)

	//Feedback
	a.Get("/api/v1/feedbacks", a.GetAllFeedbacks)
	a.Post("/api/v1/feedback", a.CreateFeedback)
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

//----------------------------------------
//App Config
func (a *App) GetAppConfig(w http.ResponseWriter, r *http.Request) {
	controller.GetAppConfig(a.DB, w)
}

func (a *App) UpdateAppConfig(w http.ResponseWriter, r *http.Request) {
	controller.UpdateAppConfig(a.DB, w, r)
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

func (a *App) GetSingleTransaction(w http.ResponseWriter, r *http.Request) {
	controller.GetTransaction(a.DB, w, r)
}

func (a *App) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	controller.UpdateTransaction(a.DB, w, r)
}

func (a *App) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	controller.DeleteTransaction(a.DB, w, r)
}

//--------------------------------------------------
//Wallets
func (a *App) GetAllWallets(w http.ResponseWriter, r *http.Request) {
	controller.GetAllWallets(a.DB, w, r)
}

func (a *App) CreateWallet(w http.ResponseWriter, r *http.Request) {
	controller.CreateWallet(a.DB, w, r)
}

func (a *App) GetSingleWallet(w http.ResponseWriter, r *http.Request) {
	controller.GetWallet(a.DB, w, r)
}

func (a *App) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	controller.UpdateWallet(a.DB, w, r)
}

func (a *App) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	controller.DeleteWallet(a.DB, w, r)
}


//--------------------------------------------------
//Wallets
func (a *App) GetAllFeedbacks(w http.ResponseWriter, r *http.Request) {
	controller.GetAllFeedbacks(a.DB, w)
}

func (a *App) CreateFeedback(w http.ResponseWriter, r *http.Request) {
	controller.CreateFeedback(a.DB, w, r)
}

//------------------------------------------------
// Run the app on it's router

//home page
var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func (a *App) Run(host string) {
	//Home Page
	//fs := http.FileServer(http.Dir("assets"))
	//a.Router.Handle("/assets/", http.StripPrefix("/assets/", fs))
	//a.Router.Handle("/assets/", http.FileServer(http.FS(contentStatic)))
	//a.Router.HandleFunc("/assets/", serveAssets)
	a.Router.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(host, a.Router))
}
