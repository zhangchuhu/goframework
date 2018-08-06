/*
 * filename      : vod_tokey.h
 * descriptor    :  
 * author        : chenyuhui
 * create time   : 2014-11-28 14:41
 * modify list   :
 * +----------------+---------------+---------------------------+
 * | date           | who           | modify summary            |
 * +----------------+---------------+---------------------------+
 */
#ifndef _VOD_TOKEY_H_
#define _VOD_TOKEY_H_

#include <stdint.h>
#include <stdlib.h>

#ifndef __cplusplus
typedef unsigned char bool;
#else
extern "C" {
#endif

typedef enum TokenType {
    TYPE_LIVE = 0,
    TYPE_VOD
} TokenType;

typedef struct Token {
    char* key;
    int keylen;
} Token;

typedef struct TokenInfo {
    uint32_t appid;
    uint64_t ttl;
    uint32_t uid;

    //点播字段
    uint32_t auth;
    char*   context;

    //直播字段
    uint32_t sid;
    uint32_t audio_send_expire;
    uint32_t audio_recv_expire;
    uint32_t video_send_expire;
    uint32_t video_recv_expire;
    uint32_t text_send_expire;
    uint32_t text_recv_expire;
} TokenInfo;

bool GenToken(uint16_t secretvsn, const char* secretkey, TokenInfo* info, TokenType type, Token* token);
bool ValidateToken(Token* token, const char* secretkey);
bool GetProperty(Token* token, TokenInfo* info);

#ifdef __cplusplus
}
#endif
#endif
