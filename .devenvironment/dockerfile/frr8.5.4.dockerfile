FROM --platform=linux/amd64 quay.io/frrouting/frr:8.5.4

COPY dockerfile/files/start_frr_base /usr/bin/start_frr
COPY dockerfile/files/run.sh /usr/bin/run_hmr.sh

ENV PATH=$PATH:/app/tmp

RUN chmod +x /usr/bin/start_frr
RUN chmod +x /usr/bin/run_hmr.sh

CMD ["/usr/bin/start_frr"]
