FROM envoyproxy/envoy:v1.26.0

ARG DIRECTORY

COPY ${DIRECTORY}/envoy.yaml /etc/envoy/envoy.yaml

CMD ["envoy", "-l", "debug", "-c", "/etc/envoy/envoy.yaml"]