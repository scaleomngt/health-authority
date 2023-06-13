package service

import (
	"errors"
	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"health-info-server/config"
	"health-info-server/gintool"
	"health-info-server/model"
	"health-info-server/utils"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var once sync.Once
var Accounts map[string]model.Account
var Index uint32
var lock sync.RWMutex

var PrivateKey = config.Config.GetString("PrivateKey")
var ViewKey = config.Config.GetString("ViewKey")
var ApiUrl = config.Config.GetString("ApiUrl")
var Contract = config.Config.GetString("Contract")

const Query = "https://vm.aleo.org/api"
const Broadcast = "https://vm.aleo.org/api/testnet3/transaction/broadcast"
const Prefix = "at"
const RPrefix = "health-server"
const FEE = "100000"

// CreateAccount 创建账户
func CreateAccount(c *gin.Context) {
	once.Do(func() {
		Accounts = make(map[string]model.Account, 0)
		Index = 0
	})

	cmd := "snarkos"
	args := []string{
		"account",
		"new"}
	result, err := utils.ExecCmdWithTimeout(60, cmd, args...)
	if err != nil {
		log.Println("err:", err, result)
		log.Println("result:", result)
		gintool.ResultFail(c, err.Error())
		return
	}
	log.Println("result:", result)
	keys := strings.Split(strings.TrimSpace(result), "\n")
	privateKey := strings.TrimSpace(strings.Split(keys[0], "Private Key")[1])
	viewKey := strings.TrimSpace(strings.Split(keys[1], "View Key")[1])
	address := strings.TrimSpace(strings.Split(keys[2], "Address")[1])
	Index = Index + 1
	account := model.Account{Index: Index, PrivateKey: privateKey, ViewKey: viewKey, Address: address}
	Accounts[account.Address] = account
	log.Println("Accounts:", Accounts)
	gintool.ResultOk(c, account)
}

// GetAccounts 获取GetAccounts
func GetAccounts(c *gin.Context) {
	if Accounts == nil {
		gintool.ResultOk(c, make([]model.Account, 0))
		return
	}

	gintool.ResultOk(c, Accounts)
}

// InitRedisId setRedis的id值
func InitRedisId(c *gin.Context) {
	id := c.Param("id")
	log.Println("id:", id)
	err := utils.SetId(id)
	if err != nil {
		log.Println("Error:", err)
		gintool.ResultFail(c, err.Error())
		return
	}
	gintool.ResultOkMsg(c, id, "success")
}

// GetData 获取数据
func GetData(c *gin.Context) {
	info := new(model.Transaction)
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println("Error:", err)
		gintool.ResultFail(c, err.Error())
		return
	}
	if info.Address == "" {
		gintool.ResultFail(c, "Address is not null")
		return
	}
	data, err := utils.GetAllHash(RPrefix + "-" + info.Address)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	gintool.ResultOkMsg(c, data, "success")
}

func CalcData(c *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	info := new(model.Transaction)
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println("Error1:", err)
		gintool.ResultFail(c, err.Error())
		return
	}

	log.Println("info: ", info)
	id := info.Id

	if info.Id == "" || info.Address == "" {
		gintool.ResultFail(c, "Id or Address is not null")
		return
	}

	type Data struct {
		Ti         string `json:"ti"`
		Ciphertext string `json:"ciphertext"`
		Classify   string `json:"classify"`
		Result     string `json:"result"`
		Owner      string `json:"owner"`
		Id         string `json:"id"`
	}

	hash, err := utils.GetOneHash(RPrefix+"-"+info.Address, info.Id)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}
	if hash != "" {
		d := &Data{}
		d.Result = hash
		d.Id = info.Id
		log.Println("hash: ", hash)
		gintool.ResultOkMsg(c, d, "success")
		return
	}

	time.Sleep(time.Duration(1000*10) * time.Millisecond)

	// 2. 获取步骤1中的output->value数据
	ciphertext, err := GetExecOutputValue(id)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	// 获取计算健康数据的入参
	value, r, err := DecryptCiphertext(ciphertext)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	classify, err := strconv.Atoi(r)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	// 获取最新record数据
	record, err := GetLatestFeeRecord()
	if err != nil {
		return
	}

	// 3.执行合约计算健康数据
	id, err = CalcHealthData(record, value, uint32(classify))
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	time.Sleep(time.Duration(1000*18) * time.Millisecond)

	// 获取计算结果
	ciphertext, err = GetExecOutputValue(id)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	// 解析结果
	value, _, err = DecryptCiphertext(ciphertext)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	data := &Data{}

	records := strings.Split(value, "\n")

	for i := 0; i < len(records); i++ {
		if strings.HasPrefix(strings.TrimSpace(records[i]), "owner:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], ".private,", "", -1))
			data.Owner = temp
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "classify:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "u32.private,", "", -1))
			data.Classify = temp
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "result:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "u32.private,", "", -1))
			data.Result = temp
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "id:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "field.private,", "", -1))
			//data.Id = temp
			log.Printf("id: %v", temp)
			data.Id = info.Id
		}
	}

	data.Ti = id
	data.Ciphertext = ciphertext

	err = utils.SetHash(RPrefix+"-"+info.Address, info.Id, data.Result)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	gintool.ResultOkMsg(c, data, "success")
}

// SubmitData 提交数据
func SubmitData(c *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	info := new(model.HealInfo)
	if err := c.ShouldBindJSON(&info); err != nil {
		log.Println("Error:", err)
		gintool.ResultFail(c, err.Error())
		return
	}
	log.Println("info: ", info)

	// 获取最新record数据
	record, err := GetLatestFeeRecord()
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	// 1.执行合约提交健康数据
	id, err := SubmitHealthData(record, info)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	err = utils.SetHash(RPrefix+"-"+info.Addr, id, "")
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	time.Sleep(time.Duration(1000*10) * time.Millisecond)

	r := &model.Transaction{
		Id:       id,
		Classify: info.Classify,
		Address:  info.Addr,
		Result:   "",
	}
	gintool.ResultOk(c, r)
}

func CalcHealthData(record, value string, classify uint32) (string, error) {
	id := ""
	method := ""
	if classify == 1 {
		method = "measure_bp"
	} else if classify == 2 {
		method = "measure_hr"
	} else if classify == 3 {
		method = "measure_pbg"
	} else {
		log.Println("没有此类别选项 classify:", classify)
		return id, errors.New("没有此类别选项")
	}

	args := []string{
		"developer",
		"execute",
		Contract,
		method,
		value,
		"--private-key",
		PrivateKey,
		"--query",
		Query,
		"--broadcast",
		Broadcast,
		"--fee",
		FEE,
		"--record",
		record}

	log.Println("args: ", args)

	cmd := "snarkos"
	result, err := utils.ExecCmdWithTimeout(60*5, cmd, args...)
	if err != nil {
		log.Println("执行合约计算健康数据 err:", err)
		log.Println("执行合约计算健康数据 result:", result)
		return id, err
	}
	log.Println("执行合约计算健康数据 result:", result)

	split := strings.Split(strings.TrimSpace(result), "\n")
	id = split[len(split)-1]
	if !strings.HasPrefix(id, Prefix) {
		log.Println("执行报错, 获取的数据有误, id: ", id)
		return "", errors.New(id)
	}

	err = utils.SetId(id)
	if err != nil {
		log.Println("utils.SetId err:", err)
		return "", err
	}
	return id, nil
}

func SubmitHealthData(record string, info *model.HealInfo) (string, error) {
	id := ""

	params := `{classify: {{classify}}u32, sbp: {{sbp}}u32, dbp: {{dbp}}u32, hr: {{hr}}u32, pbg: {{pbg}}u32, addr: {{addr}}}`
	params = strings.Replace(params, "{{classify}}", strconv.Itoa(int(info.Classify)), -1)
	params = strings.Replace(params, "{{sbp}}", strconv.Itoa(int(info.Sbp)), -1)
	params = strings.Replace(params, "{{dbp}}", strconv.Itoa(int(info.Dbp)), -1)
	params = strings.Replace(params, "{{hr}}", strconv.Itoa(int(info.Hr)), -1)
	params = strings.Replace(params, "{{pbg}}", strconv.Itoa(int(info.Pbg)), -1)
	params = strings.Replace(params, "{{addr}}", info.Addr, -1)
	args := []string{
		"developer",
		"execute",
		Contract,
		"submit",
		params,
		"--private-key",
		PrivateKey,
		"--query",
		Query,
		"--broadcast",
		Broadcast,
		"--fee",
		FEE,
		"--record",
		record}
	log.Println("执行合约提交健康数据 args:", args)

	cmd := "snarkos"
	result, err := utils.ExecCmdWithTimeout(60*5, cmd, args...)
	if err != nil {
		log.Println("执行合约提交健康数据 err:", err)
		log.Println("执行合约提交健康数据 result:", result)
		return id, err
	}
	log.Println("执行合约提交健康数据 result:", result)

	split := strings.Split(strings.TrimSpace(result), "\n")
	id = split[len(split)-1]
	if !strings.HasPrefix(id, Prefix) {
		log.Println("执行报错, 获取的数据有误, id: ", id)
		return "", errors.New(id)
	}

	log.Println("执行合约提交健康数据 id: ", id)

	err = utils.SetId(id)
	if err != nil {
		log.Println("utils.SetId err:", err)
		return id, err
	}

	return id, nil
}

// GetLatestFeeRecord 获取fee transition outputs 0 value
func GetLatestFeeRecord() (string, error) {
	id, err := utils.GetId()
	if err != nil {
		return "", err
	}

	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(ApiUrl + id)
	if err != nil {
		log.Println("获取最新record数据 error: ", err)
		return "", err
	}

	ciphertext, err := jsonparser.GetString(resp.Body(), "fee", "transition", "outputs", "[0]", "value")
	if err != nil {
		log.Println("获取最新record数据 value error: ", err)
		return "", err
	}

	record, _, err := DecryptCiphertext(ciphertext)
	if err != nil {
		log.Println("获取最新record数据 进行解密 error: ", err)
		return "", err
	}
	log.Println("获取最新record数据: ", record)
	return record, nil
}

func DecryptCiphertext(ciphertext string) (string, string, error) {
	cmd := "snarkos"

	args := []string{
		"developer",
		"decrypt",
		"--ciphertext",
		ciphertext,
		"--view-key",
		ViewKey}
	log.Println("args: ", args)
	record, err := utils.ExecCmdWithTimeout(60, cmd, args...)
	if err != nil {
		log.Println("DecryptCiphertext err:", err, record)
		log.Println("result:", record)
		return "", "", err
	}
	log.Println("DecryptCiphertext record:", strings.TrimSpace(record))

	records := strings.Split(strings.TrimSpace(record), "\n")

	classify := ""
	for i := 0; i < len(records); i++ {
		if strings.HasPrefix(strings.TrimSpace(records[i]), "classify:") {
			classify = strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "u32.private,", "", -1))
		}
	}
	return strings.TrimSpace(record), classify, nil
}

// GetExecOutputValue 获取 execution transitions value
func GetExecOutputValue(id string) (string, error) {
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		Get(ApiUrl + id)
	if err != nil {
		log.Println("发送http请求, 获取output->value数据 error: ", err)
		return "", err
	}
	cipherText, err := jsonparser.GetString(resp.Body(), "execution", "transitions", "[0]", "outputs", "[0]", "value")
	if err != nil {
		log.Println("获取json中 output->value数据 error: ", err)
		return "", err
	}
	log.Println("获取output->value数据 cipherText: ", cipherText)
	return cipherText, nil
}

func Test(c *gin.Context) {

	id := c.Param("id")
	ciphertext, err := GetExecOutputValue(id)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	value, _, err := DecryptCiphertext(ciphertext)
	if err != nil {
		gintool.ResultFail(c, err.Error())
		return
	}

	type Data struct {
		Ti         string `json:"ti"`
		Ciphertext string `json:"ciphertext"`
		Classify   string `json:"classify"`
		Result     string `json:"result"`
		Owner      string `json:"owner"`
		Id         string `json:"id"`
	}

	data := &Data{}

	records := strings.Split(value, "\n")

	for i := 0; i < len(records); i++ {
		if strings.HasPrefix(strings.TrimSpace(records[i]), "owner:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], ".private,", "", -1))
			data.Owner = temp
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "classify:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "u32.private,", "", -1))
			data.Classify = temp
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "result:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "u32.private,", "", -1))
			data.Result = temp
		} else if strings.HasPrefix(strings.TrimSpace(records[i]), "id:") {
			temp := strings.TrimSpace(strings.Replace(strings.Split(strings.TrimSpace(records[i]), ":")[1], "field.private,", "", -1))
			data.Id = temp
		}
	}

	data.Ti = id
	data.Ciphertext = ciphertext
	gintool.ResultOkMsg(c, data, "success")
}
