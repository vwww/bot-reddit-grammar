package main

import (
	"net/http"

	"google.golang.org/appengine"
)

func main() {
	appengine.Main()

	/*
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}

		http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	*/
}

func init() {
	http.HandleFunc("/_ah/warmup", func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()

		err := doInit(c)
		if err != nil {
			panic(err)
		}

		w.Write([]byte("ok"))
	})
	http.HandleFunc("/do/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("unknown action"))
	})
	http.HandleFunc("/do/post", func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()

		err := doPost(c, false)
		if err != nil {
			panic(err)
		}

		w.Write([]byte("ok"))
	})
	http.HandleFunc("/do/simulate", func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()

		err := doPost(c, true)
		if err != nil {
			panic(err)
		}

		w.Write([]byte("simulated"))
	})
	http.HandleFunc("/do/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// TODO save credentials
			r.FormValue("modhash")
			r.FormValue("cookie")

			w.Write([]byte("session info set"))
		}

		// not complete HTML document, but works in most browsers
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<form action="" method="POST">
	<input type="submit" value="Set!"><br>
	<input type="text" name="modhash" placeholder="modhash"><br>
	<input type="text" name="cookie" placeholder="cookie">
</form>
`))
	})
	http.HandleFunc("/do/comment", func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()

		err := doComment(
			c,
			r.FormValue("parent"),
			r.FormValue("text"),
		)
		if err != nil {
			panic(err)
		}

		w.Write([]byte("done"))
	})
}
