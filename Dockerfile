FROM scratch
WORKDIR /opt/perfably
ADD perfably /opt/perfably
EXPOSE 8000
ENTRYPOINT ["./perfably"]
