#!/bin/bash
COMMIT=$(git rev-parse --verify HEAD)
docker image build . \
  --build-arg "app_name=localaddr-dns" \
  -t "localaddr-dns:latest" \
  -t "localaddr-dns:${COMMIT}"