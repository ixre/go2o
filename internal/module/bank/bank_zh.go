package bank

const (
	ZH_ICBC   = "ICBC"
	ZH_ABC    = "ABC"
	ZH_CCB    = "CCB"
	ZH_BOC    = "BOC"
	ZH_CIB    = "CIB"
	ZH_CITIC  = "CITIC"
	ZH_CMB    = "CMB"
	ZH_CMBC   = "CMBC"
	ZH_CEB    = "CEB"
	ZH_COMM   = "COMM"
	ZH_GDB    = "GDB"
	ZH_SPAB   = "SPAB"
	ZH_POSTGC = "POSTGC"
	ZH_SPDB   = "SPDB"
)

type BankItem struct {
	//标识
	ID string
	//名称
	Name string
	//签名/别名
	Sign string
}

type PaymentPlatform struct {
	ID   string
	Name string
	Sign string
	Bank []*BankItem
}

func newBankItem(id string, name string, sign string) *BankItem {
	return &BankItem{
		ID:   id,
		Name: name,
		Sign: sign,
	}
}

var (
	Alipay = &PaymentPlatform{
		ID:   "AliPay",
		Name: "支付宝",
		Sign: "alipay",
		Bank: []*BankItem{
			newBankItem(ZH_ICBC, "中国工商银行", "ICBCBTB"),
			newBankItem(ZH_ABC, "中国农业银行", "ABCBTB"),
			newBankItem(ZH_CCB, "中国建设银行", "CCBBTB"),
			newBankItem(ZH_BOC, "中国银行", "BOCB2C"),
			newBankItem(ZH_CIB, "兴业银行", "CIB"),
			newBankItem(ZH_CITIC, "中信银行", "CITIC"),
			newBankItem(ZH_CMB, "招商银行", "CMB"),
			newBankItem(ZH_CMBC, "中国民生银行", "CMBC"),
			newBankItem(ZH_CEB, "中国光大银行", "CEBBANK"),
			newBankItem(ZH_COMM, "交通银行", "COMM"),
			newBankItem(ZH_GDB, "广发银行", "GDB"),
			newBankItem(ZH_SPAB, "平安银行", "SPABANK"),
			newBankItem(ZH_POSTGC, "中国邮政储蓄银行", "POSTGC"),
			newBankItem(ZH_SPDB, "浦发银行", "SPDBB2B"),
		},
	}

	KuaiBill = &PaymentPlatform{
		ID:   "KuaiQian",
		Name: "快钱",
		Sign: "kuaiqian",
	}

	ChinaPay = &PaymentPlatform{
		ID:   "ChinaPay",
		Name: "银联商务",
		Sign: "chinapay",
	}

	Tenpay = &PaymentPlatform{
		ID:   "WeiXin",
		Name: "财付通",
		Sign: "tenpay",
		Bank: []*BankItem{
			newBankItem(ZH_ICBC, "中国工商银行", "ICBC"),
			newBankItem(ZH_ABC, "中国农业银行", "ABC"),
			newBankItem(ZH_CCB, "中国建设银行", "CCB"),
			newBankItem(ZH_BOC, "中国银行", "BOC"),
			newBankItem(ZH_CIB, "兴业银行", "CIB"),
			newBankItem(ZH_CITIC, "中信银行", "CITIC"),
			newBankItem(ZH_CMB, "招商银行", "CMB"),
			newBankItem(ZH_CMBC, "中国民生银行", "CMBC"),
			newBankItem(ZH_CEB, "中国光大银行", "CEB"),
			newBankItem(ZH_COMM, "交通银行", "COMM"),
			newBankItem(ZH_GDB, "广发银行", "GDB"),
			newBankItem(ZH_SPAB, "平安银行", "PAB"),
			newBankItem(ZH_POSTGC, "中国邮政储蓄银行", "POSTGC"),
			newBankItem(ZH_SPDB, "浦发银行", "SPDB"),
		},
	}
)
