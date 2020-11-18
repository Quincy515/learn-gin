package main

import (
	"crypto/tls"
	"crypto/x509"
	"gin-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net"
)

func main() {
	//creds, err := credentials.NewServerTLSFromFile("keys/client.pem", "keys/client.key")
	//if err != nil {
	//	log.Fatal(err)
	//}

	cert, _ := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, // 服务端证书
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	rpcServer := grpc.NewServer(grpc.Creds(creds))

	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))

	listen, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(listen)

	//mux := http.NewServeMux()
	//mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Println(request)
	//	rpcServer.ServeHTTP(writer, request)
	//})
	//httpServer := &http.Server{
	//	Addr:    ":8081",
	//	Handler: mux,
	//}
	//err := httpServer.ListenAndServeTLS("keys/client.pem", "keys/client.key")
	//if err != nil {
	//	log.Fatal(err)
	//}
}
