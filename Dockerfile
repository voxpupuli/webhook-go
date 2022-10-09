FROM scratch
EXPOSE 4000
ENTRYPOINT [ "/webhook-go", "server" ]
COPY bin/webhook-go /
COPY webhook.yml.example /webhook.yml
