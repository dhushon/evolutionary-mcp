package models

import (
	"time"
)

type Tenant struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Domain     string    `json:"domain"`
	LogoSVG    string    `json:"logo_svg,omitempty"`
	BrandTitle string    `json:"brand_title,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}