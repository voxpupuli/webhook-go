#!/bin/sh

cleanInstall() {
    printf "\033[32m Post Install of an clean install\033[0m\n"
    printf "\033[32m Reload the service unit from disk\033[0m\n"
    systemctl daemon-reload ||:
    printf "\033[32m Unmask the service\033[0m\n"
    systemctl unmask webhook-go.service ||:
    printf "\033[32m Set the preset flag for the service unit\033[0m\n"
    systemctl preset webhook-go.service ||:
    printf "\033[32m Start the service with systemctl start webhook-go.service\033[0m\n"
    printf "\033[32m Enable the service with systemctl enable webhook-go.service\033[0m\n"
}

upgrade() {
    printf "\033[32m Post Install of an upgrade\033[0m\n"
}

# Step 2, check if this is a clean install or an upgrade
action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  action="install"
elif [ "$1" = "configure" ] && [ -n "$2" ]; then
  action="upgrade"
fi

case "$action" in
  "1" | "install")
    cleanInstall
    ;;
  "2" | "upgrade")
    printf "\033[32m Post Install of an upgrade\033[0m\n"
    upgrade
    ;;
  *)
    # $1 == version being installed
    printf "\033[32m Alpine\033[0m"
    cleanInstall
    ;;
esac

