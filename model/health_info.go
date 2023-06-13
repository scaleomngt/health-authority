package model

type HealInfo struct {
	Classify uint32 `json:"classify"` // 分类
	Sbp      uint32 `json:"sbp"`      // 舒张压
	Dbp      uint32 `json:"dbp"`      // 收缩压
	Hr       uint32 `json:"hr"`       // 心率
	Pbg      uint32 `json:"pbg"`      // 血糖
	Addr     string `json:"addr"`     // 地址
}

type Account struct {
	Index      uint32 `json:"index"`      // 序号
	PrivateKey string `json:"privateKey"` // Private Key
	ViewKey    string `json:"viewKey"`    // View Key
	Address    string `json:"address"`    // Address
}

type Health struct {
	Owner string `json:"owner"`
	Gates string `json:"gates"`
	Id    string `json:"id"`
	Nonce string `json:"nonce"`
}

type Transaction struct {
	Id       string `json:"id"`
	Address  string `json:"address"`
	Classify uint32 `json:"classify"`
	Result   string `json:"result"`
}
