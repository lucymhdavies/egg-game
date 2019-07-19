# In Netlify, this is running in Docker: netlify/build:xenial
# https://github.com/netlify/build-image/blob/xenial/Dockerfile

# TODO: handle GOPATH not being set
# so that I can get this to work locally

.PHONY: build_js
build_js:
	go install github.com/gopherjs/gopherjs
	go get -u
	${GOPATH}/bin/gopherjs build -o public/egg.js



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
