package iServer

import (
	"net/http"
	"upEletrcSign/trans"
	"io/ioutil"
	"encoding/json"
	"mygolib/modules/myLogger"
	"mygolib/defs"
	"mygolib/gerror"
)

type GetDoTransFunc func(msg *trans.TransMessage) (IDoAppTrans, error)

type IDoAppTrans interface {
	Init() gerror.IError
	DoTrans(*trans.TransMessage) (gerror.IError)
}

func IDoFunct(w http.ResponseWriter, r *http.Request, getTransFunc GetDoTransFunc) {
	ra, err := ioutil.ReadAll(r.Body)
	if err != nil {
		myLogger.Error("读取请求报文失败", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	myLogger.Infof("get Request msg:[%s]", string(ra))
	tr, gerr := trans.UnPackReq(ra)
	if gerr != nil {
		myLogger.Errorf("解析请求失败:[%s]", gerr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//gerr = trans.VerifyTransMessage(tr)
	//if gerr != nil {
	//	myLogger.Errorf("报文验证失败:[%s]", gerr)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	//myLogger.Infof("Pos_sign:%s", tr.MsgBody.Pos_sign)

	tr.MsgBody.UPos_sign, err = trans.DecodeBase64([]byte(tr.MsgBody.Pos_sign))
	if err != nil {
		RejectMsg(w, tr, defs.TRN_SYS_ERROR, err.Error())
		return
	}
	myLogger.Infof("UPos_sign:[%s]", tr.MsgBody.UPos_sign)

	TransFunc, err := getTransFunc(tr)
	if err != nil {
		RejectMsg(w, tr, defs.TRN_FORMAT_ERR, err.Error())
		return
	}
	gerr = TransFunc.Init()
	if gerr != nil {
		RejectMsg(w, tr, defs.TRN_FORMAT_ERR, err.Error())
		return
	}
	gerr = TransFunc.DoTrans(tr)
	if gerr != nil {
		myLogger.Error("交易处理失败", gerr)
		RejectMsg(w, tr, gerr.GetErrorCode(), gerr.GetErrorString())
		return
	}
	myLogger.Info("应答处理完成")

	tr.MsgBody.Tran_cd = tr.MsgBody.Tran_cd[:3] + "2"
	tr.MsgBody.Resp_cd = "00"
	tr.MsgBody.Resp_msg = "SUCCESS"
	rep := tr.ToString()
	myLogger.Infof("应答报文：[%s]", rep)

	w.Write([]byte(rep))
	return
}

func RejectMsg(w http.ResponseWriter, msg *trans.TransMessage, resp_cd, resp_msg string) {
	if len(msg.MsgBody.Tran_cd) > 3 {
		msg.MsgBody.Tran_cd = msg.MsgBody.Tran_cd[:3] + "2"
	}
	myLogger.Debugf("resp_cd:%s, resp_msg:%s", resp_cd, resp_msg)
	myLogger.Debug(msg.MsgBody.Tran_cd)

	msg.MsgBody.Resp_cd = resp_cd
	msg.MsgBody.Resp_msg = resp_msg
	msg.MsgBody.Pos_sign = ""
	msg.MsgBody.Sign_img = ""
	msgbody, err := json.Marshal(msg.MsgBody)
	if err != nil {
		myLogger.Error("生成应答报文失败", err)
		w.WriteHeader(500)
		return
	}
	msg.Msg_body = string(msgbody)
	//msg.Signature = trans.Md5Base64(msgbody)
	res, err := json.Marshal(msg)
	if err != nil {
		myLogger.Error("生成应答报文失败", err)
		w.WriteHeader(500)
		return
	}
	myLogger.Debug("本地拒绝成功: ", string(res))
	w.Write(res)
}
