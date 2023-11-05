package face

import (
	"os"

	"github.com/Kagami/go-face"
)

// ================================================================
//
// ================================================================
type Face struct {
	*face.Recognizer
	DirFaceRecognizationModels string
}

func New() (*Face, error) {
	return &Face{
		DirFaceRecognizationModels: os.Getenv("DIR_FACE_RECOGNIZATION_MODELS"),
	}, nil
}

// ================================================================
//
// ================================================================
func (e *Face) Open() error {
	var err error
	e.Close()
	e.Recognizer, err = face.NewRecognizer(e.DirFaceRecognizationModels)
	return err
}

func (e *Face) Close() {
	if e.Recognizer != nil {
		e.Recognizer.Close()
	}
}
