#!/bin/bash

rm *_png.go
rm *_map.go

echo
echo ----------
echo Root
echo ----------

# Base Directory
for f in *.png; do
	prefix=$(basename $f .png)
	
	file2byteslice -input $f -output ${prefix}_png.go -package sprites -var ${prefix}_png
	
	file2byteslice -input $f -output ${prefix}_Sprites_png.go -package sprites -var Sprites_png -varindex "\"${prefix}\""
done

# Subdirectories

for d in $(find . -type d -mindepth 1 -maxdepth 1); do
	map_name=$(basename $d)

	echo
	echo ----------
	echo $map_name
	echo ----------

	# Initialize the map
cat << EOF > ${map_name}_map.go
package sprites

// Sprites in this directory of are type png
import _ "image/png"

var (
	// predefine this, so it can be used in al the generated
	// *${map_name}.go files
	${map_name} = make(map[string][]byte)

	// This allows me to reference images as, e.g.
	// sprites.${map_name}["foo"]
)
EOF

	for f in $d/*png; do
		prefix=$(basename $f .png)
		file2byteslice -input $f -output ${prefix}_${map_name}_map.go -package sprites -var ${map_name} -varindex "\"${prefix}\""
	done
done
