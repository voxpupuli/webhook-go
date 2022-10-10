FROM scratch
EXPOSE 4000
COPY webhook-go /webhook-go
COPY build/webhook.yml /webhook.yml
ENTRYPOINT [ "/webhook-go", "server" ]
