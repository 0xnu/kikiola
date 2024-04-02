package db

import (
	"errors"
	"math"
)

type Vector struct {
	ID                 string
	Embedding          []float64
	Metadata           map[string]string
	Text               string
	Object             []byte
	Compressed         bool
	QuantizationParams *QuantizationParams
	PruningMask        []bool
	SparseIndices      []int
}

func (v Vector) Distance(other Vector) (float64, error) {
	if v.Compressed != other.Compressed {
		return 0, errors.New("cannot calculate distance between compressed and uncompressed vectors")
	}

	if v.Compressed {
		if len(v.Embedding) != len(other.Embedding) {
			return 0, errors.New("embedding dimensions mismatch")
		}
		dotProduct := 0.0
		normV := 0.0
		normOther := 0.0
		for i := range v.Embedding {
			vVal := v.QuantizationParams.Dequantize(v.Embedding[i])
			otherVal := other.QuantizationParams.Dequantize(other.Embedding[i])
			dotProduct += vVal * otherVal
			normV += vVal * vVal
			normOther += otherVal * otherVal
		}
		if normV == 0 || normOther == 0 {
			return 0, nil
		}
		return dotProduct / (math.Sqrt(normV) * math.Sqrt(normOther)), nil
	}

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

func (v *Vector) Normalize() {
	if !v.Compressed {
		v.normalizeUncompressed()
		return
	}

	if v.QuantizationParams != nil {
		v.normalizeQuantized()
	} else if len(v.PruningMask) > 0 {
		v.normalizePruned()
	} else if len(v.SparseIndices) > 0 {
		v.normalizeSparse()
	}
}

func (v *Vector) normalizeUncompressed() {
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

func (v *Vector) normalizeQuantized() {
	dequantizedEmbedding := make([]float64, len(v.Embedding))
	for i := range v.Embedding {
		dequantizedEmbedding[i] = v.QuantizationParams.Dequantize(v.Embedding[i])
	}
	norm := 0.0
	for _, value := range dequantizedEmbedding {
		norm += value * value
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range dequantizedEmbedding {
			dequantizedEmbedding[i] /= norm
		}
	}
	for i := range v.Embedding {
		v.Embedding[i] = v.QuantizationParams.Quantize(dequantizedEmbedding[i])
	}
}

func (v *Vector) normalizePruned() {
	norm := 0.0
	for i, value := range v.Embedding {
		if !v.PruningMask[i] {
			norm += value * value
		}
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range v.Embedding {
			if !v.PruningMask[i] {
				v.Embedding[i] /= norm
			}
		}
	}
}

func (v *Vector) normalizeSparse() {
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

func (v *Vector) Add(other Vector) error {
	if v.Compressed != other.Compressed {
		return errors.New("cannot add compressed and uncompressed vectors")
	}

	if len(v.Embedding) != len(other.Embedding) {
		return errors.New("embedding dimensions mismatch")
	}

	for i := range v.Embedding {
		v.Embedding[i] += other.Embedding[i]
	}

	return nil
}

func (v *Vector) Subtract(other Vector) error {
	if v.Compressed != other.Compressed {
		return errors.New("cannot subtract compressed and uncompressed vectors")
	}

	if len(v.Embedding) != len(other.Embedding) {
		return errors.New("embedding dimensions mismatch")
	}

	for i := range v.Embedding {
		v.Embedding[i] -= other.Embedding[i]
	}

	return nil
}

func (v *Vector) Quantize(params QuantizationParams) {
	for i := range v.Embedding {
		v.Embedding[i] = params.Quantize(v.Embedding[i])
	}
	v.Compressed = true
	v.QuantizationParams = &params
}

func (v *Vector) Prune(threshold float64) {
	v.PruningMask = make([]bool, len(v.Embedding))
	for i := range v.Embedding {
		if math.Abs(v.Embedding[i]) < threshold {
			v.Embedding[i] = 0
			v.PruningMask[i] = true
		}
	}
	v.Compressed = true
}

func (v *Vector) ToSparse() {
	var sparseEmbedding []float64
	var sparseIndices []int
	for i, value := range v.Embedding {
		if value != 0 {
			sparseEmbedding = append(sparseEmbedding, value)
			sparseIndices = append(sparseIndices, i)
		}
	}
	v.Embedding = sparseEmbedding
	v.SparseIndices = sparseIndices
	v.Compressed = true
}

type QuantizationParams struct {
	Min  float64
	Max  float64
	Bits int
}

func (p QuantizationParams) Quantize(value float64) float64 {
	normalizedValue := (value - p.Min) / (p.Max - p.Min)
	quantizedValue := math.Round(normalizedValue * (math.Pow(2, float64(p.Bits)) - 1))
	scaledValue := quantizedValue/math.Pow(2, float64(p.Bits))*p.Max - p.Min
	return scaledValue
}

func (p QuantizationParams) Dequantize(value float64) float64 {
	normalizedValue := (value - p.Min) / (p.Max - p.Min)
	dequantizedValue := normalizedValue * (math.Pow(2, float64(p.Bits)) - 1)
	scaledValue := dequantizedValue/(math.Pow(2, float64(p.Bits)))*(p.Max-p.Min) + p.Min
	return scaledValue
}
