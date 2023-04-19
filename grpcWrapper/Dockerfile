FROM envoyproxy/envoy:v1.26.0

RUN apt-get update
COPY envoy.yaml /etc/envoy/envoy.yaml
CMD ["envoy", "-c", "/etc/envoy/envoy.yaml"]