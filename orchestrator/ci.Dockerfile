FROM alpine:edge

COPY orchestrator/target/release/orchestrator /usr/bin/orchestrator
COPY orchestrator/startup.sh startup.sh

CMD sh startup.sh