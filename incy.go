/*
    incy (wincy) web server written in go

    run:
        incy serve [options]
    options:
        --https
        --cert <cert-path>
        --dir <dir-path>
        --domain <domain>
        --key <key-path>
        --port <port> 
    eg:
        //serve pwd over http at localhost on port 80
        incy serve 
        //serve pwd over https at example.com on port 443 using certbot generated certs
        incy serve --https --domain example.com --port 443 
*/
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

var https   = false
var dir     = "./"
var domain  = "localhost"
var port    = "80"
var version = "v0.1.0"

var cert = ""
var key  = ""

func main() {
    cmd := os.Args[1]
    if(cmd == "serve") {
        for p:=2; p<len(os.Args); p++ {
            if(os.Args[p] == "--https") {
                https = true;
            } else if(os.Args[p] == "--cert") {
                https = true;
                p++;
                cert = os.Args[p];
            } else if(os.Args[p] == "--domain") {
                p++;
                domain = os.Args[p];
            } else if(os.Args[p] == "--dir") {
                p++;
                dir = os.Args[p];
            } else if(os.Args[p] == "--key") {
                https = true;
                p++;
                key = os.Args[p];
            } else if(os.Args[p] == "--port") {
                p++;
                port = os.Args[p];
            } 
        }

        if((cert == "" || key == "") && cert != key) {
            fmt.Println("must specify both cert and key files")
            return
        }

        mux := http.NewServeMux()
        mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(dir))))

        if(https) {
            if(domain == "localhost") {
                fmt.Println("must specify domain for https")
                return
            }

            if(cert == "") {
                cert = "/etc/letsencrypt/live/" + domain + "/fullchain.pem"
            }

            if(key == "") {
                key = "/etc/letsencrypt/live/" + domain + "/privkey.pem"
            }

            if(port == "80") {
                port = "443";
            }

            fmt.Println("Serving " + domain + " over https on port " + port)
            for {
                err := http.ListenAndServeTLS(":" + port, cert, key, mux);
                if err != nil {
                    log.Println("ListenAndServeTLS: ", err)
                }
            }
        } else {
            fmt.Println("Serving " + domain + " over http on port " + port)
            for {
                err := http.ListenAndServe(":" + port, mux)
                if err != nil {
                    log.Println("ListenAndServe: ", err)
                }
            }
        } 
    } else if(cmd == "--version") {
        fmt.Println(version)
    } else {
        fmt.Println("do not recognize command '" + cmd + "'")
    }
}
