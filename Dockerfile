FROM scratch
COPY hack/passwd /etc/passwd
COPY kubewg /
USER 65534
ENTRYPOINT ["/kubewg"]
