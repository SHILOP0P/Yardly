package httpserver

import (
	"fmt"
	"net/http"
	"time"

	bookingpg "github.com/SHILOP0P/Yardly/backend/internal/booking/pgrepo"
	itempg "github.com/SHILOP0P/Yardly/backend/internal/item/pgrepo"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/SHILOP0P/Yardly/backend/internal/user"
	userpg "github.com/SHILOP0P/Yardly/backend/internal/user/pgrepo"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
	"github.com/SHILOP0P/Yardly/backend/internal/booking"
	"github.com/SHILOP0P/Yardly/backend/internal/item"
)

func New(port string, pool *pgxpool.Pool,itemsRepo *itempg.Repo, bookingRepo *bookingpg.Repo, userRepo *userpg.Repo, refreshesRepo *auth.RefreshRepo, jwtSvc *auth.JWT, refreshTTL time.Duration) *http.Server{
	mux := http.NewServeMux()

	RegisterBaseRotes(mux)

	authMw := auth.Middleware(jwtSvc)

	item.RegisterRoutes(mux, itemsRepo, authMw)

	booking.RegisterRoutes(mux, bookingRepo, itemsRepo, authMw)

	user.RegisterRoutes(mux, authMw, userRepo, jwtSvc)

	auth.RegisterRoutes(mux, jwtSvc, refreshesRepo, refreshTTL, userRepo, authMw)

	return &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
}

func RegisterBaseRotes(mux *http.ServeMux){
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Yardly")
	})
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})
}