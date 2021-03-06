#!/bin/sh

./cc-test-reporter before-build 
for pkg in $(go list ./... | grep -v main); do
    go test -coverprofile=$(echo $pkg | tr / -).cover $pkg
done
echo "mode: set" > c.out
grep -h -v "^mode:" ./*.cover >> c.out
rm -f *.cover

./cc-test-reporter after-build

