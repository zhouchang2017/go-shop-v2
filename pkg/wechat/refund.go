package wechat

import (
	"errors"
	wxpay "github.com/iGoogle-ink/gopay/wechat"
	"net/http"
	"strconv"
	"time"
)

type refundResponse struct {
	*wxpay.RefundResponse
	timestamp string
}

func (r *refundResponse) setTimeStamp() {
	r.timestamp = strconv.FormatInt(time.Now().Unix(), 10)
}

// 退款
func (this pay) Refund(opt *RefundOption) (*refundResponse, error) {
	bm, mapErr := opt.toMap()
	if mapErr != nil {
		return nil, mapErr
	}
	wxRsp, refundErr := this.Client.Refund(bm, opt.certFilePath, opt.keyFilePath, opt.pkcs12FilePath)
	if refundErr != nil {
		return nil, refundErr
	}
	_, err := wxpay.VerifySign(this.ApiKey, opt.signType, wxRsp)
	if err != nil {
		return nil, err
	}
	return &refundResponse{RefundResponse: wxRsp}, nil
}

// 退款通知
func (this pay) ParseRefundNotifyResult(req *http.Request) (refundNotify *wxpay.RefundNotify, err error) {
	notifyReq, err := wxpay.ParseRefundNotifyResult(req)
	if err != nil {
		return nil, err
	}
	if notifyReq.ReturnCode == "SUCCESS" {
		return wxpay.DecryptRefundNotifyReqInfo(notifyReq.ReqInfo, this.ApiKey)
	}
	return nil, errors.New(notifyReq.ReturnMsg)
}
