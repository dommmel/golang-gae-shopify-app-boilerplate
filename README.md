1. Create the standard google app engine app
2. Configure app.yaml and Procfile (ngrok subdomain)
3. Create two shopify apps. One for dev and one for production (with embedded sdk, set app url to "XXXX/admin" and callback url to "XXXX/install")
4. run "make deps" to install all dependencies
5. run "make serve" to start the local dev server. Run "make deploy" to deploy to gae
