package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"test_2/server/authentication"
	controller "test_2/server/controllers"
	"time"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		t1 := time.Now()

		next.ServeHTTP(w, r)

		t2 := time.Now()
		log.Printf(" [%s] %s %s", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				log.Print(string(debug.Stack()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func wrapHandler(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func noDirListingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		f, err := os.Stat(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		if f != nil && f.IsDir() {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func renderERROR(w http.ResponseWriter) {
	e := map[string]interface{}{
		"code":    http.StatusUnauthorized,
		"message": "Un-authorized Access",
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		log.Printf("ERROR: renderJson - %q\n", err)
	}
}

func tokenCustomerAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authMethod := authentication.GetAuthMethod()

		token, err := request.ParseFromRequest(r, request.OAuth2Extractor, authMethod)
		log.Println(token, err)
		if err == nil && token.Valid {
			c := token.Claims.(jwt.MapClaims)
			id := c["id"].(string)

			if !bson.IsObjectIdHex(id) {
				renderERROR(w)
				return
			}
			// user := struct {
			// 	ID	string
			// 	IsCustomer
			// }

			ctx := context.WithValue(r.Context(), "loggedInUserId", id)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		renderERROR(w)
	})
}

func HTTPRouteConfig() *httprouter.Router {

	router := httprouter.New()
	handler := alice.New(loggingHandler, recoverHandler)
	userhandler := alice.New(loggingHandler, recoverHandler, tokenCustomerAuthentication)

	router.GET("/", wrapHandler(http.FileServer(http.Dir("public/app"))))
	router.Handler("GET", "/public/*filepath", noDirListingHandler(http.StripPrefix("/public/", http.FileServer(http.Dir("public")))))

	router.POST("/registerUser", wrapHandler(handler.ThenFunc(controller.RegisterUser)))
	router.POST("/authenticateUser", wrapHandler(handler.ThenFunc(controller.AuthenticateUser)))

	router.GET("/getContents", wrapHandler(handler.ThenFunc(controller.Contents)))

	router.GET("/content/:id", wrapHandler(userhandler.ThenFunc(controller.GetContent)))
	router.POST("/createContent", wrapHandler(userhandler.ThenFunc(controller.CreateContent)))
	router.PUT("/editContent/:id", wrapHandler(userhandler.ThenFunc(controller.UpdateContent)))
	router.DELETE("/deleteContent/:id", wrapHandler(userhandler.ThenFunc(controller.DeleteContent)))

	router.PUT("/content/:id/vote", wrapHandler(userhandler.ThenFunc(controller.Vote)))

	router.PUT("/content/:id/comment", wrapHandler(userhandler.ThenFunc(controller.Comment)))
	router.PUT("/comment/:id/reply", wrapHandler(userhandler.ThenFunc(controller.Reply)))
	router.GET("/content/:id/comments", wrapHandler(userhandler.ThenFunc(controller.Comments)))

	return router
}
