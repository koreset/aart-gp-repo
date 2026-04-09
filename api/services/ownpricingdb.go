package services

import (
	"api/models"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"math"
	"strings"
)


func (db *OwnDb) BatchInsertPricingPoints(objArr []models.PricingPoint) (int64, error) {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return 0, errors.New("slice must not be empty")
	}

	numberOfFields := db.calculateNumberOfPricingPointFields((objArr)[0])
	return db.batchInsertPricingPointChunks(splitIntoPricingPointChunks(objArr, numberOfFields))
}

func splitIntoPricingPointChunks(objArr []models.PricingPoint, numberOfFields int) [][]models.PricingPoint {
	var chunks [][]models.PricingPoint

	chunkSize := int(math.Floor(float64(maxNumberOfBatchParameters / float32(numberOfFields))))
	numberOfObjects := len(objArr)

	if numberOfObjects < chunkSize {
		return [][]models.PricingPoint{objArr}
	}

	for i := 0; i < numberOfObjects; i += chunkSize {
		end := i + chunkSize

		if end > numberOfObjects {
			end = numberOfObjects
		}

		chunks = append(chunks, objArr[i:end])
	}

	return chunks
}

func (db *OwnDb) calculateNumberOfPricingPointFields(obj interface{}) int {
	return len(db.NewScope(obj).Fields())
}

func (db *OwnDb) batchInsertPricingPointChunks(chunks [][]models.PricingPoint) (int64, error) {
	var rowsAffected int64 = 0
	for _, chunk := range chunks {
		chunkRowsAffected, err := db.batchInsertPricingPoint(chunk)
		if err != nil {
			return 0, err
		}

		rowsAffected += chunkRowsAffected
	}

	return rowsAffected, nil
}

func (db *OwnDb) batchInsertPricingPoint(objArr []models.PricingPoint) (int64, error) {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return 0, errors.New("slice must not be empty")
	}

	mainObj := objArr[0]
	mainScope := db.NewScope(mainObj)
	mainFields := mainScope.Fields()
	quoted := make([]string, 0, len(mainFields))
	for i := range mainFields {
		// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
		// If field is ignore field, skip it.
		if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) {
			continue
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}

	placeholdersArr := make([]string, 0, len(objArr))

	for _, obj := range objArr {
		scope := db.NewScope(obj)
		fields := scope.Fields()

		placeholders := make([]string, 0, len(fields))
		for i := range fields {
			if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) {
				continue
			}
			var vars interface{}
			if (fields[i].Name == "CreatedAt" || fields[i].Name == "UpdatedAt") && fields[i].IsBlank {
				vars = gorm.NowFunc()
			} else {
				vars = fields[i].Field.Interface()
			}
			placeholders = append(placeholders, mainScope.AddToVars(vars))
		}

		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
	}

	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
		mainScope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholdersArr, ", "),
	))

	// Execute and Log
	if err := mainScope.Exec().DB().Error; err != nil {
		field, ok := mainScope.FieldByName("policy_number")
		if ok {
			fmt.Println(field.Field.Interface())
		}
		return 0, err
	}
	return mainScope.DB().RowsAffected, nil
}

