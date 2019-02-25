package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	healthURL = "http://localhost:3000/health"
	tripsURL  = "http://localhost:3000/trips/v1/medallion/67EB082BFFE72095EAF18488BEA96050?pickupdate=2013-12-31&bypasscache=true"
)

type Results struct {
	Res []Result
}

type Result struct {
	Medallion string `json:"medallion"`
	Trips     int    `json:"trips"`
}

func main() {
	ctx := context.Background()
	var netClient = &http.Client{
		Timeout: time.Second * 5,
	}
	callHealthCheck(ctx, netClient)
	callTripService(ctx, netClient)
	callClearCacheService(ctx, netClient)
}

func callHealthCheck(ctx context.Context, client *http.Client) {
	r, err := http.NewRequest("GET", healthURL, nil)
	if err != nil {
		log.Fatalf("health check request creation failed")
	}
	res, err := client.Do(r.WithContext(ctx))
	if err != nil {
		log.Fatalf("health check failed, please check the health of service")
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("health check failed healthRes.StatusCode %d", res.StatusCode)
	}
	log.Println("health check passed!!")
}

func callTripService(ctx context.Context, client *http.Client) {
	r, err := http.NewRequest("GET", tripsURL, nil)
	if err != nil {
		log.Fatalf("trip svc request creation failed")
	}
	res, err := client.Do(r.WithContext(ctx))
	if err != nil {
		log.Fatalf("trip service failed")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("failure with status %d", res.StatusCode)
	}
	var result []Result
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("failed to decode the response")
	}
	log.Printf("response=> %#v\n", result)
}

func callClearCacheService(ctx context.Context, client *http.Client) {
	r, err := http.NewRequest("GET", tripsURL, nil)
	if err != nil {
		log.Fatalf("clear cache request creation failed")
	}
	res, err := client.Do(r.WithContext(ctx))
	if err != nil {
		log.Fatalf("clear cache failed")
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("clear cache failed %d", res.StatusCode)
	}
	log.Printf("Successfully cleared cache entries")
}
