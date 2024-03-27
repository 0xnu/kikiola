package db

import (
	"errors"
	"math"
)

// Vector represents a vector in the database.
type Vector struct {
	ID        string
	Embedding []float64
	Metadata  map[string]string
}

// Distance calculates the cosine similarity between two vectors.
func (v Vector) Distance(other Vector) (float64, error) {
	if len(v.Embedding) != len(other.Embedding) {
		return 0, errors.New("embedding dimensions mismatch")
	}

	dotProduct := 0.0
	normV := 0.0
	normOther := 0.0

	for i := range v.Embedding {
		dotProduct += v.Embedding[i] * other.Embedding[i]
		normV += v.Embedding[i] * v.Embedding[i]
		normOther += other.Embedding[i] * other.Embedding[i]
	}

	if normV == 0 || normOther == 0 {
		return 0, nil
	}

	return dotProduct / (math.Sqrt(normV) * math.Sqrt(normOther)), nil
}

// Normalize normalizes the vector to unit length.
func (v *Vector) Normalize() {
	norm := 0.0
	for _, value := range v.Embedding {
		norm += value * value
	}
	norm = math.Sqrt(norm)

	if norm > 0 {
		for i := range v.Embedding {
			v.Embedding[i] /= norm
		}
	}
}

// Add adds another vector to the current vector.
func (v *Vector) Add(other Vector) error {
	if len(v.Embedding) != len(other.Embedding) {
		return errors.New("embedding dimensions mismatch")
	}

	for i := range v.Embedding {
		v.Embedding[i] += other.Embedding[i]
	}

	return nil
}

// Subtract subtracts another vector from the current vector.
func (v *Vector) Subtract(other Vector) error {
	if len(v.Embedding) != len(other.Embedding) {
		return errors.New("embedding dimensions mismatch")
	}

	for i := range v.Embedding {
		v.Embedding[i] -= other.Embedding[i]
	}

	return nil
}
