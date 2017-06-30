package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"tls-server/api"
	"tls-server/gencerts"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	mod := flag.String("config", "dev", "-config=dev run in development mode.")
	flag.Parse()

	var config map[string]string
	file, _ := os.Open(fmt.Sprintf("config.%s.json", *mod))
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("CONFIG ERROR: ", fmt.Sprintf("config.%s.json", *mod))
		panic(err)
	}

	certsdir := strings.Replace(config["addr"], ":", "-", -1)
	if _, err := os.Stat(fmt.Sprintf("certs/%s/cert.pem", certsdir)); err != nil {
		gencerts.Gen(config["addr"], true)
	}

	srv := &http.Server{
		Addr:              config["addr"],
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      8 * time.Second,
		IdleTimeout:       120 * time.Second,
		TLSConfig: &tls.Config{
			// knownGoodCipherSuites
			// Causes servers to use Go's default ciphersuite preferences,
			// which are tuned to avoid attacks. Does nothing on clients.
			PreferServerCipherSuites: true,
			// Only use curves which have assembly implementations
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519, // Go 1.8 only
			},
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

				// Best disabled, as they don't provide Forward Secrecy,
				// but might be necessary for some clients
				// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
		Handler: api.InitAPI(config),
	}

	osCh := make(chan os.Signal)
	signal.Notify(osCh, os.Interrupt, os.Kill)
	go func() {
		<-osCh

		srv.Close()
	}()

	fmt.Printf("STARTING...Listen https://%s\n", config["addr"])
	//if err := srv.ListenAndServeTLS(fmt.Sprintf("certs/%s/cert.pem", certsdir), fmt.Sprintf("certs/%s/key.pem", certsdir)); err != nil && err.Error() != "http: Server closed" {
	if err := srv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
		api.ShutdownAPI(err)
		return
	}

	api.ShutdownAPI(nil)
}

//
//go-torch -t=5 -u=http://127.0.0.1:8888 -p > profile.svg
