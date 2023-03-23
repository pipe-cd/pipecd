// Copyright 2023 The PipeCD Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const dir = "/public"

func main() {
	var (
		doneCh = make(chan error)
		mux    = http.NewServeMux()
		server = &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}
		fs = http.FileServer(http.Dir(dir))
	)
	mux.Handle("/", fs)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		fmt.Println("stopping http server")
		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("failed to shutdown http server: %v\n", err)
			return
		}
		fmt.Println("http server is stopped")
	}()

	fmt.Println("start running http server on 8080")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("failed to listen and serve http server: %v\n", err)
			doneCh <- err
		}
		doneCh <- nil
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	select {
	case <-ch:
		return
	case <-doneCh:
		return
	}
}
