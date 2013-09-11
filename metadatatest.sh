#!/bin/bash
curl -X PUT http://localhost:8080/metadata/some/path/key/a/value/b
curl http://localhost:8080/metadata/some/path/key/a
curl -X DELETE http://localhost:8080/metadata/some/path/key/a
curl http://localhost:8080/metadata/some/path/key/a