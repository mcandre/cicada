#
# An image with wget
#

FROM library/ubuntu:hirsute as builder
RUN apt-get update && \
    apt-get install -y wget

FROM builder
