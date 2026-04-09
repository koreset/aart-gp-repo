package models

type BelBuildupBaseVariable struct {
	ID             int     `gorm:"primary_key"`
	Name           string  `json:"name" csv:"name"`
	VariableChange float64 `json:"variable_change" csv:"variable_change"`
	BelBuildup     float64 `json:"bel_buildup" csv:"bel_buildup"`
}

type BelBuildupVariable struct {
	ID                       int     `gorm:"primary_key" json:"-"`
	BelBuildupVariableSetID  int     `json:"-"`
	BelBuildupBaseVariableID int     `json:"-"`
	Name                     string  `json:"name" csv:"name"`
	VariableChange           float64 `json:"variable_change" csv:"variable_change"`
	BelBuildup               float64 `json:"bel_buildup" csv:"bel_buildup"`
	RAChange                 float64 `json:"ra_change" csv:"ra_change"`
	RABuildup                float64 `json:"ra_buildup" csv:"ra_buildup"`
	CSMChange                float64 `json:"csm_change" csv:"csm_change"`
	CSMBuildup               float64 `json:"csm_buildup" csv:"csm_buildup"`
	LossComponentChange      float64 `json:"loss_component_change" csv:"loss_component_change"`
	LossComponentUnwind      float64 `json:"loss_component_unwind" csv:"loss_component_unwind"`
	LossComponentBuildup     float64 `json:"loss_component_buildup" csv:"loss_component_buildup"`
	Notes                    string  `json:"notes" csv:"notes"`
	ReBelChange              float64 `json:"re_bel_change" csv:"re_bel_change"`
	ReBelBuildup             float64 `json:"re_bel_buildup" csv:"re_bel_buildup"`
	ReRAChange               float64 `json:"re_ra_change" csv:"re_ra_change"`
	ReRABuildup              float64 `json:"re_ra_buildup" csv:"re_ra_buildup"`
	ReCSMBuildup             float64 `json:"re_csm_buildup" csv:"re_csm_buildup"`
}

type BelBuildupVariableSet struct {
	ID                  int                  `json:"id" gorm:"primary_key"`
	ConfigurationName   string               `json:"configuration_name"`
	BelBuildupVariables []BelBuildupVariable `json:"lic_variables"`
}
