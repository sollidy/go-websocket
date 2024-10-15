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
