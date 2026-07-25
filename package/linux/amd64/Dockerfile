FROM ubuntu:26.04

ARG UBUNTU_MIRROR=http://archive.ubuntu.com/ubuntu
ARG IMAGE_VERSION=dev

LABEL org.opencontainers.image.source="https://github.com/PastureStack/internal-dns" \
      org.opencontainers.image.description="Metadata-driven internal DNS service" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.version="${IMAGE_VERSION}"

RUN set -eux; \
    rm -f /usr/bin/pebble; \
    find /etc/apt -type f \( -name '*.list' -o -name '*.sources' \) -exec sed -i \
        -e "s|http://security.ubuntu.com/ubuntu|${UBUNTU_MIRROR}|g" \
        -e "s|http://archive.ubuntu.com/ubuntu|${UBUNTU_MIRROR}|g" {} +; \
    printf 'Acquire::Retries "5";\nAcquire::http::Timeout "30";\nAcquire::https::Timeout "30";\nAcquire::http::Pipeline-Depth "0";\n' > /etc/apt/apt.conf.d/80pasturestack-retries; \
    apt-get update; \
    apt-get install -y --no-install-recommends ca-certificates curl openssl tini; \
    rm -rf /var/lib/apt/lists/*; \
    mkdir -p /etc/internal-dns /licenses/cache /licenses/dependencies; \
    touch /etc/internal-dns/answers.json

COPY bin/linux/amd64/internal-dns /usr/bin/internal-dns
COPY package/linux/amd64/platform-entrypoint.sh package/linux/amd64/update-platform-ca /usr/bin/
COPY LICENSE ORIGIN.md THIRD_PARTY_NOTICES.md /licenses/
COPY cache/LICENSE cache/AUTHORS cache/CONTRIBUTORS /licenses/cache/
COPY LICENSES/ /licenses/dependencies/

RUN chmod 0755 /usr/bin/internal-dns /usr/bin/platform-entrypoint.sh /usr/bin/update-platform-ca

ENTRYPOINT ["/usr/bin/tini", "-s", "--", "/usr/bin/platform-entrypoint.sh"]
CMD ["internal-dns"]
