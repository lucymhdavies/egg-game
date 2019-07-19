# In Netlify, this is running in Docker: netlify/build:xenial
# https://github.com/netlify/build-image/blob/xenial/Dockerfile

.PHONY: build_js
build_js:
	GO111MODULE=on go install github.com/gopherjs/gopherjs
	gopherjs build -o public/egg.js



#
# TODO: Other Stuff Later
#


# TODO: run
#  go run .

# TODO: run js
# go get -u github.com/gopherjs/gopherjs
# gopherjs serve


# TODO: test
# lol


# TODO: assets/sprites
# i.e. go into resources/sprites and run convert.sh

# TODO: build
# darwin / windows / linux
# ios?
# JSGO?
