package persistence

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// Operator is a single row of the mcc-mnc.csv dataset.
type Operator struct {
	Mcc         string
	Mnc         string
	Iso         string
	Country     string
	CountryCode string
	Network     string
}

// LoadOperators reads the mcc-mnc.csv dataset and returns the list of
// operators, left-padding MCC/MNC codes with zeros as required by the
// 3GPP EPDG FQDN format (mnc<MNC>.mcc<MCC>.pub.3gppnetwork.org).
func LoadOperators(path string) ([]Operator, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("LoadOperators(): failed to open %q: %w", path, err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("LoadOperators(): failed to parse csv: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("LoadOperators(): the CSV file is empty")
	}

	// Skip header row if present
	start := 0
	if len(records[0]) >= 6 && records[0][0] == "MCC" {
		start = 1
	}

	var operators []Operator
	for _, row := range records[start:] {
		if len(row) < 6 {
			log.Println("Row discarded due to an insufficient number of columns: ", len(row))
			continue
		}

		mcc, ok := padCode(row[0], "MCC")
		if !ok {
			continue
		}
		mnc, ok := padCode(row[1], "MNC")
		if !ok {
			continue
		}

		operators = append(operators, Operator{
			Mcc:         mcc,
			Mnc:         mnc,
			Iso:         row[2],
			Country:     row[3],
			CountryCode: row[4],
			Network:     row[5],
		})
	}

	return operators, nil
}

// padCode left-pads MCC/MNC codes with zeros to a length of 3 digits.
func padCode(code string, label string) (string, bool) {
	switch len(code) {
	case 0:
		log.Printf("Row discarded due to an invalid %s.\n", label)
		return "", false
	case 1:
		return "00" + code, true
	case 2:
		return "0" + code, true
	default:
		return code, true
	}
}
