package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const (
	healthURL            = "http://localhost:3000/health"
	tripsURLbyPickUpdate = "http://localhost:3000/trips/v1/medallion/67EB082BFFE72095EAF18488BEA96050/pickupdate/2013-12-31?bypasscache=true"
	tripsURLbyMedallions = "http://localhost:3000/trips/v1/medallions/67EB082BFFE72095EAF18488BEA96050,D7D598CD99978BD012A87A76A7C891B7?bypasscache=true"
	clearCacheURL        = "http://localhost:3000/trips/v1/cache/contents"
)

// Results represents all output.
type Results struct {
	Res []Result
}

// Result represents output.
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
	callTripServiceWithPickUpDate(ctx, netClient)
	callTripServiceWithMedallions(ctx, netClient)
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

func callTripServiceWithPickUpDate(ctx context.Context, client *http.Client) {
	r, err := http.NewRequest("GET", tripsURLbyPickUpdate, nil)
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
	var result Result
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("failed to decode the response")
	}
	log.Printf("response=> %#v\n", result)
}

func callTripServiceWithMedallions(ctx context.Context, client *http.Client) {
	r, err := http.NewRequest("GET", tripsURLbyMedallions, nil)
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
	r, err := http.NewRequest("DELETE", clearCacheURL, nil)
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
