package face

import (
	"net/http"
	"os"

	"github.com/Kagami/go-face"
	"github.com/hexcraft-biz/her"
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
func (e Face) RecognizeSingle(jpegbytes []byte) (*face.Face, *her.Error) {
	faces, err := e.Recognizer.Recognize(jpegbytes)
	if err != nil {
		return nil, her.NewError(http.StatusInternalServerError, err, nil)
	} else if len(faces) != 1 {
		return nil, her.NewErrorWithMessage(http.StatusBadRequest, "Can only have one face in the image", nil)
	}

	return &faces[0], nil
}

func (e Face) Recognize(jpegbytes []byte) ([]face.Face, *her.Error) {
	faces, err := e.Recognizer.Recognize(jpegbytes)
	return faces, her.NewError(http.StatusInternalServerError, err, nil)
}

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
