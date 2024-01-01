// Copyright 2024 The PipeCD Authors.
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

package encrypt

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/pipe-cd/pipecd/pkg/app/pipectl/client"
	"github.com/pipe-cd/pipecd/pkg/app/server/service/apiservice"
	"github.com/pipe-cd/pipecd/pkg/cli"
)

// 10MB
const maxDataSize int64 = 10485760

type command struct {
	clientOptions *client.Options

	pipedID        string
	inputFile      string
	base64Encoding bool

	stdout io.Writer
}

func NewCommand() *cobra.Command {
	c := &command{
		clientOptions: &client.Options{},
		stdout:        os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt the plaintext entered in either stdin or the --input-file flag.",
		Example: `  pipectl encrypt --piped-id=xxx --api-key=yyy --address=foo.xz <secret.txt
  cat secret.txt | pipectl encrypt --piped-id=xxxxt --api-key=yyy --address=foo.xz
  pipectl encrypt --input-file=secret.txt --piped-id=xxxxt --api-key=yyy --address=foo.xz`,
		RunE: cli.WithContext(c.run),
	}

	cmd.Flags().StringVar(&c.pipedID, "piped-id", c.pipedID, "The id of Piped to which the application using the ciphertext belongs.")
	cmd.Flags().StringVar(&c.inputFile, "input-file", c.inputFile, "The path to the file to be encrypted.")
	cmd.Flags().BoolVar(&c.base64Encoding, "use-base64-encoding", c.base64Encoding, "Whether the plaintext should be base64 encoded before encrypting or not. (default false)")
	cmd.MarkFlagRequired("piped-id")

	c.clientOptions.RegisterPersistentFlags(cmd)

	return cmd
}

func (c *command) run(ctx context.Context, input cli.Input) error {
	cli, err := c.clientOptions.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize client: %w", err)
	}
	defer cli.Close()

	// Prioritize the file passed via the "--input-file" flag.
	var source io.Reader
	if c.inputFile != "" {
		fd, err := os.Open(c.inputFile)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", c.inputFile, err)
		}
		defer fd.Close()
		source = fd
	} else {
		source = input.Stdin
	}

	// Prevent accidental loading of large data into memory.
	reader := io.LimitReader(source, maxDataSize+1)
	buf := bytes.Buffer{}
	n, err := buf.ReadFrom(reader)
	if err != nil {
		return fmt.Errorf("failed to read the data: %w", err)
	}
	if n > maxDataSize {
		return fmt.Errorf("input data exceeds set limit 10 MB")
	}

	req := &apiservice.EncryptRequest{
		PipedId:        c.pipedID,
		Plaintext:      buf.String(),
		Base64Encoding: c.base64Encoding,
	}

	resp, err := cli.Encrypt(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}

	fmt.Fprintln(c.stdout, resp.Ciphertext)
	return nil
}
