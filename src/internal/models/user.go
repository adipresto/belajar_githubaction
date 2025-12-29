package models

import "time"

type PreferenceType string
type WeightUnitType string
type HeightUnitType string

const (
	PreferenceCardio PreferenceType = "CARDIO"
	PreferenceWeight PreferenceType = "WEIGHT"

	WeightUnitKg WeightUnitType = "KG"
	WeightUnitLb WeightUnitType = "LBS"

	HeightUnitCm HeightUnitType = "CM"
	HeightUnitIn HeightUnitType = "INCH"
)

type User struct {
	ID         int             `db:"id"`
	Name       *string         `db:"name"`
	Email      string          `db:"email"`
	Password   string          `db:"password"`
	Preference *PreferenceType `db:"preference"`
	WeightUnit *WeightUnitType `db:"weight_unit"`
	HeightUnit *HeightUnitType `db:"height_unit"`
	Weight     *int            `db:"weight"`
	Height     *int            `db:"height"`
	ImageURI   *string         `db:"image_uri"`
	CreatedAt  time.Time       `db:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at"`
}
