package controller

import (
	"helloworld-blockchain-go/application/service"
	"helloworld-blockchain-go/application/vo"
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/util/JsonUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"io"
	"io/ioutil"
	"net/http"
)

type WalletApplicationController struct {
	blockchainNetCore        *netcore.BlockchainNetCore
	walletApplicationService *service.WalletApplicationService
}

func NewWalletApplicationController(blockchainNetCore *netcore.BlockchainNetCore, walletApplicationService *service.WalletApplicationService) *WalletApplicationController {
	var b WalletApplicationController
	b.blockchainNetCore = blockchainNetCore
	b.walletApplicationService = walletApplicationService
	return &b
}

func (w *WalletApplicationController) CreateAccount(rw http.ResponseWriter, req *http.Request) {
	accountTemp := AccountUtil.RandomAccount()
	accountVo := vo.AccountVo{accountTemp.PrivateKey, accountTemp.PublicKey, accountTemp.PublicKeyHash, accountTemp.Address}
	var response vo.CreateAccountResponse
	response.Account = &accountVo

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (w *WalletApplicationController) CreateAndSaveAccount(rw http.ResponseWriter, req *http.Request) {
	account := w.blockchainNetCore.GetBlockchainCore().GetWallet().CreateAndSaveAccount()
	accountVo := vo.AccountVo{account.PrivateKey, account.PublicKey, account.PublicKeyHash, account.Address}
	var response vo.CreateAndSaveAccountResponse
	response.Account = &accountVo

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (w *WalletApplicationController) SaveAccount(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.SaveAccountRequest{}).(*vo.SaveAccountRequest)

	privateKey := request.PrivateKey
	if StringUtil.IsNullOrEmpty(privateKey) {
		//return Response.createFailResponse("账户私钥不能为空。");
	}
	accountTemp := AccountUtil.AccountFromPrivateKey(privateKey)
	w.blockchainNetCore.GetBlockchainCore().GetWallet().SaveAccount(accountTemp)
	var response vo.SaveAccountResponse
	response.AddAccountSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (w *WalletApplicationController) DeleteAccount(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.DeleteAccountRequest{}).(*vo.DeleteAccountRequest)

	address := request.Address
	if StringUtil.IsNullOrEmpty(address) {
		//return Response.createFailResponse("请填写需要删除的地址");
	}
	w.blockchainNetCore.GetBlockchainCore().GetWallet().DeleteAccountByAddress(address)
	var response vo.DeleteAccountResponse
	response.DeleteAccountSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (w *WalletApplicationController) QueryAllAccounts(rw http.ResponseWriter, req *http.Request) {
	wallet := w.blockchainNetCore.GetBlockchainCore().GetWallet()
	allAccounts := wallet.GetAllAccounts()

	var accountVos []*vo.AccountVo2
	if allAccounts != nil {
		for _, account := range allAccounts {
			var accountVo vo.AccountVo2
			accountVo.Address = account.Address
			accountVo.PrivateKey = account.PrivateKey
			accountVo.Value = wallet.GetBalanceByAddress(account.Address)
			accountVos = append(accountVos, &accountVo)
		}
	}

	var balance uint64
	for _, accountVo := range accountVos {
		balance = balance + accountVo.Value
	}

	var response vo.QueryAllAccountsResponse
	response.Accounts = accountVos
	response.Balance = balance

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (w *WalletApplicationController) AutoBuildTransaction(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), model.AutoBuildTransactionRequest{}).(*model.AutoBuildTransactionRequest)

	response := w.blockchainNetCore.GetBlockchainCore().AutoBuildTransaction(request)
	/*	if autoBuildTransactionResponse.IsBuildTransactionSuccess() {
			return Response.createSuccessResponse("构建交易成功", autoBuildTransactionResponse)
		} else {
			return Response.createFailResponse(autoBuildTransactionResponse.getMessage())
		}*/

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (w *WalletApplicationController) SubmitTransactionToBlockchainNetwork(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), vo.SubmitTransactionToBlockchainNetworkRequest{}).(*vo.SubmitTransactionToBlockchainNetworkRequest)

	response := w.walletApplicationService.SubmitTransactionToBlockchainNetwork(request)

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
