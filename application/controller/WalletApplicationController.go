package controller

import (
	"helloworld-blockchain-go/application/service"
	"helloworld-blockchain-go/application/vo/account"
	"helloworld-blockchain-go/application/vo/transaction"
	"helloworld-blockchain-go/core/Model/ModelWallet"
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
	accountVo := account.AccountVo{accountTemp.PrivateKey, accountTemp.PublicKey, accountTemp.PublicKeyHash, accountTemp.Address}
	var response account.CreateAccountResponse
	response.Account = &accountVo

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (w *WalletApplicationController) CreateAndSaveAccount(rw http.ResponseWriter, req *http.Request) {
	accountTemp := w.blockchainNetCore.GetBlockchainCore().GetWallet().CreateAndSaveAccount()
	accountVo := account.AccountVo{accountTemp.PrivateKey, accountTemp.PublicKey, accountTemp.PublicKeyHash, accountTemp.Address}
	var response account.CreateAndSaveAccountResponse
	response.Account = &accountVo

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (w *WalletApplicationController) SaveAccount(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), account.SaveAccountRequest{}).(*account.SaveAccountRequest)

	privateKey := request.PrivateKey
	if StringUtil.IsNullOrEmpty(privateKey) {
		//return Response.createFailResponse("账户私钥不能为空。");
	}
	accountTemp := AccountUtil.AccountFromPrivateKey(privateKey)
	w.blockchainNetCore.GetBlockchainCore().GetWallet().SaveAccount(accountTemp)
	var response account.SaveAccountResponse
	response.AddAccountSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (w *WalletApplicationController) DeleteAccount(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), account.DeleteAccountRequest{}).(*account.DeleteAccountRequest)

	address := request.Address
	if StringUtil.IsNullOrEmpty(address) {
		//return Response.createFailResponse("请填写需要删除的地址");
	}
	w.blockchainNetCore.GetBlockchainCore().GetWallet().DeleteAccountByAddress(address)
	var response account.DeleteAccountResponse
	response.DeleteAccountSuccess = true

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}

func (w *WalletApplicationController) QueryAllAccounts(rw http.ResponseWriter, req *http.Request) {
	wallet := w.blockchainNetCore.GetBlockchainCore().GetWallet()
	allAccounts := wallet.GetAllAccounts()

	var accountVos []*account.AccountVo2
	if allAccounts != nil {
		for _, accountTemp := range allAccounts {
			var accountVo account.AccountVo2
			accountVo.Address = accountTemp.Address
			accountVo.PrivateKey = accountTemp.PrivateKey
			accountVo.Value = wallet.GetBalanceByAddress(accountTemp.Address)
			accountVos = append(accountVos, &accountVo)
		}
	}

	var balance uint64
	for _, accountVo := range accountVos {
		balance = balance + accountVo.Value
	}

	var response account.QueryAllAccountsResponse
	response.Accounts = accountVos
	response.Balance = balance

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
func (w *WalletApplicationController) AutoBuildTransaction(rw http.ResponseWriter, req *http.Request) {
	result, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(result), ModelWallet.AutoBuildTransactionRequest{}).(*ModelWallet.AutoBuildTransactionRequest)

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
	request := JsonUtil.ToObject(string(result), transaction.SubmitTransactionToBlockchainNetworkRequest{}).(*transaction.SubmitTransactionToBlockchainNetworkRequest)

	response := w.walletApplicationService.SubmitTransactionToBlockchainNetwork(request)

	s := CreateSuccessResponse("", response)
	rw.Header().Set("content-type", "text/json")
	io.WriteString(rw, s)
}
