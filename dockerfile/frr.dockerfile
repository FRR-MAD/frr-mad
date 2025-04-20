FROM quay.io/frrouting/frr:8.5.4

COPY dockerfile/files/start_frr_base /usr/bin/start_frr

RUN chmod +x /usr/bin/start_frr

CMD ["/usr/bin/start_frr"]
