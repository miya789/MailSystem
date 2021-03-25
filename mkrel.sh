#!/bin/sh

progname=$1

for os in ${OSS}
do
    if [ $os = "windows" ]; then
        fullname="${progname}_${CI_COMMIT_TAG}_${os}_${GOARCH}.exe"
    else
        fullname="${progname}_${CI_COMMIT_TAG}_${os}_${GOARCH}"

    fi
	linkname="${fullname}(SHA256 $(cut -f1 -d' ' ${fullname}.sha256))"
	linkurl="${BASEURL}/jobs/${CI_JOB_ID}/artifacts/${fullname}"

    curl --insecure -v -X POST -H "PRIVATE-TOKEN: ${PRIVATE_TOKEN}" "${BASEURL}/releases/${CI_COMMIT_TAG}/assets/links" \
        -d name=${fullname} \
        -d url=${linkurl}
done
