package handle

import (
	"net/http"
	"upEletrcSign/iServer"
	"upEletrcSign/trans"
	"github.com/gogap/errors"
)

func DoHandle(w http.ResponseWriter, r *http.Request) {

	iServer.IDoFunct(w, r, DoSvr)
}

func DoSvr(msg *trans.TransMessage) (iServer.IDoAppTrans, error) {

	switch msg.MsgBody.Tran_cd {
	case "8261":
		return &iServer.T8262{}, nil
	default:
		return nil, errors.New("不识别的交易码: " + msg.MsgBody.Tran_cd)
	}
}
