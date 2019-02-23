#!/bin/bash

for f in *.png; do
	prefix=$(basename $f .png)
	
	file2byteslice -input $f -output ${prefix}_png.go -package sprites -var ${prefix}_png
done
