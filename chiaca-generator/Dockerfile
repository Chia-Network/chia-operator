FROM registry.gitlab.com/cmmarslender/kubectl-helm:v3 AS kubectl
FROM ghcr.io/chia-network/chia-tools:latest AS chia-tools

FROM debian:bookworm-slim

COPY --from=kubectl /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=chia-tools /chia-tools /usr/local/bin/chia-tools
COPY start.sh /start.sh
CMD ["/bin/bash", "/start.sh"]
