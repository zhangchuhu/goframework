#include "vod_token.h"
#include "YCTokenAppSecretProvider.h"
#include "YCTokenPropertyProvider.h"
#include "YCToken.h"
#include "YCTokenBuilder.h"
#include "uri_codec.h"
#include <string>
#include <stdlib.h>
//#include <iostream>
//#include<iomanip>
using namespace yctoken;
//using namespace std;

class MyAppSecretProvider:public YCloudAppSecretProvider
{
public:
    MyAppSecretProvider(uint16_t secretvsn, const std::string &key)
        : m_secretvsn(secretvsn), m_secretkey(key)
    {
    }

    std::map<uint16_t,std::string> getAppsecret(uint32_t& appKey)
    {
        std::map<uint16_t,std::string> secretMap;
        secretMap.insert(std::pair<int,std::string>(m_secretvsn, m_secretkey));
        return secretMap;
    }
private:
    uint16_t m_secretvsn;
    std::string m_secretkey;
};

bool GenToken(uint16_t secretvsn, const char* secretkey, TokenInfo* info, TokenType type, Token* token)
{
    if (secretkey == NULL || info == NULL || token == NULL || (type == TYPE_VOD && info->context == NULL)) {
        if (token != NULL) {
            token->key = NULL;
            token->keylen = 0;
        }
        return false;
    }
    std::string strkey(secretkey);
    MyAppSecretProvider myAppsecretProvider(secretvsn, strkey); 
    YCTokenBuilder builder(myAppsecretProvider);

    uint32_t cttl = info->ttl;
    YCTokenPropertyProvider propProvider(info->appid, cttl);
    std::string uidname("UID");
    propProvider.addTokenExtendProperty(uidname, info->uid);
    if (type == TYPE_VOD) {
        std::string authname("AUTH");
        propProvider.addTokenExtendProperty(authname, info->auth);
        std::string ctxname("CONTEXT");
        std::string strctx(info->context);
        propProvider.addTokenExtendProperty(ctxname, strctx);
    } else {
        std::string name;
        name = "SID";
        propProvider.addTokenExtendProperty(name, info->sid);
        name = "AUDIO_SEND";
        propProvider.addTokenExtendProperty(name, info->audio_send_expire);
        name = "AUDIO_RECV";
        propProvider.addTokenExtendProperty(name, info->audio_recv_expire);
        name = "VIDEO_SEND";
        propProvider.addTokenExtendProperty(name, info->video_send_expire);
        name = "VIDEO_RECV";
        propProvider.addTokenExtendProperty(name, info->video_recv_expire);
        name = "TEXT_SEND";
        propProvider.addTokenExtendProperty(name, info->text_send_expire);
        name = "TEXT_RECV";
        propProvider.addTokenExtendProperty(name, info->text_recv_expire);
    }
    std::string strtoken;
    try {
        strtoken = builder.buildBinaryToken(propProvider);
    }
    catch (...) {
    }
    token->key = urlsafe_base64_encode((const unsigned char*)strtoken.c_str(), strtoken.length());
    token->keylen = strlen(token->key); 
    return true;
}

bool ValidateToken(Token* token, const char* secretkey)
{
    if (token == NULL || secretkey == NULL) {
        return false;
    }
    std::string strkey(secretkey);
    MyAppSecretProvider myAppsecretProvider(0, strkey); 
    YCTokenBuilder builder(myAppsecretProvider);

    bool ret = false;
    try {
        int declen = token->keylen;
        char* deckey = (char*)urlsafe_base64_decode(token->key, &declen);
        std::string strtoken(deckey, declen);
        free(deckey);
        YCToken* pToken = builder.validateTokenBytes(strtoken);
        ret = !pToken->isExpired();
        delete pToken;
    }
    catch (...) {}
    return ret;
}

bool GetProperty(Token* token, TokenInfo* info)
{
    info->context = NULL;
    if (token == NULL || info == NULL) {
        return false;
    }
    std::string strkey;
    MyAppSecretProvider myAppsecretProvider(0, strkey); 
    YCTokenBuilder builder(myAppsecretProvider);

    bool ret = false;
    try {
        int declen = token->keylen;
        char* deckey = (char*)urlsafe_base64_decode(token->key, &declen);
        std::string strtoken(deckey, declen);
        YCToken* pToken = builder.getTokenInfo(strtoken);
        info->appid = pToken->getAppKey();
        info->ttl = pToken->getExpireTime();
        pToken->fetchExtendPropertyValue("AUTH", info->auth);
        pToken->fetchExtendPropertyValue("UID", info->uid);
        pToken->fetchExtendPropertyValue("SID", info->sid);
        pToken->fetchExtendPropertyValue("AUDIO_SEND", info->audio_send_expire);
        pToken->fetchExtendPropertyValue("AUDIO_RECV", info->audio_recv_expire);
        pToken->fetchExtendPropertyValue("VIDEO_SEND", info->video_send_expire);
        pToken->fetchExtendPropertyValue("VIDEO_RECV", info->video_recv_expire);
        pToken->fetchExtendPropertyValue("TEXT_SEND", info->text_send_expire);
        pToken->fetchExtendPropertyValue("TEXT_RECV", info->text_recv_expire);
        std::string ctx;
        if (pToken->fetchExtendPropertyValue("CONTEXT", ctx)) {
            size_t len = ctx.length();
            info->context = (char*)malloc(len + 1);
            memset((void*)info->context, 0, len + 1);
            memcpy((void*)info->context, ctx.c_str(), len);
        } else {
            info->context = NULL;
        }
        delete pToken;
        free(deckey);
        ret = true;
    }
    //catch (YCTokenException exp) { cout << exp.errorCode() << endl; }
    catch (...) {}
    return ret;
}
