package models

type Superhero struct {
	ID            int    `json:"id"`
	SuperheroName string `json:"superhero_name"`
	FullName      string `json:"full_name"`
	GenderID      int    `json:"gender_id"`
	EyeColourID   int    `json:"eye_colour_id"`
	HairColourID  int    `json:"hair_colour_id"`
	SkinColourID  int    `json:"skin_colour_id"`
	RaceID        int    `json:"race_id"`
	PublisherID   int    `json:"publisher_id"`
	AlignmentID   int    `json:"alignment_id"`
	HeightCm      int    `json:"height_cm"`
	WeightKg      int    `json:"weight_kg"`
}