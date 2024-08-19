package app

import (
	"context"
	"fmt"
	"house-service/config"
	"house-service/infrastructure/db/adapter"
	infrastructure "house-service/infrastructure/db/repo"
	"house-service/internal/transport/handlers"
	mdware "house-service/internal/transport/middleware"
	"house-service/internal/usecase"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.Config) {
	// Подключение к базе данных
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db.Db)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("can't connect to PostgreSQL: %v", err.Error())
	}
	defer pool.Close()

	// Создание retry адаптера
	retryAdapter := adapter.NewPostgresRetryAdapter(pool, 3, time.Second*3)

	// Инициализация репозиториев
	houseRepo := infrastructure.NewPostgresHouseRepo(pool, retryAdapter)
	userRepo := infrastructure.NewPostgresUserRepo(pool, retryAdapter)
	apartmentRepo := infrastructure.NewPostgresApartmentRepo(pool, retryAdapter)

	// Инициализация usecases
	houseUsecase := usecase.NewHouseUsecase(houseRepo)
	userUsecase := usecase.NewUserUsecase(userRepo)
	flatUsecase := usecase.NewApartmentUsecase(apartmentRepo, houseRepo)

	// Инициализация хендлеров
	houseHandler := handlers.NewHouseHandler(houseUsecase, time.Duration(cfg.DbTimeoutSec)*time.Second)
	userHandler := handlers.NewUserHandler(userUsecase, time.Duration(cfg.DbTimeoutSec)*time.Second)
	flatHandler := handlers.NewApartmentHandler(flatUsecase, time.Duration(cfg.DbTimeoutSec)*time.Second)

	// Настройка роутера
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// Роутинг
	r.Post("/house/create", mdware.AuthMiddleware(mdware.AccessMiddleware(houseHandler.Create)))
	r.Get("/house/{id}", mdware.AuthMiddleware(houseHandler.GetApartmentsByID))
	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)
	r.Post("/flat/update", mdware.AuthMiddleware(mdware.AccessMiddleware(flatHandler.Update)))
	r.Post("/flat/create", mdware.AuthMiddleware(flatHandler.Create))

	fmt.Println("Server started at :8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		fmt.Println(err)
	}
}
