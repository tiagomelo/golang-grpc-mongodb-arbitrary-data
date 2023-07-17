// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/config"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/server"
	"github.com/tiagomelo/golang-grpc-mongodb-arbitrary-data/store"
)

// run is the main entry point for the gRPC server.
// It sets up the server, initializes the necessary dependencies, and starts the server to listen for incoming requests.
func run(log *log.Logger) error {
	log.Println("main: initializing gRPC server")
	defer log.Println("main: Completed")

	ctx := context.Background()

	// =========================================================================
	// Config reading
	const envFilePath = ".env"
	cfg, err := config.Read(envFilePath)
	if err != nil {
		return errors.Wrap(err, "reading config")
	}

	// =========================================================================
	// Database support
	db, err := store.Connect(ctx, cfg.MongodbHostName, cfg.MongodbDatabase, cfg.MongodbPort)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}

	// =========================================================================
	// Listener init
	port := fmt.Sprintf(":%d", cfg.GrpcServerṔort)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "tcp listening")
	}

	// =========================================================================
	// Server init
	srv := server.New(db)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main: gRPC server listening on %s", port)
		serverErrors <- srv.GrpcSrv.Serve(lis)
	}()

	// =========================================================================
	// Shutdown
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Println("main: received signal for shutdown: ", sig)
		srv.GrpcSrv.Stop()
	}

	return nil
}

func main() {
	log := log.New(os.Stdout, "GRPC SERVER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	if err := run(log); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
