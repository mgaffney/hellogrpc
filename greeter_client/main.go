/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/mgaffney/grpc/helloworld/helloworld"
	"google.golang.org/grpc"
)

const (
	serviceAddr = "localhost:50051"
)

func greet(name string) (string, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure())
	if err != nil {
		return "", fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		return "", fmt.Errorf("could not greet: %v", err)
	}
	return r.Message, nil
}

func main() {
	// use PORT environment variable, or default to 8080
	port := "8080"
	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	// register hello function to handle all requests
	server := http.NewServeMux()
	server.HandleFunc("/", hello)

	// start the web server on port and accept requests
	log.Printf("Server listening on port %s", port)
	err := http.ListenAndServe(":"+port, server)
	log.Fatal(err)
}

// hello responds to the request with a plain-text "Hello, world" message.
func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	host, _ := os.Hostname()
	name := r.FormValue("name")
	if name == "" {
		name = "world"
	}
	msg, err := greet(name)
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "Error getting greeting: %v\n", err)
		msg = "sorry"
	}
	fmt.Fprintf(w, "Greeting: %s\n", msg)
	fmt.Fprintf(w, "Version: 1.0.0\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)
}
