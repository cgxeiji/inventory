#!/usr/bin/env bash

package="inventory-server"
version="latest"
platforms=("windows/amd64" "darwin/amd64" "linux/arm" "linux/amd64")

for platform in "${platforms[@]}"
do
    p_split=(${platform//\// })
    GOOS=${p_split[0]}
    GOARCH=${p_split[1]}

    output_name='bin/'$package'_'$version'_'$GOOS'-'$GOARCH

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name

    if [ $? -ne 0 ]; then
        echo 'There was an error compiling!'
        exit 1
    fi
done
