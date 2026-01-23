package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"

	r "github.com/czxrny/veh-sense-backend/batch-receiver/internal/domain/upload/repository"
	"github.com/czxrny/veh-sense-backend/batch-receiver/internal/model"
	"github.com/czxrny/veh-sense-backend/shared/models"
)

const (
	EventAggressiveAccel model.RideEventType = "acceleration"
	EventAggressiveBrake model.RideEventType = "braking"
	EventHighRPM         model.RideEventType = "high_rpm"
	EventHighEngineLoad  model.RideEventType = "high_engine_load"
	EventOverspeed       model.RideEventType = "overspeed"
)

type Thresholds struct {
	// Acceleration and breaking in m/s
	AggressiveAccel float64
	AggressiveBrake float64

	// RPM
	HighRPM int

	// Engine load %
	HighEngineLoad int

	// Overspeed (km/h)
	OverspeedKmh int
}

var thresholds = Thresholds{
	AggressiveAccel: 2.6,  // >= 2.6 m/s^2 (0-100kmh in about 10-11s)
	AggressiveBrake: -3.0, // <= -3.0 m/s^2
	HighRPM:         3000, // np. >= 3000
	HighEngineLoad:  80,   // np. >= 80
	OverspeedKmh:    140,  // np. >= 140kmh
}

type Service struct {
	reportRepo *r.ReportRepository
	dataRepo   *r.ReportDataRepository
	userRepo   *r.UserInfoRepository
}

func NewService(reportRepo *r.ReportRepository, dataRepo *r.ReportDataRepository, userRepo *r.UserInfoRepository) *Service {
	return &Service{reportRepo: reportRepo, dataRepo: dataRepo, userRepo: userRepo}
}

func (s *Service) UploadRide(ctx context.Context, authInfo models.AuthInfo, req model.UploadRideRequest) (*models.Report, error) {
	rawFrames, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return nil, fmt.Errorf("invalid base64: %w", err)
	}

	frames, err := parseFrames(rawFrames)
	if err != nil {
		return nil, err
	}
	if len(frames) < 2 {
		return nil, errors.New("not enough frames to create report")
	}

	// sort in case the report is not already sorted
	sort.Slice(frames, func(i, j int) bool { return frames[i].Timestamp < frames[j].Timestamp })

	report, events, err := buildReportAndEvents(frames)
	if err != nil {
		return nil, err
	}

	report.UserID = authInfo.UserID
	report.OrganizationID = authInfo.OrganizationID
	report.VehicleID = int(req.VehicleID)
	report.ID = 0

	err = s.reportRepo.Add(ctx, report)
	if err != nil {
		return nil, err
	}

	eventData, err := toGzip(events)
	if err != nil {
		return nil, err
	}

	dataRecord := &model.RawRideRecord{
		ReportID:  report.ID,
		Data:      rawFrames,
		EventData: eventData,
	}

	err = s.dataRepo.Add(ctx, dataRecord)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.userRepo.GetByID(ctx, authInfo.UserID)
	if err != nil {
		return nil, err
	}

	userInfo.NumberOfRides++
	userInfo.TotalKilometers += int(report.KilometersTravelled)

	err = s.userRepo.Update(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func parseFrames(data []byte) ([]model.ObdFrame, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid gzip: %w", err)
	}
	defer gz.Close()

	var frames []model.ObdFrame
	if err := json.NewDecoder(gz).Decode(&frames); err != nil {
		return nil, fmt.Errorf("invalid OBD JSON: %w", err)
	}

	return frames, nil
}

func buildReportAndEvents(frames []model.ObdFrame) (*models.Report, []model.RideEvent, error) {
	if len(frames) < 2 {
		return nil, nil, errors.New("not enough frames")
	}

	start := frames[0].Timestamp
	end := frames[len(frames)-1].Timestamp
	if end <= start {
		return nil, nil, errors.New("invalid timestamps (end <= start)")
	}

	var (
		sumSpeed float64
		maxSpeed float64
		km       float64

		accelCount, brakeCount       int
		accelAggCount, brakeAggCount int
		accelPeak, brakePeak         float64

		events = make([]model.RideEvent, 0, 64)
	)

	previousFrame := frames[0]
	numberOfFramesWhileMoving := len(frames)

	for i := 1; i < len(frames); i++ {
		currentFrame := frames[i]

		// ---- SPEED STATS ----
		currentSpeed := float64(currentFrame.VehicleSpeed)

		// If not moving - do not count to avg speed
		if currentSpeed > 0 {
			sumSpeed += currentSpeed
		} else {
			numberOfFramesWhileMoving--
		}

		if currentSpeed > maxSpeed {
			maxSpeed = currentSpeed
		}

		// Calculating the delta between previous and current frame (in seconds)
		dt := float64(currentFrame.Timestamp-previousFrame.Timestamp) / 1000.0

		// ---- DISTANCE ----
		v1 := kmhToMS(float64(previousFrame.VehicleSpeed))
		v2 := kmhToMS(currentSpeed)
		km += ((v1 + v2) / 2) * (dt / 1000.0)

		// ---- ACCELERATION / BRAKING ----
		acc := (v2 - v1) / dt

		// Threshold for small changes in speed
		if acc > 0.2 {
			accelCount++
			if acc > accelPeak {
				accelPeak = acc
			}
			if acc >= thresholds.AggressiveAccel {
				accelAggCount++
				events = append(events, model.RideEvent{
					Timestamp: currentFrame.Timestamp,
					Type:      EventAggressiveAccel,
					Value:     acc,
				})
			}
			// same for braking
		} else if acc < -0.2 {
			brakeCount++
			if acc < brakePeak {
				brakePeak = acc
			}
			if acc <= thresholds.AggressiveBrake {
				brakeAggCount++
				events = append(events, model.RideEvent{
					Timestamp: currentFrame.Timestamp,
					Type:      EventAggressiveBrake,
					Value:     acc,
				})
			}
		}

		// ---- RPM ----
		if currentFrame.Rpm >= thresholds.HighRPM {
			events = append(events, model.RideEvent{
				Timestamp: currentFrame.Timestamp,
				Type:      EventHighRPM,
				Value:     float64(currentFrame.Rpm),
			})
		}

		// ---- ENGINE LOAD ----
		if currentFrame.EngineLoad >= thresholds.HighEngineLoad {
			events = append(events, model.RideEvent{
				Timestamp: currentFrame.Timestamp,
				Type:      EventHighEngineLoad,
				Value:     float64(currentFrame.EngineLoad),
			})
		}

		// ---- OVERSPEED ----
		if thresholds.OverspeedKmh > 0 && currentFrame.VehicleSpeed >= thresholds.OverspeedKmh {
			events = append(events, model.RideEvent{
				Timestamp: currentFrame.Timestamp,
				Type:      EventOverspeed,
				Value:     float64(currentFrame.VehicleSpeed),
			})
		}

		previousFrame = currentFrame
	}

	avgSpeed := sumSpeed / float64(numberOfFramesWhileMoving)

	accStyle := classifyStyle(
		accelCount,
		accelAggCount,
		accelPeak,
		thresholds.AggressiveAccel,
	)

	brkStyle := classifyStyle(
		brakeCount,
		brakeAggCount,
		math.Abs(brakePeak),
		math.Abs(thresholds.AggressiveBrake),
	)

	report := &models.Report{
		StartTime:           frames[0].Timestamp,
		StopTime:            frames[0].Timestamp,
		AccelerationStyle:   accStyle,
		BrakingStyle:        brkStyle,
		AverageSpeed:        avgSpeed,
		MaxSpeed:            maxSpeed,
		KilometersTravelled: km,
	}

	return report, events, nil
}

func kmhToMS(kmh float64) float64 { return kmh / 3.6 }

func classifyStyle(
	totalMoves int,
	aggressiveMoves int,
	peak float64,
	aggressiveThreshold float64,
) string {
	if totalMoves < 5 {
		return "calm"
	}

	ratio := float64(aggressiveMoves) / float64(totalMoves)

	switch {
	case ratio >= 0.20 || peak >= aggressiveThreshold*1.3:
		return "aggressive"

	case ratio >= 0.07 || peak >= aggressiveThreshold:
		return "normal"

	default:
		return "calm"
	}
}

func toGzip[T any](data []T) ([]byte, error) {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if _, err := gz.Write(dataJSON); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
