alias ll="ls -alF"

alias fl=fleetctl
alias flu="fleetctl list-units"
alias fld="fleetctl destroy"
alias flm="fleetctl list-machines -full"
alias flf="fleetctl list-unit-files"

function flr() {
  fleetctl unload $*
  fleetctl start $*
}

alias watch="watch -n 1"
alias w="watch"

alias wflu="watch fleetctl list-units"
alias wflf="watch fleetctl list-unit-files"
alias wflm="watch fleetctl list-machines"

alias wdf="watch df -h"

alias wdps="watch docker ps -a --no-trunc"
alias wdimg="watch docker images"

alias htop="toolbox htop"

# systemctl
alias sy=systemctl
alias syr="sudo systemctl restart"
alias sys="systemctl status"
alias jou="journalctl -u"
alias jout="journalctl --since today -u"
alias jot="journalctl --since today"
alias jof="journalctl -fu"
alias jok="journalctl -k"
alias jfk="journalctl -fk"
alias jf="journalctl -f"

# docker
alias dr="docker run --rm -ti" 
alias dk="docker kill"
alias dp="docker push"
function db() {
  docker build -t $1 .
}
alias de="docker exec -ti"
alias dlg="docker logs"
alias dlf="docker logs -f"

function des() {
  docker exec -ti $1 /bin/sh
}
function deb() {
  docker exec -ti $1 /bin/bash
}

alias clps="docker ps -aq|xargs docker rm"
alias climg="docker images -q|xargs docker rmi"
alias clall="clps;climg"

alias dus="du -sh"

function drs() {
  docker run -ti --rm --entrypoint=/bin/sh $*
}

function drb() {
  docker run -ti --rm --entrypoint=/bin/bash $*
}


alias ssh="ssh -v"


# etcdctl
alias etl="etcdctl"
alias etllr="etcdctl ls --recursive"
alias etldr="etcdctl rm --recursive"

alias http="docker run -it --rm clue/httpie"

PROMPT=""
PROMPT+='%(?.`if [ $timer_show != 0 ]
then echo "😄  ${timer_show}s\n "
fi`.😡  😡  😡  (%?:${timer_show}s%) 😡  😡  😡 
)
'

. /etc/environment

export PS1="\[\033[01;32m\][\H:\w]\[\033[0m\]\[\033[01;34m\] \$\[\033[00m\] "


alias dfh="df -h"

# DNS
alias dig="docker run wkruse/dnsutils dig"
