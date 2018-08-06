#! /bin/bash
destpath=${GOPATH}/src/code.yy.com/yytars/goframework

# update

#get machine and architect
#Linux and Darwin support
unames=`uname -s`
#x86_64 only support
unamem=`uname -m`

if [[ "$unames" = "Linux" || "$unames" = "Darwin" ]];then
  echo 'setup env for '$unames
else
  echo 'only Linux or Darwin support'
  exit
fi

# git check
if [[ -z `which git` ]];then
  echo "please install git first"
  exit
fi

# go env setup
whichgo=`which go`

if [ -z ${whichgo} ];then
  if [ ! -f goinstall.sh ];then
    echo 'contact 909012229 to get goinstall.sh'
    exit
  fi
  echo 'use goinstall.sh to install golang first, then run setup.sh again'
  if [[ "$unames" = "Linux" && "$unamem" = "x86_64" ]];then
    echo 'for linux 64bit:bash goinstall.sh --64'
    exit
  fi
  
  if [[ "$unames" = "Linux" && "$unamem" = "i686" ]];then
    echo 'for linux 64bit:bash goinstall.sh --32'
    exit
  fi
  
  if [[ "$unames" = "Darwin" && "$unamem" = "x86_64" ]];then
    echo 'for mac 64bit:bash goinstall.sh --darwin'
  fi
  
  echo "only Linux 64/32 or Darwin 64 support"
  exit
fi

## go version check
goversion=`go version | awk '{print $3}' | sed 's/go//g'`
OLD_IFS="$IFS" ; IFS="." ; versionarr=($goversion) ; IFS="$OLD_IFS";
if [ ${versionarr[0]} -lt 1 ] ; then
	echo 'your golang version is lower than go1.9, please remove it and try again'
	exit
elif [ ${versionarr[0]} -eq 1 ] && [ ${versionarr[1]} -lt 9 ] ; then
	echo 'your golang version is lower than go1.9, please remove it and try again'
	exit
fi

# check go env 
if test -z ${GOPATH} -o -z ${GOROOT}
then
  echo 'GOPATH or GOROOT env not set,please declare in .bash_profile'
  exit
else
  echo 'GOPATH='${GOPATH}
  echo 'GOROOT='${GOROOT}
fi

#check GOPATH/bin and src

if [ ! -d ${GOPATH}/bin ];then
  echo 'mkdir '${GOPATH}/bin
  mkdir ${GOPATH}/bin
fi

if [ ! -d ${GOPATH}/src ];then
  echo 'mkdir '${GOPATH}/src
  mkdir ${GOPATH}/src
fi

# setup protoc-gen-tars
echo 'setup protoc-gen-tars'
go get -u code.yy.com/yytars/protoc-gen-tars
go install code.yy.com/yytars/protoc-gen-tars

#check pb3
if test -e ${GOPATH}/bin/protoc3
then
  echo 'protoc3 ready, version:'
  echo `protoc3 --version`
else
  if [[ "$unames" = "Linux" && "$unamem" = "x86_64" ]];then
      cp ${GOPATH}/src/code.yy.com/yytars/protoc-gen-tars/tools/protoc-3.5.1-linux-x86_64/bin/protoc ${GOPATH}/bin/protoc3
      chmod +x ${GOPATH}/bin/protoc3
  fi

  if [[ "$unames" = "Darwin" && "$unamem" = "x86_64" ]];then
      cp ${GOPATH}/src/code.yy.com/yytars/protoc-gen-tars/tools/protoc-3.5.1-osx-x86_64/bin/protoc ${GOPATH}/bin/protoc3
      chmod +x ${GOPATH}/bin/protoc3
  fi
fi

#setup yytars-gen
cp yytars-gen.sh ${GOPATH}/bin/

echo '********************************************'
echo 'Cheers!!!, env ready'
