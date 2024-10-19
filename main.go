package main

import (
	//"database/sql"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abhi3832/jwt_auth/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"// special import kind of thing
)

type apiConfig struct{
	DB * database.Queries
}
// we do so, in order to have to all the database functions
// the handler function still remains the same , we just get access to all the database function in addition


func main(){

	// fetch the port number from the env file
	godotenv.Load(".env")
	portstring := os.Getenv("PORT")

	if portstring == ""{
		log.Fatal("Port is not found in the env !")
	}else{
		fmt.Println("Server Running on Port :", portstring)
	}

	db_url := os.Getenv("DB_URL")

	if db_url == ""{
		log.Fatal("Error Fetching DB URL from the env !")
	}

	conn, err := sql.Open("postgres", db_url)// connect to database

	if err != nil{
		log.Fatal("Can't Connect to DB !")
	}

	apicfg := apiConfig{
		DB : database.New(conn),
	}


	// create and configure the main router
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: 	[]string {"https://*", "http://*"},
		AllowedMethods: 	[]string {"GET", "POST","PUT","DELETE","OPTIONS"},
		AllowedHeaders: 	[]string {"*"},
		ExposedHeaders: 	[]string {"Link"},
		AllowCredentials: false,
		MaxAge: 			300,
	}))

	// create a new sub-router v1
	v1Router := chi.NewRouter()
	
	// requests
	v1Router.Get("/health", serverCheckHandler)
	v1Router.Post("/create_user", apicfg.userSignupHandler)
	v1Router.Post("/login", apicfg.userLoginHandler)
	v1Router.Get("/getuser",apicfg.middlewareAuthByToken(apicfg.getUser))
	v1Router.Post("/refreshtokens", apicfg.refreshTokens)

	














	
	// mount all the routes to the main router
	router.Mount("/v1", v1Router)


	// bind the port and router to the server
	srv := &http.Server{
		Handler: router,
		Addr: ":"+portstring,
	}

	log.Printf("Server Started on Port : %v", portstring)
	err1 := srv.ListenAndServe()
	if err1 != nil{
		log.Fatal(err1)
	}

	

}