package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	auth "backend/auth"
	api "backend/interface"
	resourses "backend/serverless"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lpernett/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTION"},
		AllowedHeaders:   []string{"Content-Type", "X-TRADIFY-UID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ня :3"))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("Нету :("))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			auth.Login(w, r, db)
		})

		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			auth.Register(w, r, db)
		})
	})

	r.Route("/resourses", func(r chi.Router) {
		r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
			resourses.User(w, r, db)
		})
		r.Get("/trades", func(w http.ResponseWriter, r *http.Request) {
			resourses.Trades(w, r, db)
		})
		r.Get("/panels", func(w http.ResponseWriter, r *http.Request) {
			resourses.Panels(w, r, db)
		})
		r.Get("/notifications", func(w http.ResponseWriter, r *http.Request) {
			resourses.Notifications(w, r, db)
		})
	})

	r.Route("/interface", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				api.UserChangeBySocial(w, r, db)
			})
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
				api.UserChangeBySection(w, r, db)
			})
		})
		r.Route("/keys", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				api.KeyCreate(w, r, db)
			})
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
				api.KeyChangeTitle(w, r, db)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				api.KeyDelete(w, r, db)
			})
		})
		r.Route("/trades", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				api.TradesCreate(w, r, db)
			})
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
				api.TradeChangeBySection(w, r, db)
			})
			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				api.TradesDeleteByIds(w, r, db)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				api.TradesDeleteByKey(w, r, db)
			})
		})
		r.Route("/panels", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				api.PanelsCreate(w, r, db)
			})
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
				api.PanelsChangeWidgets(w, r, db)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				api.PanelsDelete(w, r, db)
			})
		})
		r.Route("/notifications", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				api.NotificationsCreate(w, r, db)
			})
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
				api.NotificationsRead(w, r, db)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				api.NotificationsDelete(w, r, db)
			})
		})
	})

	fmt.Println("норм")

	if err := http.ListenAndServe(":8000", r); err != nil {
		panic(err)
	}
}
