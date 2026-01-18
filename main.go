// Copyright 2021 Marko Bevc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

const (
	version = "v1.0.0"
)

// Version
var Version string = "v1.0.0"

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	log.Println("Starting helloapp application..." + fmt.Sprintf(" | version: %s", Version))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello everyone!\n")
		log.Println("Hello request")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		var status int = 200
		var cacheStatus, dbStatus string

		log.Println("Health request")

		client := redis.NewClient(&redis.Options{
			Addr:     getEnv("REDIS_HOST", "localhost") + ":" + getEnv("REDIS_PORT", "6379"),
			Password: "",
			DB:       0,
		})

		pong, err := client.Ping().Result()
		cacheStatus = fmt.Sprintf("Redis: %s\n", pong)

		if err != nil {
			status = 500
			cacheStatus = fmt.Sprintf("Redis: %s", err)
		}

		db, err := sql.Open("mysql", getEnv("DB_USER", "root")+":"+getEnv("DB_PASS", "password1")+"@tcp("+getEnv("DB_HOST", "localhost")+":"+getEnv("DB_PORT", "3306")+")/"+getEnv("DB_NAME", "mysql"))
		dbStatus = fmt.Sprintf("MySQL: OK\n")
		if err := db.Ping(); err != nil {
			status = 503
			dbStatus = fmt.Sprintf("MySQL: %s", err.Error())
		}
		defer db.Close()

		output := fmt.Sprintf("Status: %s\n%s\n%s", strconv.Itoa(status), strings.TrimSuffix(cacheStatus, "\n"), strings.TrimSuffix(dbStatus, "\n"))
		log.Print(output)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		mapOutput := map[string]string{"status": strconv.Itoa(status), "Cache": strings.TrimSuffix(cacheStatus, "\n"), "DB": strings.TrimSuffix(dbStatus, "\n")}
		mapOJ, _ := json.Marshal(mapOutput)
		fmt.Fprintf(w, "%s\n", (string(mapOJ)))

		//fmt.Fprintf(w, "%s\n", output)
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", Version)
		log.Println("Version request")
	})

	s := http.Server{Addr: ":8080"}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}
