#!/usr/bin/env bash
#
#
set -u -e -o pipefail

case "$*" in
  "-h"|"--help"|"help")
    echo "da.go -h|--help|help   -- This message"
    echo "da.go packages dirs    -- List packages in this directory."
    echo "da.go audit            -- Fix up go code."
    echo "da.go upgrade          -- Upgrade go modules."
    ;;

  "package dirs")
    # while read -r LINE ; do
    #   echo "$(dirname "$LINE")"
    # done < <(find . -mindepth 2 -maxdepth 2 -type f -iname '*.go' -and -not -path '*/tmp/*' | sort)
    {
      find . -mindepth 2 -maxdepth 2 -type f -iname '*.go' -and -not -path '*/tmp/*'
      find . -mindepth 1 -maxdepth 1 -type f -iname '*.go'
    } | sort | xargs dirname | uniq
   # find . -mindepth 1 -maxdepth 2 -type f -iname '*.go'
  ;;

  "audit")
    go mod verify
    go mod tidy
    while read -r pkg; do
      echo "=== in $pkg ==="
      gofmt -s -w ./"$pkg"
      go vet  ./"$pkg"
      go tool staticcheck ./"$pkg"
      go tool govulncheck ./"$pkg"
      echo
    done < <("$0" package dirs)

    echo -e "\e[1;32m====== Done tidying up. ======\e[0m" >&2
  ;;

  # doc: CMD upgrade
  upgrade)
    if which go-mod-upgrade &>/dev/null ; then
      go-mod-upgrade
    fi
    set -x
    go get -u all
    go get -t -u all
    go get tool
    go mod tidy
    # go list -m -u all  # lists all possible upgrades, including indirect.
    ;;

  # =============================================================================
  *)
    echo "!!! Unknown command: $*" >&2
    exit 1
    ;;
esac
