package auth

// // import "github.com/go-chi/chi/v5"

// // func Routes() *chi.Mux {
// // 	router := chi.NewRouter()

// // 	router.Post("/login", HandleLogin)
// // 	router.Post("/signup", HandleSignup)

// // 	return router
// // }

// package auth

// import (
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// // Handler holds your database pool (and later, maybe a logger or config)
// type Handler struct {
// 	DB *pgxpool.Pool
// }

// // NewHandler acts as a constructor to inject the dependencies
// func NewHandler(db *pgxpool.Pool) *Handler {
// 	return &Handler{DB: db}
// }

// // Routes returns the chi router for auth endpoints
// func (h *Handler) Routes() chi.Router {
// 	r := chi.NewRouter()

// 	// Since handleLogin is a method on *Handler, it automatically has access to h.DB!
// 	r.Get("/test", h.handleLoginTest)

// 	return r
// }

// func (h *Handler) handleLoginTest(w http.ResponseWriter, r *http.Request) {
// 	// Example of how you would use it:
// 	// err := h.DB.QueryRow(context.Background(), "SELECT ...")

// 	w.Write([]byte("Auth endpoint hitting the DB securely without global state!"))
// }
