#!/bin/bash

VERSION="$(hugo version)"
EXPECTED="hugo v0.92.1-85E2E862+extended linux/amd64*"

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

