package main

import (
	"io"
	"log"
	"login/conn"
	routes "login/routes"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// routes.StartGin()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:        false,
		AllowOrigins:           []string{"http://localhost:3000", "http://localhost:3001", "https://localhost:3000"},
		AllowMethods:           []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders:           []string{"Content-Length", "Content-Type"},
		AllowCredentials:       true,
		ExposeHeaders:          []string{"Content-Type"},
		AllowWildcard:          true,
		AllowBrowserExtensions: false,
	}))
	err := conn.InitializeDB()
	if err != nil {
		log.Fatal("session err: ", err)
	}
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	// caCert, _ := ioutil.ReadFile("ca.crt")
	// caCertPool, _ := x509.SystemCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)
	// tls := &tls.Config{
	// 	ClientCAs:             caCertPool,
	// 	ClientAuth:            tls.RequireAndVerifyClientCert,
	// 	GetCertificate:        utils.CertReqFunc("cert.pem", "key.pem"),
	// 	VerifyPeerCertificate: utils.CertificateChains,
	// }
	// server := http.Server{
	// 	Addr:      ":8080",
	// 	Handler:   router,
	// 	TLSConfig: tls,
	// }
	// err = server.ListenAndServeTLS("cert.pem", "key.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// router.RunTLS(":8080", "cert.pem", "key.pem")
	router.Run(":8080")
}
