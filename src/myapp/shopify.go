package app

import (
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"gopkg.in/dommmel/go-shopify.v2"
	"html/template"
	"net/http"
	"os"
)

var app *shopify.App

func init() {
	var key, secret, hostname string

	key = os.Getenv("SHOPIFY_API_KEY")
	secret = os.Getenv("SHOPIFY_API_SECRET")
	hostname = os.Getenv("HOSTNAME")

	if appengine.IsDevAppServer() {
		key = os.Getenv("DEV_SHOPIFY_API_KEY")
		secret = os.Getenv("DEV_SHOPIFY_API_SECRET")
		hostname = os.Getenv("DEV_HOSTNAME")
	}

	redirect := hostname + "/install"

	app = &shopify.App{
		RedirectURI:     redirect,
		APIKey:          key,
		APISecret:       secret,
		IgnoreSignature: true,
	}

}

func serveAppProxy(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if app.AppProxySignatureOk(r.URL) {
		http.ServeFile(w, r, "views/app_proxy.html")
	} else {
		http.Error(w, "Unauthorized", 401)
	}
}

// initial page served when visited as embedded app inside Shopify admin
func serveAdmin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)
	params := r.URL.Query()
	var shopName string

	// signed request from Shopify?
	if app.AdminSignatureOk(r.URL) {
		log.Infof(c, "signed request!")
		shopName = params["shop"][0]
	} else {
		log.Infof(c, "no current_shop")
		// not logged in and not signed request
		http.Error(w, "Unauthorized", 401)
		return
	}

	if len(params["shop"]) != 1 {
		http.Error(w, "Expected 'shop' param", 400)
		log.Errorf(c, "No shop parameter")
		return
	}

	// shop, _ := session.Values["current_shop"].(string)

	// if we don't have an access token for the shop, obtain one now.
	// if _, ok := tokens[shop]; !ok {
	// 	http.Redirect(w, r, app.AuthorizeURL(shop, "read_themes,write_themes"), 302)
	// 	return
	// }

	// they're logged in

	type AdminVars struct {
		Shop   string
		APIKey string
	}
	v := AdminVars{Shop: shopName, APIKey: app.APIKey}

	t, _ := template.ParseFiles("views/admin.html")
	t.Execute(w, v)
}

func serveInstall(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)
	params := r.URL.Query()

	if len(params["error"]) == 1 {
		log.Errorf(c, "Install error: %s", params["error"])
	} else if len(params["code"]) == 1 {

		// auth callback from shopify
		if app.AdminSignatureOk(r.URL) != true {
			http.Error(w, "Invalid signature", 401)
			log.Errorf(c, "Invalid signature from Shopify")
			return
		}

		if len(params["shop"]) != 1 {
			http.Error(w, "Expected 'shop' param", 400)
			log.Errorf(c, "Invalid signature from Shopify")
			return
		}

		shop := params["shop"][0]
		token, _ := app.AccessToken(c, shop, params["code"][0])

		// persist this token
		log.Infof(c, "token is %s", token)

		// log in user
		// session := getSession(r)
		// session.Values["current_shop"] = shop
		// err := session.Save(r, w)
		// if err != nil {
		// 	panic(err)
		// }

		log.Infof(c, "logged in as %s, redirecting to admin", shop)

		http.Redirect(w, r, `/admin?shop=`+shop, 302)

	} else if len(params["install_shop"]) == 1 {
		// install request, redirect to Shopify
		shop := params["install_shop"][0]
		log.Infof(c, "starting oauth flow")

		http.Redirect(w, r, app.AuthorizeURL(shop, "read_themes,write_themes"), 302)
	}
}
