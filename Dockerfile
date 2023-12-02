FROM traefik:2.10.5

COPY traefik.yml /etc/traefik/traefik.yml
COPY configurations /configurations
COPY plugins-local plugins-local