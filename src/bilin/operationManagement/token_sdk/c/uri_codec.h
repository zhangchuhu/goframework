/*
 * filename      : uri_codec.h
 * descriptor    :  
 * author        : chenyuhui
 * create time   : 2014-12-02 14:20
 * modify list   :
 * +----------------+---------------+---------------------------+
 * | date           | who           | modify summary            |
 * +----------------+---------------+---------------------------+
 */
#ifndef _URI_CODEC_H_
#define _URI_CODEC_H_

char * base64_encode(const unsigned char *value, int vlen);	
unsigned char * base64_decode(const char *value, int *rlen);
char * urlsafe_base64_encode(const unsigned char *value, int vlen);	
unsigned char * urlsafe_base64_decode(const char *value, int *rlen);

#endif
