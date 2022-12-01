// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	"github.com/golang/glog"

	spannerClient "github.com/google/hashr/tool/clients/gcp"
	postgresClient "github.com/google/hashr/tool/clients/postgres"
	"github.com/google/hashr/tool/common"
	_ "github.com/lib/pq"
)

var (
	spannerDBPath = flag.String("spanner_db_path", "", "Path to spanner DB.")
	storage       = flag.String("storage", "", fmt.Sprintf("Type of storage backend to query: %s,%s", postgresClient.Name, spannerClient.Name))
	export        = flag.Bool("raw", false, "Input hashes are raw bytes and not hex-encoded")

	// Postgres DB flags
	postgresHost     = flag.String("postgres_host", "localhost", "PostgreSQL instance address.")
	postgresPort     = flag.Int("postgres_port", 5432, "PostgresSQL instance port.")
	postgresUser     = flag.String("postgres_user", "hashr", "PostgresSQL user.")
	postgresPassword = flag.String("postgres_password", "hashr", "PostgresSQL password.")
	postgresDBName   = flag.String("postgres_db", "hashr", "PostgresSQL database.")
)

func main() {
	ctx := context.Background()
	flag.Parse()
	var client common.Client = nil

	switch *storage {
	case postgresClient.Name:
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			*postgresHost, *postgresPort, *postgresUser, *postgresPassword, *postgresDBName)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			glog.Exitf("Error initializing Postgres connection: %v", err)
		}
		defer db.Close()

		clientInstance, err := postgresClient.NewClient(db)
		if err != nil {
			glog.Exitf("Error initializing Postgres client: %v", err)
		}
		client = clientInstance

	case spannerClient.Name:
		clientInstance, err := spannerClient.NewClient(ctx, spannerDBPath)
		if err != nil {
			glog.Exitf("Error initializing Spanner client: %v", err)
		}
		client = clientInstance
	default:
		glog.Exitf("Unknown storage type: %s", *storage)
	}

	// TODO: read lines from stdin in batches and lookup in database
	hashes := make([]string, 0)
	hashes = append(hashes, "12")
	hashes = append(hashes, "34")
	hashes = append(hashes, "56")
	fmt.Println("{}", hashes)
	res, err := client.LookupSamples(hashes)
	if err != nil {
		glog.Exitf("Failed to lookup samples")
	}

	// TODO: print results in a good way
	fmt.Println("{}", res)
}
