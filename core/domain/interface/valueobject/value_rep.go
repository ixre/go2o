/**
 * Copyright 2015 @ z3q.net.
 * name : value_rep.go
 * author : jarryliu
 * date : 2016-05-27 15:28
 * description :
 * history :
 */
package valueobject

type(
    // 微信API设置
    WxApiConfig struct {

        /**===== 微信公众平台设置 =====**/

        //APP ID
        AppId            string
        //APP 密钥
        AppSecret        string
        //通信密钥
        MpToken          string
        //通信AES KEY
        MpAesKey         string
        //原始ID
        OriId            string

        /**===== 用于微信支付 =====**/

        //商户编号
        MchId            string
        //商户接口密钥
        MchApiKey        string
        //微信支付的证书路径(上传)
        MchCertPath      string
        //微信支付的证书公钥路径(上传)
        MchCertKeyPath   string
        
        //MchPayNotifyPath string //微信支付异步通知的路径
    }



    IValueRep interface{
        // 获取微信接口配置
        GetWxApiConfig()*WxApiConfig
        // 保存微信接口配置
        SaveWxApiConfig(v *WxApiConfig)error
    }
)
