FROM scratch
COPY hack/passwd /etc/passwd
COPY kubewg-operator /
USER 65534
ENTRYPOINT ["/kubewg-operator"]
