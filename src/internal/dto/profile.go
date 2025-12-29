package dto

import (
	"errors"

	"github.com/arieffadhlan/go-fitbyte/internal/models"
)

type ProfileUpdateRequest struct {
	Preference *string `json:"preference" validate:"omitempty,oneof=CARDIO WEIGHT"`
	WeightUnit *string `json:"weightUnit" validate:"omitempty,oneof=KG LBS"`
	HeightUnit *string `json:"heightUnit" validate:"omitempty,oneof=CM INCH"`
	Weight     *int    `json:"weight" validate:"omitempty,min=10,max=1000"`
	Height     *int    `json:"height" validate:"omitempty,min=3,max=250"`
	Name       *string `json:"name" validate:"omitempty,min=2,max=60"`
	ImageURI   *string `json:"imageUri"`
}

type ProfileResponse struct {
	Preference *string `json:"preference"`
	WeightUnit *string `json:"weightUnit"`
	HeightUnit *string `json:"heightUnit"`
	Weight     *int    `json:"weight"`
	Height     *int    `json:"height"`
	Email      string  `json:"email"`
	Name       *string `json:"name"`
	ImageURI   *string `json:"imageUri"`
}

func (p *ProfileUpdateRequest) Validate() error {
	if p.Preference != nil {
		if *p.Preference != "CARDIO" && *p.Preference != "WEIGHT" {
			return errors.New("preference must be CARDIO or WEIGHT")
		}
	}

	if p.WeightUnit != nil {
		if *p.WeightUnit != "KG" && *p.WeightUnit != "LBS" {
			return errors.New("weightUnit must be KG or LBS")
		}
	}

	if p.HeightUnit != nil {
		if *p.HeightUnit != "CM" && *p.HeightUnit != "INCH" {
			return errors.New("heightUnit must be CM or INCH")
		}
	}

	if p.Weight != nil {
		if *p.Weight < 10 || *p.Weight > 1000 {
			return errors.New("weight must be between 10 and 1000")
		}
	}

	if p.Height != nil {
		if *p.Height < 3 || *p.Height > 250 {
			return errors.New("height must be between 3 and 250")
		}
	}

	if p.Name != nil {
		nameLen := len(*p.Name)
		if nameLen < 2 || nameLen > 60 {
			return errors.New("name must be between 2 and 60 characters")
		}
	}

	return nil
}

func (p *ProfileResponse) FromUser(user *models.User) {
	p.Email = user.Email
	p.Name = user.Name
	p.ImageURI = user.ImageURI
	p.Weight = user.Weight
	p.Height = user.Height

	if user.Preference != nil {
		pref := string(*user.Preference)
		p.Preference = &pref
	}

	if user.WeightUnit != nil {
		unit := string(*user.WeightUnit)
		p.WeightUnit = &unit
	}

	if user.HeightUnit != nil {
		unit := string(*user.HeightUnit)
		p.HeightUnit = &unit
	}
}
