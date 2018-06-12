#!/usr/bin/env bash

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

rm -fr *.go testdata

eval $( go env )

cp -r "$GOROOT/src/encoding/json/"* .

sed -i.bak 's|"json:|"jdb:|g' *.go
sed -i.bak 's|`json:|`jdb:|g' *.go
sed -i.bak 's|encoding/json|github.com/silas/jdb/internal/json|g' *.go
sed -i.bak 's|parseTag|ParseTag|g' *.go
sed -i.bak 's|sf.Tag.Get("json")|sf.Tag.Get("jdb")|g' encode.go
sed -i.bak 's|name, opts := ParseTag(tag)|name, opts := ParseTag(tag);if strings.HasPrefix(name, "-") {continue}|g' encode.go
find . -name '*.bak' -delete

go fmt github.com/silas/jdb/internal/json/...
