FROM --platform=linux/amd64 quay.io/frrouting/frr:10.3.0

COPY dockerfile/files/start_frr_base /usr/bin/start_frr

RUN chmod +x /usr/bin/start_frr

CMD ["/usr/bin/start_frr"]
