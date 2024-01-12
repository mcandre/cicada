FROM fedora
RUN yum update && \
    yum install -y wget
