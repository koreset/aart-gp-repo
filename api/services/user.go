package services

import (
	appLog "api/log"
	"api/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// StoreUserToken stores a user token in the database
func StoreUserToken(tokenString string, subject string) error {
	logger := appLog.WithFields(map[string]interface{}{
		"subject": subject,
		"action":  "StoreUserToken",
	})

	logger.Debug("Storing user token")

	var userToken models.UserToken
	userToken.TokenString = tokenString
	userToken.Subject = subject

	err := DB.Where("subject=?", subject).Assign(userToken).FirstOrCreate(&userToken).Error
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to store user token")
		return err
	}

	logger.Info("User token stored successfully")
	return nil
}

// StoreActivity stores an activity record in the database
func StoreActivity(activity models.Activity) error {
	logger := appLog.WithFields(map[string]interface{}{
		"user_email":    activity.UserEmail,
		"user_name":     activity.UserName,
		"activity_type": activity.Type,
		"object_type":   activity.ObjectType,
		"object_id":     activity.ObjectID,
		"action":        "StoreActivity",
	})

	logger.Debug("Storing activity")

	err := DB.Save(&activity).Error
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to store activity")
		return err
	}

	logger.Info("Activity stored successfully")
	return nil
}

// CreateActivity creates and stores an activity record
func CreateActivity(ctx context.Context, activityType string, description string, objectId int, objectType string) {
	logger := appLog.WithContext(ctx).WithFields(map[string]interface{}{
		"activity_type": activityType,
		"object_type":   objectType,
		"object_id":     objectId,
		"action":        "CreateActivity",
	})

	logger.Debug("Creating activity")

	var activity models.Activity
	userData, ok := ctx.Value("keys").(map[string]interface{})
	if !ok {
		logger.Error("Failed to get user data from context")
		return
	}

	userName, ok := userData["userName"].(string)
	if !ok {
		logger.Error("Failed to get user name from context")
		return
	}

	userEmail, ok := userData["userEmail"].(string)
	if !ok {
		logger.Error("Failed to get user email from context")
		return
	}

	activity.UserName = userName
	activity.UserEmail = userEmail
	activity.Date = time.Now()
	activity.Type = activityType
	activity.Description = description
	activity.ObjectID = objectId
	activity.ObjectType = objectType

	logger = logger.WithFields(map[string]interface{}{
		"user_name":  userName,
		"user_email": userEmail,
	})

	err := StoreActivity(activity)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to store activity")
		return
	}

	logger.Info("Activity created and stored successfully")
}

// FindOrgUsers finds users for an organization
func FindOrgUsers(company models.Organisation) (users []models.OrgUser) {
	logger := appLog.WithFields(map[string]interface{}{
		"organisation": company.Name,
		"action":       "FindOrgUsers",
	})

	logger.Info("Finding organization users")

	// Use caching to improve performance
	var orgUsers []models.OrgUser

	err := DB.Where("organisation = ?", company.Name).Find(&orgUsers).Error

	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to query organization users from database")
	} else {
		logger.WithField("user_count", len(orgUsers)).Debug("Found users in database")
	}

	if len(orgUsers) > 0 {
		logger.Info("Returning users from database")
		return orgUsers
	}

	// need to make a call the keygen api to get the org users
	logger.Info("No users found in database, fetching from license server")

	var licenseServerUrl string
	licenseServerUrl = "https://licenses.aart-enterprise.com"

	// instantiate an http client to call the license server
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	jsonBody, err := json.Marshal(&company)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to marshal organization to JSON")
		return nil
	}

	logger.Debug("Making request to license server")
	req, err := http.NewRequest("POST", licenseServerUrl+"/get-org-users", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to create HTTP request")
		return nil
	}

	req.Header.Set("Accept", "application/json")

	startTime := time.Now()
	resp, err := client.Do(req)
	requestDuration := time.Since(startTime)

	if err != nil {
		logger.WithFields(map[string]interface{}{
			"error":       err.Error(),
			"duration_ms": requestDuration.Milliseconds(),
		}).Error("Failed to execute HTTP request")
		return nil
	}

	defer resp.Body.Close()

	logger.WithFields(map[string]interface{}{
		"status_code": resp.StatusCode,
		"duration_ms": requestDuration.Milliseconds(),
	}).Debug("Received response from license server")

	if resp.StatusCode != http.StatusOK {
		logger.WithField("status", resp.Status).Error("License server returned non-OK status")
		return nil
	}

	var result []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to decode response from license server")
		return nil
	}

	logger.WithField("user_count", len(result)).Info("Successfully retrieved users from license server")

	users = make([]models.OrgUser, 0)
	// iterate over the result and create a list of OrgUser
	for _, user := range result {
		var orgUser models.OrgUser

		name, ok := user["name"].(string)
		if !ok {
			logger.WithField("user", user).Warn("User name not found or not a string, skipping user")
			continue
		}

		email, ok := user["user"].(string)
		if !ok {
			logger.WithField("user", user).Warn("User email not found or not a string, skipping user")
			continue
		}

		licenseId, ok := user["license_id"].(string)
		if !ok {
			logger.WithField("user", user).Warn("License ID not found or not a string, skipping user")
			continue
		}

		orgUser.Name = name
		orgUser.Email = email
		orgUser.LicenseId = licenseId
		orgUser.GPRole = "None"
		orgUser.Organisation = company.Name
		users = append(users, orgUser)
	}

	if len(users) > 0 {
		// we got here, now we need to save the users to the database
		logger.Info("Saving users to database")
		err = DB.CreateInBatches(&users, 100).Error
		if err != nil {
			logger.WithField("error", err.Error()).Error("Failed to save users to database")
		} else {
			logger.Info("Successfully saved users to database")
		}
	} else {
		logger.Warn("No valid users found to save to database")
	}

	return users
}
