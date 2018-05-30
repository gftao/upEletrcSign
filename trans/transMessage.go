package trans

import "encoding/json"

type MessageHead struct {
	Encoding    string `json:"encoding"`
	Sign_method string `json:"sign_method"`
	Signature   string `json:"signature"`
	Version     string `json:"version"`
}

type TransMessage struct {
	MessageHead
	Msg_body string       `json:"msg_body"`
	MsgBody  *TransParams `json:"-"`
}

type TransParams struct {
	commonParams
	Pos_sign          string    `json:"pos_sign,omitempty"`
	UPos_sign         string    `json:"-"`
	Term_inf          *TermInfo `json:"term_inf,omitempty"`
	Orig_sys_order_id string    `json:"orig_sys_order_id,omitempty"`
	Orig_term_seq     string    `json:"orig_term_seq,omitempty"`
	//电子签名
	Sign_img string `json:"sign_img,omitempty"`
}

func (t *TransMessage) ToString() string {
	t.SetMsgBody()
	res, err := json.Marshal(t)
	if err != nil {
		return "{}"
	}
	//log.Println(string(res))
	return string(res)
}

func (t *TransMessage) SetMsgBody() {
	btMsgBody, err := json.Marshal(t.MsgBody)
	if err != nil {
		t.Msg_body = "{}"
		return
	}
	t.Msg_body = string(btMsgBody)
}

type commonParams struct {
	Tran_cd          string `json:"tran_cd,omitempty"`
	Prod_cd          string `json:"prod_cd,omitempty"`
	Biz_cd           string `json:"biz_cd,omitempty"`
	Term_seq         string `json:"term_seq,omitempty"`
	Term_batch       string `json:"term_batch,omitempty"`
	Mcht_cd          string `json:"mcht_cd,omitempty"`
	Mcht_nm          string `json:"mcht_nm,omitempty"`
	Term_id          string `json:"term_id,omitempty"`
	Term_flag        string `json:"term_flag,omitempty"`
	Tran_dt_tm       string `json:"tran_dt_tm,omitempty"`
	Order_id         string `json:"order_id,omitempty"`
	Order_timeout    string `json:"order_timeout,omitempty"`
	Sys_order_id     string `json:"sys_order_id,omitempty"`  //产品平台订单号
	Cld_order_id     string `json:"-"`                       //云前置自己的订单号
	Tran_order_id    string `json:"tran_order_id,omitempty"` //交易订单号，汇宜系统内唯一
	Acct_order_id    string `json:"-"`
	Order_desc       string `json:"order_desc,omitempty"`
	Req_reserved     string `json:"req_reserved,omitempty"`
	Resp_cd          string `json:"resp_cd,omitempty"`
	Resp_msg         string `json:"resp_msg,omitempty"`
	ActiveCode       string `json:"active_code,omitempty"`
	Tran_amt         string `json:"tran_amt,omitempty"`
	Curr_cd          string `json:"curr_cd,omitempty"`
	Pre_auth_id      string `json:"pre_auth_id,omitempty"`
	Sett_dt          string `json:"sett_dt"`
	Ins_id_cd        string `json:"ins_id_cd,omitempty"`
	Iss_ins_id_cd    string `json:"iss_ins_id_cd"`
	Trans_in_acct_no string `json:"trans_in_acct_no,omitempty"`
	Chn_ins_id_cd    string `json:"chn_ins_id_cd,omitempty"`
}

type TermInfo struct {
	Ip_addr    string `json:"ip_addr,omitempty"`
	Gps_addr   string `json:"gps_addr,omitempty"`
	Term_prod  string `json:"term_prod,omitempty"`
	Term_model string `json:"term_model,omitempty"`
	Brand_ksn  string `json:"brand_ksn,omitempty"`
	Brand_sn   string `json:"brand_sn,omitempty"`
	Term_tp    string `json:"term_tp,omitempty"`
	Term_sn    string `json:"term_sn,omitempty"`
	Term_rand  string `json:"term_rand,omitempty"`
	Term_enc   string `json:"term_enc,omitempty"`
	Term_ver   string `json:"term_ver,omitempty"`
}

type Spos struct {
	Pos []*Pos `json:"spos,omitempty"`
}
type Pos struct {
	Bold         string `json:"spos->bold,omitempty"`
	Content      string `json:"content,omitempty"`
	Content_type string `json:"content-type,omitempty"`
	Position     string `json:"position,omitempty"`
	Size         string `json:"size,omitempty"`
}
