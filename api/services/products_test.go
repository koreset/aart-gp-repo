package services

import (
	"api/models"
	"mime/multipart"
	"testing"
)

func TestSaveModelPoints2(t *testing.T) {
	type args struct {
		v       *multipart.FileHeader
		product models.Product
		year    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveModelPoints(tt.args.v, tt.args.product, tt.args.year, ""); (err != nil) != tt.wantErr {
				t.Errorf("SaveModelPoints() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
