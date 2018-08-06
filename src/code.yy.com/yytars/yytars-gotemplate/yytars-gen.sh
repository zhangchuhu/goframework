#!/bin/bash

if [[ -z $GOPATH ]]; then
    echo "Please set GOPATH in ~/.bashrc"
    exit 1
fi

if [[ -z $GOROOT ]]; then
    echo "Please set GOROOT in ~/.bashrc"
    exit 1
fi

if [[ $(pwd) != $GOPATH/src ]]; then
    echo -n 'cd ${GOPATH}/src: '
    cd $GOPATH/src
fi

echo 
echo "GOPATH = $GOPATH"
echo "GOROOT = $GOROOT"
echo "GOVERSION = $(go version)"

unset AppName
while [[ ! ${AppName} =~ ^[A-Za-z]+$ ]]; do
    echo "What is your AppName? [应用名只能包含英文字母]"
    read AppName
done

unset ServiceName
while [[ ! ${ServiceName} =~ ^[A-Za-z]+[A-Za-z0-9]*$ ]]; do
    echo "What is your ServiceName? [服务名只能包含英文字母，数字，并以字母开头]"
    read ServiceName
done

unset ServantName
while [[ ! ${ServantName} =~ ^[A-Z]+[A-Za-z0-9]*$ ]]; do
    echo "What is your ServantName? [服务名只能包含英文字母，数字，并以大写字母开头]"
    read ServantName
done

echo
echo "AppName = ${AppName}"
echo "ServiceName = ${ServiceName}"
echo "ServantName = ${ServantName}"
echo

echo "Are you Okay?"
select yn in "Yes" "No"; do
    case $yn in
        Yes ) break;;
        No ) exit 1;;
    esac
done

#get machine and architect
#Linux and Darwin support
unames=`uname -s`

TEMP=code.yy.com/yytars/yytars-gotemplate
mkdir -p $AppName/protocol

#copy makefile.taf to app folder
if [[ ! -e $AppName/makefile.taf ]]; then
echo 'copy makefile.taf to app folder'
cp $TEMP/makefile.taf $AppName/
fi

cat $TEMP/protocol/AppName/ServiceName/ServiceName.proto > $AppName/protocol/$ServiceName.proto
if [[ "$unames" = "Darwin" ]];then
    sed -i '' -e "s/%{AppName}/$AppName/g" $AppName/protocol/$ServiceName.proto
    sed -i '' -e "s/%{ServiceName}/$ServiceName/g" $AppName/protocol/$ServiceName.proto
    sed -i '' -e "s/%{ServantName}/$ServantName/g" $AppName/protocol/$ServiceName.proto
else
    sed -i "s/%{AppName}/$AppName/g" $AppName/protocol/$ServiceName.proto
    sed -i "s/%{ServiceName}/$ServiceName/g" $AppName/protocol/$ServiceName.proto
    sed -i "s/%{ServantName}/$ServantName/g" $AppName/protocol/$ServiceName.proto
fi

mkdir -p $AppName/$ServiceName
cp -a $TEMP/AppName/ServiceName/. $AppName/$ServiceName
mv $AppName/$ServiceName/handler/HelloWorldObj.go $AppName/$ServiceName/handler/${ServantName}Obj.go
if [[ "$unames" = "Darwin" ]];then
    find $AppName/$ServiceName -type f -exec sed -i '' -e "s/%{AppName}/$AppName/g" {} \;
    find $AppName/$ServiceName -type f -exec sed -i '' -e "s/%{ServiceName}/$ServiceName/g" {} \;
    find $AppName/$ServiceName -type f -exec sed -i '' -e "s/%{ServantName}/$ServantName/g" {} \;
else
    find $AppName/$ServiceName -type f -exec sed -i "s/%{AppName}/$AppName/g" {} \;
    find $AppName/$ServiceName -type f -exec sed -i "s/%{ServiceName}/$ServiceName/g" {} \;
    find $AppName/$ServiceName -type f -exec sed -i "s/%{ServantName}/$ServantName/g" {} \;
fi

echo "Done."
echo 'following below steps:'
echo '1 cd '${GOPATH}/src/$AppName
echo '2 dep init;dep ensure'
echo '3 cd $ServiceName'
echo '4 make;make tar'
echo '5 start publish your code from http://58.215.138.213:8080'

