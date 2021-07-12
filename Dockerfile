FROM quay.io/prometheus/busybox:latest
LABEL maintainer="kobtea9696@gmail.com"

COPY pkg/config/sample.yml /etc/iapetus.yml
COPY dist/iapetus_linux_amd64/iapetus /bin/iapetus

EXPOSE 19090
ENTRYPOINT ["/bin/iapetus"]
CMD ["--config=/etc/iapetus.yml"]
