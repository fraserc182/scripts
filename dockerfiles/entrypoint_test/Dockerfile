FROM alpine:latest

# Update the package repositories, upgrade installed packages, install extra packages, and clean up the package cache in one step
RUN apk update && \
    apk upgrade && \
    apk add --no-cache \
    bash

# Copy necessary files and make the entrypoint executable
COPY --chmod=0755 ito-entrypoint.sh /docker-entrypoint.d/00-ito-entrypoint.sh

# Set the default entrypoint
ENTRYPOINT ["/docker-entrypoint.d/00-ito-entrypoint.sh"]
