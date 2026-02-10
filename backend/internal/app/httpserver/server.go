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

	"github.com/SHILOP0P/Yardly/backend/internal/admin"

	"github.com/SHILOP0P/Yardly/backend/internal/auth"
	"github.com/SHILOP0P/Yardly/backend/internal/booking"
	"github.com/SHILOP0P/Yardly/backend/internal/item"
	"github.com/SHILOP0P/Yardly/backend/internal/favorite"
	"github.com/SHILOP0P/Yardly/backend/internal/httpx"
)



func New(port string, pool *pgxpool.Pool,itemsRepo *itempg.Repo, bookingRepo *bookingpg.Repo, userRepo *userpg.Repo, refreshesRepo *auth.RefreshRepo, favoriteRepo favorite.Repo, adminRepo admin.Repo, jwtSvc *auth.JWT, refreshTTL time.Duration) *http.Server{
	mux := http.NewServeMux()

	RegisterBaseRotes(mux)

	authMw := auth.Middleware(jwtSvc)

	protectedChain := func (h http.Handler) http.Handler{
		return authMw(auth.RequireNotBanned(h))
	}


	adminChain := func(h http.Handler) http.Handler {
		return authMw(auth.RequireNotBanned(auth.RequireAdmin(h)))
	}

	// superAdminChain := func(h http.Handler) http.Handler {
	// 	return authMw(auth.RequireSuperAdmin(authUsers)(h))
	// }

	item.RegisterRoutes(mux, itemsRepo, protectedChain)
	booking.RegisterRoutes(mux, bookingRepo, itemsRepo, protectedChain)
	user.RegisterRoutes(mux, authMw, userRepo, jwtSvc)
	auth.RegisterRoutes(mux, jwtSvc, refreshesRepo, refreshTTL, userRepo, authMw)
	favorite.RegisterRoutes(mux, favoriteRepo, authMw)
	admin.RegisterRoutes(mux, adminRepo, adminChain)

	return &http.Server{
		Addr: ":" + port,
		Handler: httpx.CORS(mux),
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

