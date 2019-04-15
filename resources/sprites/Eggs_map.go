package sprites

// Sprites in this directory of are type png
import _ "image/png"

var (
	// predefine this, so it can be used in al the generated
	// *Eggs.go files
	Eggs = make(map[string][]byte)

	// This allows me to reference images as, e.g.
	// sprites.Eggs["foo"]
)
