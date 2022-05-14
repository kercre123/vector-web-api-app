# wire-vector-custom-api

This is the source for an HTTP API designed to run on a dev/OSKR-unlocked Vector, as well as an HTML app to communicate with it easily.

It communicates with both Vector's HTTPS SDK API and webViz.

## Building

The build system is from Digital Dream Lab's [vector-cloud](https://github.com/digital-dream-labs/vector-cloud)

Install golang and docker

`git clone https://github.com/kercre123/wire-vector-custom-api`

`cd wire-vector-custom-api`

`make docker-builder`

`make custom-web`

To install it onto Vector, put ./build/custom-web into Vector's /bin, then everything else goes into their corresponsing folder.

./sbin/custom-web-interface goes to Vector's /sbin/, etc

## Current Efforts

The main current effort is to implement all SDK and webViz communication in golang.

We are also working on making the web app more responsive and more visually pleasing.

## Credits

[xanathon](https://github.com/xanathon) for making the web app look much better

[digital-dream-labs](https://github.com/digital-dream-labs) for creating OSKR and keeping Vector alive
