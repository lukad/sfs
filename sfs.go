/*
sfs - A simple http file server with digest access authentication support
Copyright (C) 2013  Luka Dornhecker

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	auth "bitbucket.org/taruti/http_auth"
	"fmt"
	flag "github.com/ogier/pflag"
	"log"
	"net/http"
)

var (
	flagAddress  = flag.StringP("listen", "l", ":8080", "Listen address")
	flagDigest   = flag.BoolP("digest", "d", false, "Use digest access authentication")
	flagUser     = flag.StringP("user", "u", "", "Username for authentication")
	flagPassword = flag.StringP("password", "p", "", "Password for authentication")
	flagRoot     = flag.StringP("root", "r", "./", "Root directory for the file server")
	flagLog      = flag.BoolP("log", "", false, "Log to stdout")
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		handler.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL, r.RemoteAddr)
	})
}

func makeBasicAuthHandleFunc(handler http.Handler) http.Handler {
	digestServer := auth.NewDigest("realm", func(user, realm string) string {
		return auth.CalculateHA1(user, realm, *flagPassword)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if digestServer.Auth(w, r) {
			handler.ServeHTTP(w, r)
		}
	})
}

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: sfs [options]\nWhere options are:")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	fileServer := http.FileServer(http.Dir(*flagRoot))
	var handler http.Handler
	if *flagDigest {
		handler = makeBasicAuthHandleFunc(fileServer)
	} else {
		handler = fileServer
	}
	if *flagLog {
		http.Handle("/", Log(handler))
	} else {
		http.Handle("/", handler)
	}
	if err := http.ListenAndServe(*flagAddress, nil); err != nil {
		fmt.Println(err)
	}
}
