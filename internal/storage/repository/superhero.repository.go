package repository

import (
	"context"
	"go-ws/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SuperheroRepository struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func NewSuperheroRepository(db *pgxpool.Pool, ctx context.Context) *SuperheroRepository {
	return &SuperheroRepository{db: db, ctx: ctx}
}

func (r *SuperheroRepository) FindById(id int) (*models.Superhero, error) {
	query := `
		SELECT id, superhero_name, full_name, gender_id, eye_colour_id, hair_colour_id, 
		skin_colour_id, race_id, publisher_id, alignment_id, height_cm, weight_kg
		FROM superhero.superhero
		WHERE id = $1
	`
	var superhero models.Superhero
	err := r.db.QueryRow(r.ctx, query, id).Scan(
		&superhero.ID, &superhero.SuperheroName, &superhero.FullName, &superhero.GenderID,
		&superhero.EyeColourID, &superhero.HairColourID, &superhero.SkinColourID,
		&superhero.RaceID, &superhero.PublisherID, &superhero.AlignmentID,
		&superhero.HeightCm, &superhero.WeightKg,
	)
	if err != nil {
		return nil, err
	}
	return &superhero, nil
}

func (r *SuperheroRepository) FindByIdWithDetailed(id int) (*models.SuperheroWithDetails, error) {
	query := `
		SELECT s.id, s.superhero_name, s.full_name, 
			g.gender, e.colour as eye_colour, h.colour as hair_colour,
			sk.colour as skin_colour, r.race, p.publisher_name, a.alignment,
			s.height_cm, s.weight_kg
		FROM superhero.superhero s
		LEFT JOIN superhero.gender g ON s.gender_id = g.id
		LEFT JOIN superhero.colour e ON s.eye_colour_id = e.id
		LEFT JOIN superhero.colour h ON s.hair_colour_id = h.id
		LEFT JOIN superhero.colour sk ON s.skin_colour_id = sk.id
		LEFT JOIN superhero.race r ON s.race_id = r.id
		LEFT JOIN superhero.publisher p ON s.publisher_id = p.id
		LEFT JOIN superhero.alignment a ON s.alignment_id = a.id
		WHERE s.id = $1
	`
	var superhero models.SuperheroWithDetails
	err := r.db.QueryRow(r.ctx, query, id).Scan(
		&superhero.ID, &superhero.SuperheroName, &superhero.FullName,
		&superhero.Gender, &superhero.EyeColour, &superhero.HairColour,
		&superhero.SkinColour, &superhero.Race, &superhero.Publisher,
		&superhero.Alignment, &superhero.HeightCm, &superhero.WeightKg,
	)
	if err != nil {
		return nil, err
	}
	return &superhero, nil
}
