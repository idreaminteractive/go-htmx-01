#!/bin/sh
docker run -p 8080:8080 --rm --device /dev/fuse --cap-add SYS_ADMIN latest