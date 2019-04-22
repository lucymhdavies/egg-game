package sprites

// Sprites in this directory of are type png
import _ "image/png"

var (
	// predefine this, so it can be used in al the generated
	// *Sprites_png.go files
	Sprites_png = make(map[string][]byte)

	// This allows me to reference images as, e.g.
	// sprites.Sprites_png["foo"]
)

// ImageMap is a generic map of strings to byte slices, meant for storing images
type ImageMap = map[string][]byte
