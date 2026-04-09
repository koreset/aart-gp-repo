package main

import (
	"api/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// Set seed for random number generation
	rand.Seed(time.Now().UnixNano())

	// Generate 100 scheme names
	schemeNames := generateSchemeNames(100)

	// Generate 10,000 data points
	data := generateHistoricalCredibilityData(2000, schemeNames)

	// Write data to JSON file
	writeToJSON(data, "../historical_credibility_data.json")

	// Write data to CSV file
	writeToCSV(data, "../historical_credibility_data.csv")

	fmt.Printf("Successfully generated %d data points with %d scheme names\n", len(data), len(schemeNames))
}

func generateSchemeNames(count int) []string {
	prefixes := []string{"Corporate", "Enterprise", "Group", "Business", "Company", "Organization", "Institution", "Association", "Federation", "Union"}
	types := []string{"Insurance", "Health", "Life", "Pension", "Benefit", "Protection", "Coverage", "Assurance", "Security", "Welfare"}
	suffixes := []string{"Plan", "Scheme", "Program", "Policy", "Package", "Solution", "Arrangement", "System", "Framework", "Structure"}

	names := make([]string, count)
	for i := 0; i < count; i++ {
		prefix := prefixes[rand.Intn(len(prefixes))]
		typeStr := types[rand.Intn(len(types))]
		suffix := suffixes[rand.Intn(len(suffixes))]
		names[i] = fmt.Sprintf("%s %s %s %d", prefix, typeStr, suffix, i+1)
	}
	return names
}

func generateHistoricalCredibilityData(count int, schemeNames []string) []models.HistoricalCredibilityData {
	data := make([]models.HistoricalCredibilityData, count)

	basisOptions := []string{"Standard", "Enhanced", "Premium", "Basic", "Custom"}
	quoteTypeOptions := []string{"New Business", "Renewal", "Amendment", "Adjustment", "Review"}

	currentYear := time.Now().Year()

	for i := 0; i < count; i++ {
		// Generate random scheme index
		schemeIndex := rand.Intn(len(schemeNames))

		// Generate random data
		data[i] = models.HistoricalCredibilityData{
			ID:                       i + 1,
			QuoteID:                  rand.Intn(10000) + 1,
			Basis:                    basisOptions[rand.Intn(len(basisOptions))],
			CreationDate:             randomDate(2018, currentYear),
			QuoteType:                quoteTypeOptions[rand.Intn(len(quoteTypeOptions))],
			SchemeName:               schemeNames[schemeIndex],
			SchemeID:                 schemeIndex + 1,
			Year:                     rand.Intn(5) + (currentYear - 5),
			DurationInForce:          float64(rand.Intn(10) + 1),
			MemberCount:              rand.Intn(1000) + 10,
			ClaimCount:               rand.Intn(100),
			ExperiencePeriod:         float64(rand.Intn(10) + 1),
			CalculatedCredibility:    rand.Float64(),
			ManuallyAddedCredibility: rand.Float64() * 0.5,
		}
	}

	return data
}

func randomDate(minYear, maxYear int) time.Time {
	min := time.Date(minYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(maxYear, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func writeToJSON(data []models.HistoricalCredibilityData, filename string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Printf("Data successfully written to %s\n", filename)
}

func writeToCSV(data []models.HistoricalCredibilityData, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	header := []string{
		"id", 
		"quote_id", 
		"basis", 
		"creation_date", 
		"quote_type", 
		"scheme_name", 
		"scheme_id", 
		"year", 
		"duration_in_force", 
		"member_count", 
		"claim_count", 
		"experience_period", 
		"calculated_credibility", 
		"manually_added_credibility",
	}
	if err := writer.Write(header); err != nil {
		fmt.Printf("Error writing CSV header: %v\n", err)
		return
	}

	// Write data rows
	for _, record := range data {
		row := []string{
			strconv.Itoa(record.ID),
			strconv.Itoa(record.QuoteID),
			record.Basis,
			record.CreationDate.Format(time.RFC3339),
			record.QuoteType,
			record.SchemeName,
			strconv.Itoa(record.SchemeID),
			strconv.Itoa(record.Year),
			strconv.FormatFloat(record.DurationInForce, 'f', 2, 64),
			strconv.Itoa(record.MemberCount),
			strconv.Itoa(record.ClaimCount),
			strconv.FormatFloat(record.ExperiencePeriod, 'f', 2, 64),
			strconv.FormatFloat(record.CalculatedCredibility, 'f', 6, 64),
			strconv.FormatFloat(record.ManuallyAddedCredibility, 'f', 6, 64),
		}
		if err := writer.Write(row); err != nil {
			fmt.Printf("Error writing CSV row: %v\n", err)
			return
		}
	}

	fmt.Printf("Data successfully written to %s\n", filename)
}
