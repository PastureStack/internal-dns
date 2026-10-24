# syntax=docker/dockerfile:1.7
FROM ghcr.io/pasturestack/ubuntu-security-backports:v0.1.0@sha256:2742a3d429e4e14bc587f6dedc51ded4246cf90bbaedcd86c67d0885e32a6f2a AS ubuntu_security_backports
FROM ubuntu:26.04@sha256:2260313b31c8c011cd2eebe728008efac1b3982be73eb71348ea2648d2c0e09b

ARG UBUNTU_SNAPSHOT=20260825T000000Z

COPY security/install-ubuntu-security-updates.sh /usr/local/lib/pasturestack/install-ubuntu-security-updates.sh
COPY security/ubuntu-security-packages.lock /usr/share/pasturestack/ubuntu-security-packages.lock

ADD --checksum=sha256:5b6326f1932cfd7be57130561dc78192f09ed9f58bfbe6606ed0c50510c3b6cd \
    https://archive.ubuntu.com/ubuntu/pool/main/c/ca-certificates/ca-certificates_20260816ubuntu0.26.04.1_all.deb \
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
    rm -rf /tmp/ca-bootstrap; \
    rm -f /usr/bin/pebble; \
    rm -f /etc/apt/sources.list /etc/apt/sources.list.d/*.list /etc/apt/sources.list.d/*.sources; \
    printf 'Types: deb\nURIs: https://snapshot.ubuntu.com/ubuntu/%s\nSuites: resolute resolute-updates resolute-backports resolute-security\nComponents: main universe restricted multiverse\nSigned-By: /usr/share/keyrings/ubuntu-archive-keyring.gpg\nSnapshot: no\n' "${UBUNTU_SNAPSHOT}" > /etc/apt/sources.list.d/pasturestack-snapshot.sources; \
    printf 'Acquire::Retries "5";\nAcquire::http::Timeout "30";\nAcquire::https::Timeout "30";\nAcquire::http::Pipeline-Depth "0";\nAcquire::https::CaInfo "/etc/ssl/certs/ca-certificates.crt";\nAcquire::https::Verify-Peer "true";\nAcquire::https::Verify-Host "true";\nAcquire::AllowInsecureRepositories "false";\nAPT::Get::AllowUnauthenticated "false";\n' > /etc/apt/apt.conf.d/80pasturestack-snapshot; \
    apt-get update; \
    apt-get install -y --no-install-recommends /tmp/ca-certificates.deb curl openssl tini; \
    test "$(dpkg-query -W -f='${Version}' ca-certificates)" = '20260816ubuntu0.26.04.1'; \
    sh /usr/local/lib/pasturestack/install-ubuntu-security-updates.sh /usr/share/pasturestack/ubuntu-security-packages.lock /etc/ssl/certs/ca-certificates.crt; \
    dpkg-query -W -f='${binary:Package}\t${Version}\n' | LC_ALL=C sort > /usr/share/pasturestack/ubuntu-security-packages.tsv; \
    rm -f /tmp/ca-certificates.deb; \
    rm -rf /var/lib/apt/lists/*; \
    mkdir -p /etc/internal-dns /licenses/cache /licenses/dependencies; \
    touch /etc/internal-dns/answers.json

COPY bin/linux/amd64/internal-dns /usr/bin/internal-dns
COPY package/linux/amd64/platform-entrypoint.sh package/linux/amd64/update-platform-ca /usr/bin/
COPY LICENSE ORIGIN.md THIRD_PARTY_NOTICES.md /licenses/
COPY cache/LICENSE cache/AUTHORS cache/CONTRIBUTORS /licenses/cache/
COPY LICENSES/ /licenses/dependencies/

RUN chmod 0755 /usr/bin/internal-dns /usr/bin/platform-entrypoint.sh /usr/bin/update-platform-ca

# Install only qualified packages already present after the final Ubuntu update.
# This read-only build mount leaves no installer, archive or extra runtime tool.
RUN --mount=from=ubuntu_security_backports,source=/out,target=/opt/pasturestack/ubuntu-security-backports,ro \
    set -eux; \
    evidence=/usr/share/pasturestack/security/ubuntu-backports; \
    mkdir -p "${evidence}"; \
    /opt/pasturestack/ubuntu-security-backports/install-existing-backports --install-in-docker-build > "${evidence}/install.json"; \
    /opt/pasturestack/ubuntu-security-backports/install-existing-backports --verify-installed > "${evidence}/verify.json"; \
    for document in manifest.json qualification.json openvex.json source-inputs.lock.json; do \
        cp "/opt/pasturestack/ubuntu-security-backports/${document}" "${evidence}/${document}"; \
    done; \
    printf '%s\n' 'ghcr.io/pasturestack/ubuntu-security-backports:v0.1.0@sha256:2742a3d429e4e14bc587f6dedc51ded4246cf90bbaedcd86c67d0885e32a6f2a' > "${evidence}/bundle-image.txt"; \
    (cd "${evidence}" && sha256sum manifest.json qualification.json openvex.json source-inputs.lock.json install.json verify.json bundle-image.txt > SHA256SUMS); \
    dpkg-query -W -f='${binary:Package}\t${Version}\n' | LC_ALL=C sort > /usr/share/pasturestack/ubuntu-security-packages.tsv; \
    rm -f /var/cache/ldconfig/aux-cache

ENTRYPOINT ["/usr/bin/tini", "-s", "--", "/usr/bin/platform-entrypoint.sh"]
CMD ["internal-dns"]
