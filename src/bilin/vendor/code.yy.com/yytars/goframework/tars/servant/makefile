CC          = go
CFLAGSG         = test -v

test: servant

adapter:
	${CC} ${CFLAGSG} AdapterProxy.go AdapterProxy_test.go AppProtocol.go 

manager: adapter
	${CC} ${CFLAGSG} EndpointManager_test.go EndpointManager.go AdapterProxy.go AppProtocol.go

obj: manager 
	${CC} ${CFLAGSG} ObjectProxy_test.go ObjectProxy.go EndpointManager.go AdapterProxy.go AppProtocol.go Message.go ServantProxy.go Communicator.go

servant: obj
	${CC} ${CFLAGSG} ServantProxy_test.go ServantProxy.go ObjectProxy.go EndpointManager.go AdapterProxy.go AppProtocol.go Message.go Communicator.go

msg:
	${CC} ${CFLAGSG}  Message.go Message_test.go ObjectProxy.go ServantProxy.go AdapterProxy.go Communicator.go EndpointManager.go Application.go ServantHandle.go Config.go AppProtocol.go AdminF.go NodeF.go
server:
	nohup /home/svn/titustian/test/tafgo/NetClip/NetClip &
