#!/bin/sh -v
export GOPROXY=https://goproxy.cn
NPWD=`pwd`


cd Cinder/plugin/navmesh/gonavmesh
make clean
make 
cd -

cd Cinder/plugin/physxgo/physxcwrap
make clean
make
cd -


cd $NPWD/../Cinder/Agent
go build -o ../../../bin
cd $NPWD/../Cinder/Login
go build -o ../../../bin
cd $NPWD/../Daisy/DBAgent
go build -o ../../../bin
cd $NPWD/../Daisy/Game
go build -o ../../../bin
cd $NPWD/../Daisy/Battle
go build -o ../../../bin
