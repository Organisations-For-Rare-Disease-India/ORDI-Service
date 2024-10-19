package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/umahmood/haversine"
)

type Location struct {
	Latitude  float64
	Longitude float64
	District  string
}

type PincodeServiceConfig struct {
	filepath string
}
type PincodeService struct {
	pincodeData map[string]Location
}

func NewDefaultPincodeService() (*PincodeService, error) {
	return NewPincodeService(
		&PincodeServiceConfig{
			filepath: "data/pincode.csv",
		},
	)
}

func NewPincodeService(config *PincodeServiceConfig) (*PincodeService, error) {
	ps := &PincodeService{
		pincodeData: make(map[string]Location),
	}

	err := ps.loadPincodeData(config.filepath)
	if err != nil {
		return nil, fmt.Errorf("error loading pincode data: %w", err)
	}

	return ps, nil
}

// LoadPincodedata loads the pincode information to memory
func (ps *PincodeService) loadPincodeData(locationPath string) error {
	absPath, err := filepath.Abs(locationPath)
	if err != nil {
		return fmt.Errorf("could not get absolute path: %v", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("could not open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read CSV file: %v", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found in CSV")
	}
	records = records[1:] // Skipping the header row

	for _, record := range records {
		pincode := record[4]
		latitude, _ := strconv.ParseFloat(record[9], 64)
		longitude, _ := strconv.ParseFloat(record[10], 64)
		district := record[7]

		location := Location{
			Latitude:  latitude,
			Longitude: longitude,
			District:  district,
		}

		// Store the location in the map
		ps.pincodeData[pincode] = location
	}

	return nil
}

// Method to compute haversine distance in kilometers between 2 pincodes
func (ps *PincodeService) ComputeDistancePincodes(A, B string) (float64, error) {
	location1, err := ps.getLocation(A)
	if err != nil {
		return 0, err
	}

	location2, err := ps.getLocation(B)
	if err != nil {
		return 0, err
	}

	start := haversine.Coord{Lat: location1.Latitude, Lon: location1.Longitude}
	end := haversine.Coord{Lat: location2.Latitude, Lon: location2.Longitude}

	_, distance := haversine.Distance(start, end)
	return distance / 1000, nil // Return distance in kilometers

}

// Get location from pincode
func (ps *PincodeService) getLocation(pincode string) (Location, error) {
	location, found := ps.pincodeData[pincode]
	if !found {
		return Location{}, fmt.Errorf("pincode %s not found", pincode)
	}
	return location, nil
}
