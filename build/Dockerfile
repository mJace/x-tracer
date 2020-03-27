FROM sheenam3/ubuntu18.04
COPY parse/probeparser/probeparser.go /
COPY build/entrypoint.sh /
RUN chmod +x /entrypoint.sh
COPY bin/x-agent /
RUN chmod +x /x-agent
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/x-agent"]
