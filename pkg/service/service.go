package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/nikhil-github/api-cab-data/pkg/output"
)

const keyNotFound = "Key not found in cache"

// TripService embeds dependencies for counting trips.
type TripService struct {
	cacheGetter CacheGetter
	cacheSetter CacheSetter
	dbGetter    Getter
	logger      *zap.Logger
}

// Getter provides method to get trip count from DB.
type Getter interface {
	TripsByMedallionsOnPickUpDate(ctx context.Context, medallion string, pickUpDate time.Time) (output.Result, error)
	TripsByMedallion(ctx context.Context, medallions []string) ([]output.Result, error)
}

// CacheGetter provides method to get trip count from Cache.
type CacheGetter interface {
	Get(ctx context.Context, key string) (int, error)
}

// CacheSetter provides method to cache trip counts.
type CacheSetter interface {
	Set(ctx context.Context, key string, val int)
}

// New creates a new Tripservice.
func New(g Getter, cg CacheGetter, cs CacheSetter, l *zap.Logger) *TripService {
	return &TripService{dbGetter: g, cacheGetter: cg, cacheSetter: cs, logger: l}
}

// TripsByMedallionsOnPickUpDate get the number of trips for each medallion by pickup date.
// Check cache entries first before finding in DB.
// Query DB for cache misses.
func (s *TripService) TripsByMedallionsOnPickUpDate(ctx context.Context, medallion string, pickUpDate time.Time, byPassCache bool) (output.Result, error) {
	var err error
	var result output.Result
	if byPassCache {
		result, err = s.getFromDBByPickUpDate(ctx, medallion, pickUpDate)
		if err != nil {
			s.logger.Error("Error finding trips", zap.String("medallion", medallion), zap.Time("pickupdate", pickUpDate))
			return output.Result{}, err
		}
	} else {
		var tripCount int
		tripCount, err = s.cacheGetter.Get(ctx, key(medallion, pickUpDate))
		if err != nil && err.Error() == keyNotFound {
			result, err = s.getFromDBByPickUpDate(ctx, medallion, pickUpDate)
			if err != nil {
				s.logger.Error("Error finding trips", zap.String("medallion", medallion), zap.Time("pickupdate", pickUpDate))
				return output.Result{}, err
			}
		} else {
			result = output.Result{Medallion: medallion, Trips: tripCount}
		}
	}
	return result, nil
}

func (s *TripService) getFromDBByPickUpDate(ctx context.Context, medallion string, pickUpDate time.Time) (output.Result, error) {
	result, err := s.dbGetter.TripsByMedallionsOnPickUpDate(ctx, medallion, pickUpDate)
	if err != nil {
		return output.Result{}, err
	}
	go s.cacheSetter.Set(ctx, key(medallion, pickUpDate), result.Trips)
	return result, err
}

// TripsByMedallion get the number of trips for medallions.
// Check cache entries first before finding in DB.
// By pass cache with byPassCache flag equals true.
// Cache key is the medallion.
// Query DB for cache misses.
func (s *TripService) TripsByMedallion(ctx context.Context, medallions []string, byPassCache bool) ([]output.Result, error) {
	var results []output.Result
	var dbResults []output.Result
	var dbMedallions []string
	var err error
	if byPassCache {
		results, err = s.getFromDBByMedallion(ctx, medallions)
		if err != nil {
			s.logger.Error("Error finding trips for medallions", zap.Strings("medallions", medallions))
			return []output.Result{}, err
		}
	} else {
		var tripCount int
		for _, med := range medallions {
			tripCount, err = s.cacheGetter.Get(ctx, med)
			if err != nil && err.Error() == keyNotFound {
				dbMedallions = append(dbMedallions, med)
			} else {
				results = append(results, output.Result{Medallion: med, Trips: tripCount})
			}
		}
		if len(dbMedallions) > 0 {
			dbResults, err = s.getFromDBByMedallion(ctx, dbMedallions)
			if err != nil {
				s.logger.Error("Error finding trips for medallions", zap.Strings("medallions", medallions))
				return []output.Result{}, err
			}
		}
		results = append(results, dbResults...)
	}

	return results, nil

}

func (s *TripService) getFromDBByMedallion(ctx context.Context, medallions []string) ([]output.Result, error) {
	results, err := s.dbGetter.TripsByMedallion(ctx, medallions)
	if err != nil {
		return nil, err
	}
	go s.cacheMedallions(ctx, results)
	return results, nil
}

func (s *TripService) cacheMedallions(ctx context.Context, res []output.Result) {
	for _, r := range res {
		s.cacheSetter.Set(ctx, r.Medallion, r.Trips)
	}
}

// key is built by concatenate medallion + pickUpDate.
func key(medallion string, pickUpDate time.Time) string {
	return fmt.Sprintf("%s%d%d%d", medallion, pickUpDate.Year(), pickUpDate.Month(), pickUpDate.Day())
}
