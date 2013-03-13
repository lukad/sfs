package main

import (
	auth "bitbucket.org/taruti/http_auth"
	"fmt"
	flag "github.com/ogier/pflag"
	"net/http"
)

var (
	flagAddress  = flag.StringP("listen", "l", ":8080", "Listen address")
	flagDigest   = flag.BoolP("digest", "d", false, "Use digest access authentication")
	flagUser     = flag.StringP("user", "u", "", "Username for authentication")
	flagPassword = flag.StringP("password", "p", "", "Password for authentication")
	flagRoot     = flag.StringP("root", "r", "./", "Root Directory for the file server")
)

func makeBasicAuthHandleFunc(handler http.Handler) http.HandlerFunc {
	digestServer := auth.NewDigest("realm", func(user, realm string) string {
		return auth.CalculateHA1(user, realm, *flagPassword)
	})
	return func(w http.ResponseWriter, r *http.Request) {
		if digestServer.Auth(w, r) {
			handler.ServeHTTP(w, r)
		}
	}
}

func main() {
	flag.Parse()
	fileServer := http.FileServer(http.Dir(*flagRoot))
	if *flagDigest {
		http.HandleFunc("/", makeBasicAuthHandleFunc(fileServer))
	} else {
		http.Handle("/", fileServer)
	}
	if err := http.ListenAndServe(*flagAddress, nil); err != nil {
		fmt.Println(err)
	}
}
