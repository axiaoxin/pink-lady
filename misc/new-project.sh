#! /usr/bin/env bash
CYAN="\033[1;36m"
GREEN="\033[0;32m"
WHITE="\033[1;37m"
NOTICE_FLAG="${CYAN}❯"
QUESTION_FLAG="${GREEN}?"

main() {
    gopath=`go env GOPATH`
    echo -e "${NOTICE_FLAG} New project will be create in ${WHITE}${gopath}/src/"
    echo -e "${NOTICE_FLAG} You should enter the project full name such like <github.com/username/projectname>"
    echo -ne "${QUESTION_FLAG} ${CYAN}Enter your new project full name${CYAN}: "
    read projname
    projname_dir=`dirname ${projname}`
    mkdir -p ${gopath}/src/${projname_dir}
    echo -ne "${QUESTION_FLAG} ${CYAN}Do you want to the demo code[${WHITE}Y/n${CYAN}]: "
    read rmdemo

    # get skeleton
    echo -e "${NOTICE_FLAG} Downloading the skeleton..."
    go get -u -d github.com/axiaoxin/pink-lady/app
    # replace project name
    echo -e "${NOTICE_FLAG} Generating the project..."
    cp -r ${gopath}/src/github.com/axiaoxin/pink-lady ${gopath}/src/${projname}
    cd ${gopath}/src/${projname} && rm -rf .git && cp ${gopath}/src/${projname}/app/config.yaml.example ${gopath}/src/${projname}/app/config.yaml
    sed -i "s|github.com/axiaoxin/pink-lady|${projname}|g"  `grep "github.com/axiaoxin/pink-lady" --include *.go -rl .`

    # remove demo
    if [ "${rmdemo}" == "n" ] || [ "${rmdemo}" == "N" ]; then
        rm -rf app/apis/demo
        rm -rf app/services/demo
        rm -rf app/models/demo
        sed -i "/demo routes start/,/demo routes end/d" app/apis/routes.go
        sed -i '/app\/apis\/demo"$/d' app/apis/routes.go
    fi
    echo -e "${NOTICE_FLAG} Create ${projname} succeed."

    # init a git repo
    echo -ne "${QUESTION_FLAG} ${CYAN}Do you want to init a git repo[${WHITE}N/y${CYAN}]: "
    read initgit
    if [ "${initgit}" == "y" ] || [ "${rmdemo}" == "Y" ]; then
        cd ${gopath}/src/${projname} && git init && git add . && git commit -m "init project with pink-lady"
        cp ${gopath}/src/${projname}/misc/pre-commit.githook ${gopath}/src/${projname}/.git/hooks/pre-commit
        chmod +x ${gopath}/src/${projname}/.git/hooks/pre-commit
    fi
}
main
