# syntax=docker/dockerfile:1.7
FROM ubuntu:26.04@sha256:2260313b31c8c011cd2eebe728008efac1b3982be73eb71348ea2648d2c0e09b

ARG UBUNTU_SNAPSHOT=20260825T000000Z

ADD --checksum=sha256:6077d27c6b6f8b23590cb01ff877ed8c804a67a5442cc32b5a33da10d2bd0e90 \
    https://snapshot.ubuntu.com/ubuntu/20260825T000000Z/pool/main/c/ca-certificates/ca-certificates_20260601~26.04.1_all.deb \
    /tmp/ca-certificates.deb
ARG IMAGE_VERSION=dev

LABEL org.opencontainers.image.source="https://github.com/PastureStack/internal-dns" \
      org.opencontainers.image.description="Metadata-driven internal DNS service" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.version="${IMAGE_VERSION}"

RUN set -eux; \
    mkdir -p /tmp/ca-bootstrap /etc/ssl/certs; \
    dpkg-deb --extract /tmp/ca-certificates.deb /tmp/ca-bootstrap; \
    find /tmp/ca-bootstrap/usr/share/ca-certificates -type f -name '*.crt' | LC_ALL=C sort | while IFS= read -r certificate; do sed -e '$a\' "${certificate}"; done > /etc/ssl/certs/ca-certificates.crt; \
    rm -rf /tmp/ca-bootstrap /tmp/ca-certificates.deb; \
    rm -f /usr/bin/pebble; \
    rm -f /etc/apt/sources.list /etc/apt/sources.list.d/*.list /etc/apt/sources.list.d/*.sources; \
    printf 'Types: deb\nURIs: https://snapshot.ubuntu.com/ubuntu/%s\nSuites: resolute resolute-updates resolute-backports resolute-security\nComponents: main universe restricted multiverse\nSigned-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg\nSnapshot: no\n' "${UBUNTU_SNAPSHOT}" > /etc/apt/sources.list.d/pasturestack-snapshot.sources; \
    printf 'Acquire::Retries "5";\nAcquire::http::Timeout "30";\nAcquire::https::Timeout "30";\nAcquire::http::Pipeline-Depth "0";\nAcquire::https::CaInfo "/etc/ssl/certs/ca-certificates.crt";\nAcquire::https::Verify-Peer "true";\nAcquire::https::Verify-Host "true";\nAcquire::AllowInsecureRepositories "false";\nAPT::Get::AllowUnauthenticated "false";\n' > /etc/apt/apt.conf.d/80pasturestack-snapshot; \
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
