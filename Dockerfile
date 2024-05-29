FROM scratch
COPY hack/passwd /etc/passwd
COPY operator /
USER 65534
ENTRYPOINT ["/operator"]
