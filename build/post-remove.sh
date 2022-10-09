#!/bin/sh

remove() {
    printf "\033[32m removing webhook-go\033[0m\n"
    systemctl stop webhook-go.service
    systemctl disable webhook-go.service
    rm -rf /etc/systemd/system/webhook-go.service
    systemctl daemon-reload
}

purge() {
    printf "\033[32m Purgins config files\033[0m\n"
    rm -rf /etc/voxpupuli/webhook.yml
}

upgrade() {
    echo ""
}

echo "$@"

action="$1"

case "$action" in
  "0" | "remove")
    remove
    ;;
  "1" | "upgrade")
    upgrade
    ;;
  "purge")
    purge
    ;;
  *)
    printf "\033[32m Alpine\033[0m"
    remove
    ;;
esac

