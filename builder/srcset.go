package builder

import (
	"fmt"
	"math"
	"path"
	"regexp"
	"strconv"
	"strings"
)

/*
http://w3c.github.io/html/semantics-embedded-content.html#element-attrdef-img-srcset

The src attribute must be present, and must contain a valid non-empty URL potentially surrounded by spaces referencing a non-interactive, optionally animated, image resource that is neither paged nor scripted.

The srcset attribute may also be present. If present, its value must consist of one or more image candidate strings, each separated from the next by a U+002C COMMA character (,). If an image candidate string contains no descriptors and no space characters after the URL, the following image candidate string, if there is one, must begin with one or more space characters.

An image candidate string consists of the following components, in order, with the further restrictions described below this list:
1. Zero or more space characters.
2. A valid non-empty URL that does not start or end with a U+002C COMMA character (,), referencing a non-interactive, optionally animated, image resource that is neither paged nor scripted.
3. Zero or more space characters.
4. Zero or one of the following:
  - A width descriptor, consisting of: a space character, a valid non-negative integer giving a number greater than zero representing the width descriptor value, and a U+0077 LATIN SMALL LETTER W character.
  - A pixel density descriptor, consisting of: a space character, a valid floating-point number giving a number greater than zero representing the pixel density descriptor value, and a U+0078 LATIN SMALL LETTER X character.
5. Zero or more space characters.
*/

// ImageCandidate within a srcset html attribute
type ImageCandidate struct {
	URL               string
	WidthDescriptor   int
	DensityDescriptor float32
}

// GetFileName to be used for the specific image candidate
func (ic ImageCandidate) GetFileName() string {
	if ic.URL == "" {
		return ""
	}

	_, fullName := path.Split(ic.URL)
	ext := path.Ext(fullName)
	name := fullName[:len(fullName)-len(ext)]

	if ic.WidthDescriptor != 0 {
		return fmt.Sprintf("%s-%dw%s", name, ic.WidthDescriptor, ext)
	}

	if ic.DensityDescriptor != 0 {
		return fmt.Sprintf("%s-%.2fx%s", name, ic.DensityDescriptor, ext)
	}

	return fullName
}

// Srcset for an image or source html element
type Srcset []ImageCandidate

var matchExtraWhiteSpace = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

func parseSrcset(srcsetStr string) (ss Srcset) {

	ics := strings.Split(srcsetStr, ",")

	for _, ic := range ics {

		s := strings.Trim(ic, " ")

		if s == "" {
			continue
		}

		s = matchExtraWhiteSpace.ReplaceAllString(s, " ")

		sSplit := strings.Split(s, " ")

		if len(sSplit) <= 1 || len(sSplit) > 2 {
			ss = append(ss, ImageCandidate{s, 0, 0})
			continue
		}

		src, desc := sSplit[0], sSplit[1]

		switch desc[len(desc)-1] {
		case 'w':
			if wd, err := strconv.Atoi(desc[:len(desc)-1]); err == nil {
				ss = append(ss, ImageCandidate{src, wd, 0})
			}
		case 'x':
			if wd, err := strconv.ParseFloat(desc[:len(desc)-1], 32); err == nil {
				ss = append(ss, ImageCandidate{src, 0, float32(wd)})
			}
		}

		// ignore other descriptors
	}

	return
}

// Min image candidate (smallest width or pixel dencity)
// if there are width descriptors present, pixel dencity descriptors are ignored
func (ss Srcset) Min() ImageCandidate {
	if len(ss) == 0 {
		return ImageCandidate{}
	}

	w, d := math.MaxInt32, float32(math.MaxFloat32)
	var minByW, minByD ImageCandidate
	for _, ic := range ss {
		if ic.WidthDescriptor > 0 && ic.WidthDescriptor < w {
			w = ic.WidthDescriptor
			minByW = ic
		}
		if ic.DensityDescriptor > 0 && ic.DensityDescriptor < d {
			d = ic.DensityDescriptor
			minByD = ic
		}
	}

	if w != math.MaxInt32 {
		return minByW
	}

	if d != math.MaxFloat32 {
		return minByD
	}

	return ss[0]
}
