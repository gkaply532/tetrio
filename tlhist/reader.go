package tlhist

import (
	"encoding/csv"
	"io"
	"slices"
)

type Reader struct {
	csv *csv.Reader
	rec Record
}

func NewReader(r io.Reader) Reader {
	reader := csv.NewReader(r)
	reader.ReuseRecord = true
	return Reader{csv: reader}
}

func (r Reader) NextRaw() ([]string, error) {
	return r.csv.Read()
}

// Next reads, parses and returns the next Record. The memory of the returned
// Record is reused between calls to Next.
func (r Reader) Next() (*Record, error) {
	fields, err := r.csv.Read()
	if err != nil {
		return nil, err
	}

	err = r.rec.Unslice(fields)
	if err != nil {
		return nil, err
	}

	return &r.rec, nil
}

func (r Reader) All(skipHeader bool) ([]Record, error) {
	if skipHeader {
		_, err := r.NextRaw()
		if err != nil {
			return nil, err
		}
	}

	rec, err := r.Next()
	if err != nil {
		return nil, err
	}

	result := make([]Record, 0, rec.GamesPlayed)
	result = append(result, *rec)

	for {
		rec, err = r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		result = append(result, *rec)
	}

	result = slices.Clip(result)
	return result, nil
}
