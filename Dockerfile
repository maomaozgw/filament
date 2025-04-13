FROM ubuntu:latest

COPY ./ /filament/
EXPOSE 8000
WORKDIR /filament/
CMD ["/filament/filament", "serve"]