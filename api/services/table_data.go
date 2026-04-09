package services

import "api/models"

func GetDataTableYearVersions(tableType string) ([]models.TableYearVersion, error) {
	var yearVersions []models.TableYearVersion

	// Get the data table versions for the specified table type
	var tableName string
	switch tableType {
	case "finance_variables":
		// get the disticnt years and versions from the database
		tableName = "finance_variables"
	case "paa_finance":
		tableName = "paa_finances"
	case "ra_factors":
		tableName = "risk_adjustment_factors"
	case "transition_adjustments":
		tableName = "balance_sheet_records"
	case "exp_age_bands":
		tableName = "exp_age_bands"
	case "exp_current_mortalities":
		tableName = "exp_current_mortalities"
	case "exp_current_lapses":
		tableName = "exp_current_lapses"
	case "modified_gmm_parameters":
		tableName = "modified_gmm_parameters"
	}

	var rows []struct {
		Year    int
		Version string
	}
	err := DB.Table(tableName).
		Select("year, version"). // Selects the 'year' and 'version' columns
		Distinct().              // Ensures each (year, version) pair is unique
		Order("year asc").       // Orders by year primarily
		Order("version asc").    // Then by version (ensures versions are sorted within each year)
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	// Initialize with the first year's data
	currentYear := rows[0].Year
	var currentVersions []string

	for _, row := range rows {
		if row.Year != currentYear {
			// New year encountered, save the previous year's data
			if len(currentVersions) > 0 { // Ensure we don't add an empty version set for a prior year if it somehow occurred
				yearVersions = append(yearVersions, models.TableYearVersion{
					Year:    currentYear,
					Version: currentVersions,
				})
			}
			// Reset for the new current year
			currentYear = row.Year
			currentVersions = []string{} // Initialize a new slice for versions
		}
		// Add the version to the current year's list.
		// Since the query has `Distinct()`, we don't need to check for duplicate versions here for the same year.
		currentVersions = append(currentVersions, row.Version)
	}

	// Append the last processed year's data
	if len(currentVersions) > 0 {
		yearVersions = append(yearVersions, models.TableYearVersion{
			Year:    currentYear,
			Version: currentVersions,
		})
	}

	return yearVersions, nil

}

func GetDataTableYears(tableType string) ([]int, error) {
	var years []int

	// Get the data table versions for the specified table type
	var tableName string
	switch tableType {
	case "finance_variables":
		tableName = "finance_variables"
	case "paa_finance":
		tableName = "paa_finances"
	case "ra_factors":
		tableName = "risk_adjustment_factors"
	case "transition_adjustments":
		tableName = "balance_sheet_records"
	case "exp_age_bands":
		tableName = "exp_age_bands"
	case "exp_current_mortalities":
		tableName = "exp_current_mortalities"
	case "exp_current_lapses":
		tableName = "exp_current_lapses"
	case "modified_gmm_parameters":
		tableName = "modified_gmm_parameters"
	case "modiified_gmm_shocks":
		tableName = "modified_gmm_shocks"
	case "premium_earning_patterns":
		tableName = "premium_earning_patterns"
	case "reinsurance_parameters":
		tableName = "reinsurance_parameters"
	case "paa_lapses":
		tableName = "paa_lapses"
	}

	err := DB.Table(tableName).
		Select("DISTINCT year"). // Selects distinct years
		Order("year asc").       // Orders by year
		Scan(&years).Error
	if err != nil {
		return nil, err
	}

	return years, nil
}
