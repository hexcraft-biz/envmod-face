package face

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
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
func (e Face) RecognizeSingle(jpegbytes []byte) (*face.Face, her.Error) {
	faces, err := e.Recognizer.Recognize(jpegbytes)
	if err != nil {
		return nil, her.NewError(http.StatusInternalServerError, err, nil)
	} else if len(faces) != 1 {
		return nil, her.NewErrorWithMessage(http.StatusBadRequest, "Can only have one face in the image", nil)
	}

	return &faces[0], nil
}

func (e Face) Recognize(jpegbytes []byte) ([]face.Face, her.Error) {
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

// ================================================================
//
// ================================================================
const (
	DimensionCount    = 128
	FaceDistThreshold = 0.15
)

type Descriptor face.Descriptor

func (d Descriptor) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, d); err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

func (d *Descriptor) Scan(src any) error {
	if src != nil {
		buf := bytes.NewReader(src.([]byte))
		return binary.Read(buf, binary.LittleEndian, d)
	}

	return nil
}

func (d Descriptor) DistWithFace(f *Descriptor) float64 {
	sum, diff := float64(0), float64(0)
	for i := 0; i < DimensionCount; i += 1 {
		diff = float64(d[i] - (*f)[i])
		sum += diff * diff
	}

	return sum
}

// ================================================================
type Threshold float64

func (t *Threshold) Validate() {
	if *t == 0.0 {
		*t = FaceDistThreshold
	} else if *t < 0.0 {
		*t = 0.01
	} else if *t > 0.99 {
		*t = 0.99
	}
}
