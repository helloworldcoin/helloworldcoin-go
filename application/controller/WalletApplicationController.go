package controller

/*
 @author king 409060350@qq.com
*/
import (
	"helloworld-blockchain-go/application/service"
	"helloworld-blockchain-go/application/vo"
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/netcore"
	"helloworld-blockchain-go/util/StringUtil"
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
	account := AccountUtil.RandomAccount()
	accountVo := vo.AccountVo{account.PrivateKey, account.PublicKey, account.PublicKeyHash, account.Address}
	var response vo.CreateAccountResponse
	response.Account = &accountVo

	SuccessHttpResponse(rw, "", response)
}
func (w *WalletApplicationController) CreateAndSaveAccount(rw http.ResponseWriter, req *http.Request) {
	account := w.blockchainNetCore.GetBlockchainCore().GetWallet().CreateAndSaveAccount()
	accountVo := vo.AccountVo{account.PrivateKey, account.PublicKey, account.PublicKeyHash, account.Address}
	var response vo.CreateAndSaveAccountResponse
	response.Account = &accountVo

	SuccessHttpResponse(rw, "", response)
}

func (w *WalletApplicationController) SaveAccount(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.SaveAccountRequest{}).(*vo.SaveAccountRequest)

	privateKey := request.PrivateKey
	if StringUtil.IsNullOrEmpty(privateKey) {
		FailedHttpResponse(rw, "账户私钥不能为空。")
		return
	}
	account := AccountUtil.AccountFromPrivateKey(privateKey)
	w.blockchainNetCore.GetBlockchainCore().GetWallet().SaveAccount(account)
	var response vo.SaveAccountResponse
	response.AddAccountSuccess = true

	SuccessHttpResponse(rw, "", response)
}

func (w *WalletApplicationController) DeleteAccount(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.DeleteAccountRequest{}).(*vo.DeleteAccountRequest)

	address := request.Address
	if StringUtil.IsNullOrEmpty(address) {
		FailedHttpResponse(rw, "请填写需要删除的地址。")
		return
	}
	w.blockchainNetCore.GetBlockchainCore().GetWallet().DeleteAccountByAddress(address)
	var response vo.DeleteAccountResponse
	response.DeleteAccountSuccess = true

	SuccessHttpResponse(rw, "", response)
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

	SuccessHttpResponse(rw, "", response)
}
func (w *WalletApplicationController) AutoBuildTransaction(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, model.AutoBuildTransactionRequest{}).(*model.AutoBuildTransactionRequest)

	response := w.blockchainNetCore.GetBlockchainCore().AutoBuildTransaction(request)

	if response.BuildTransactionSuccess {
		SuccessHttpResponse(rw, "构建交易成功", response)
		return
	} else {
		FailedHttpResponse(rw, response.Message)
		return
	}
}

func (w *WalletApplicationController) SubmitTransactionToBlockchainNetwork(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.SubmitTransactionToBlockchainNetworkRequest{}).(*vo.SubmitTransactionToBlockchainNetworkRequest)

	response := w.walletApplicationService.SubmitTransactionToBlockchainNetwork(request)

	SuccessHttpResponse(rw, "", response)
}
