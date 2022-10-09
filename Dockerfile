FROM scratch
ENTRYPOINT [ "/webhook-go", "server" ]
COPY bin/webhook-go /
COPY webhook.yml.example /webhook.yml
