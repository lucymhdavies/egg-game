#!/bin/bash

rm *_png.go

for f in *.png; do
	prefix=$(basename $f .png)
	
	file2byteslice -input $f -output ${prefix}_png.go -package sprites -var ${prefix}_png
	
	file2byteslice -input $f -output ${prefix}_Sprites_png.go -package sprites -var Sprites_png -varindex "\"${prefix}\""
done
