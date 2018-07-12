package builder

import (
	"reflect"
	"testing"
)

func TestParseSrcset(t *testing.T) {

	expectations := map[string]Srcset{
		"image.jpg": Srcset{
			ImageCandidate{"image.jpg", 0, 0},
		},
		"image.jpg,": Srcset{
			ImageCandidate{"image.jpg", 0, 0},
		},
		"image.jpg, image2.jpg     1.2x": Srcset{
			ImageCandidate{"image.jpg", 0, 0},
			ImageCandidate{"image2.jpg", 0, 1.2},
		},
		"image.jpg 11w, image.jpg 22w": Srcset{
			ImageCandidate{"image.jpg", 11, 0},
			ImageCandidate{"image.jpg", 22, 0},
		},
		"image.jpg 1.1x,      image.jpg 2.2x": Srcset{
			ImageCandidate{"image.jpg", 0, 1.1},
			ImageCandidate{"image.jpg", 0, 2.2},
		},
		"    image.jpg 1.1x     ,      image.jpg     2.2x    ,": Srcset{
			ImageCandidate{"image.jpg", 0, 1.1},
			ImageCandidate{"image.jpg", 0, 2.2},
		},
		"    image.jpg 1.1x     ,      image.jpg     2w    ": Srcset{
			ImageCandidate{"image.jpg", 0, 1.1},
			ImageCandidate{"image.jpg", 2, 0},
		},
	}

	for in, expOut := range expectations {
		out := parseSrcset(in)
		if !reflect.DeepEqual(out, expOut) {
			t.Errorf("\nunexpected result for parsing '%s'\nexpecting: %+v\n received: %+v\n", in, expOut, out)
		}
	}

}

func TestSrcsetMin(t *testing.T) {
	expectations := map[string]ImageCandidate{
		"":                                            ImageCandidate{},
		"image.jpg":                                   ImageCandidate{"image.jpg", 0, 0},
		"image.jpg, image2.jpg 22w":                   ImageCandidate{"image2.jpg", 22, 0},
		"image.jpg, image2.jpg 22w, image2.jpg 33w":   ImageCandidate{"image2.jpg", 22, 0},
		"image.jpg, image2.jpg 2.2x, image2.jpg 3.3x": ImageCandidate{"image2.jpg", 0, 2.2},
	}

	for in, expOut := range expectations {
		ss := parseSrcset(in)
		out := ss.Min()
		if !reflect.DeepEqual(out, expOut) {
			t.Errorf("\nunexpected result for min of srcset '%s'\nexpecting: %+v\n received: %+v\n", in, expOut, out)
		}
	}
}

func TestGetFileName(t *testing.T) {
	expectations := map[ImageCandidate]string{
		ImageCandidate{}:                                "",
		ImageCandidate{"image.jpg", 0, 0}:               "image.jpg",
		ImageCandidate{"/somewhere/image.jpg", 0, 0}:    "image.jpg",
		ImageCandidate{"/somewhere/image.jpg", 10, 0}:   "image-10w.jpg",
		ImageCandidate{"/somewhere/image.jpg", 0, 1.1}:  "image-1.10x.jpg",
		ImageCandidate{"/somewhere/image.jpg", 10, 1.1}: "image-10w.jpg", // ignore x if we have w
	}

	for in, expOut := range expectations {
		out := in.GetFileName()
		if out != expOut {
			t.Errorf("\nunexpected result for file name for image candidate '%+v'\nexpecting: %+v\n received: %+v\n", in, expOut, out)
		}
	}
}
