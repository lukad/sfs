/*
sfs - A simple http file server with digest access authentication support
Copyright (C) 2013,2014  Luka Dornhecker

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
	"github.com/wsxiaoys/terminal/color"
	"log"
	"net/http"
	"net/url"
)

var (
	flagAddress  = flag.StringP("listen", "l", ":8080", "Listen address")
	flagDigest   = flag.BoolP("digest", "d", false, "Use digest access authentication")
	flagUser     = flag.StringP("user", "u", "", "Username for authentication")
	flagPassword = flag.StringP("password", "p", "", "Password for authentication")
	flagRoot     = flag.StringP("root", "r", "./", "Root directory for the file server")
	flagLog      = flag.BoolP("log", "", true, "Log to stdout")
	flagNoColor  = flag.BoolP("no-color", "n", false, "Don't log with colors")
)

func Log(handler http.Handler) http.Handler {
	var logf func(string, *url.URL, string)
	if *flagNoColor {
		logf = func(method string, u *url.URL, address string) {
			log.Printf("%s %s %s", method, u, address)
		}
	} else {
		logf = func(method string, u *url.URL, address string) {
			logFormat := "@{y}%s @{m}%s @{|}%s"
			log.Print(color.Sprintf(logFormat, method, u, address))
		}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logf(r.Method, r.URL, r.RemoteAddr)
		handler.ServeHTTP(w, r)
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
