CURPATH:=$(shell pwd)
OBJDIR:=obj
LIBDIR:=lib

all: c go

c: sdk $(LIBDIR)/libvod_token_c.so

$(LIBDIR)/libvod_token_c.so: $(OBJDIR)/vod_token.o $(OBJDIR)/uri_codec.o
	g++ -shared -o $@ $^

$(OBJDIR)/vod_token.o: c/vod_token.cpp
	g++ -g -fPIC -I ./ycloud_token/include -c -o $@ $^

$(OBJDIR)/uri_codec.o: c/uri_codec.c
	g++ -fPIC -I ./ycloud_token/include -c -o $@ $^

go: c $(LIBDIR)/vod_token.a

$(LIBDIR)/vod_token.a:
	GOPATH=$(GOPATH):$(CURPATH)/go;go build -o $@ vod_token

sdk: $(LIBDIR)/libycloud_token.so

$(LIBDIR)/libycloud_token.so:
	cd ycloud_token/Release;make
	mv -f ycloud_token/Release/libycloud_token.so $@

clean:
	rm -rf $(OBJDIR)/*
	rm -rf $(LIBDIR)/*
	cd ycloud_token/Release;make clean
	make clean -f test/makefile

test::
	make -f test/makefile
example::
	make -f example/makefile
