package data

import (
	"errors"
)
var (
    ErrRecordNotFound = errors.New("record not found")
)
type Models struct {
    // Set the Movies field to be an interface containing the methods that both the
    // 'real' model and mock model need to support.
    Movies interface {
        Insert(movie *Movie) error
        Get(id int64) (*Movie, error)
        Update(movie *Movie) error
        Delete(id int64) error
    }
}

// // Create a helper function which returns a Models instance containing the mock models
// // only.
// func NewMockModels() Models {
// 	return Models{
// 			Movies: MockMovieModel{},
// 	}
// }