#!/bin/bash

VERSION="$(hugo version)"
EXPECTED="Hugo Static Site Generator v0.73.0-428907CC/extended linux/amd64*"

if [[ $VERSION == $EXPECTED ]]; then
    echo -e "PASSED: Correct hugo version."
    echo "  want: ${EXPECTED}"
    echo "  got : ${VERSION}"
else
    echo "FAILED: Wrong hugo version."
    echo "  want: ${EXPECTED}"
    echo "  got : ${VERSION}"
    exit 1
fi

