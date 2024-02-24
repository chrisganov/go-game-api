package data

import (
	"database/sql"
	"fmt"

	validator "go_game_api.com/internal/validators"
	"go_game_api.com/pkg/utils"
)

type scoreLevel string

const (
	EasyLevel   scoreLevel = "EASY"
	MediumLevel scoreLevel = "MEDIUM"
	HardLevel   scoreLevel = "HARD"
)

var scoreLevelArr = []scoreLevel{EasyLevel, MediumLevel, HardLevel}

type Score struct {
	Id        int        `json:"id"`
	Level     scoreLevel `json:"level"`
	Moves     int        `json:"moves"`
	Time      int        `json:"time"`
	UserId    int        `json:"userId"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}

type ScoreInput struct {
	Level scoreLevel `json:"level"`
	Moves int        `json:"moves"`
	Time  int        `json:"time"`
}

type ScoreModel struct {
	DB *sql.DB
}

func ValidateScoreInput(v *validator.Validator, scoreInput *ScoreInput) {
	v.Check(scoreInput.Time > 0, "time", "Time cannot be 0 or less")

	v.Check(scoreInput.Moves > 0, "moves", "Moves cannot be 0 or less")

	v.Check(scoreInput.Level != "", "level", "Level cannot be empty")
	v.Check(utils.ArrContains(scoreLevelArr, scoreInput.Level), "level", "Invalid level")
}

const scoreColumns = "id, level, moves, time, user_id, created_at, updated_at"

func (m ScoreModel) Insert(score *Score) error {
	args := []interface{}{score.Level, score.Moves, score.Time, score.UserId}

	query := fmt.Sprintf(`
		INSERT INTO scores (level, moves, time, user_id)
		values ($1, $2, $3, $4)
		RETURNING %s
	`, scoreColumns)

	err := m.DB.QueryRow(query, args...).Scan(&score.Id, &score.Level, &score.Moves, &score.Time, &score.UserId, &score.CreatedAt, &score.UpdatedAt)

	return err
}

func (m ScoreModel) GetAll() ([]Score, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM scores
	`, scoreColumns)

	rows, err := m.DB.Query(query)

	if err != nil {
		fmt.Printf("Error kur %v", err)
		return nil, err
	}

	defer rows.Close()

	var scores []Score

	for rows.Next() {
		var score Score

		err := rows.Scan(&score.Id, &score.Level, &score.Moves, &score.Time, &score.UserId, &score.CreatedAt, &score.UpdatedAt)

		if err != nil {
			// TODO
			continue
		}

		scores = append(scores, score)
	}

	return scores, nil
}
