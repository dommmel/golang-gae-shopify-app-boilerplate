package app

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func init() {
	router := httprouter.New()
	router.GET("/install", serveInstall)
	router.GET("/admin", serveAdmin)
	router.GET("/app_proxy/", serveAppProxy)
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "views/home.html")
	})
	http.Handle("/", router)
}
