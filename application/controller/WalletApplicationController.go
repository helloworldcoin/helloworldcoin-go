package controller

/*
 @author x.king xdotking@gmail.com
*/
import (
	"helloworldcoin-go/application/service"
	"helloworldcoin-go/application/vo"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/netcore"
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

	success(rw, response)
}
func (w *WalletApplicationController) CreateAndSaveAccount(rw http.ResponseWriter, req *http.Request) {
	account := w.blockchainNetCore.GetBlockchainCore().GetWallet().CreateAndSaveAccount()
	accountVo := vo.AccountVo{account.PrivateKey, account.PublicKey, account.PublicKeyHash, account.Address}
	var response vo.CreateAndSaveAccountResponse
	response.Account = &accountVo

	success(rw, response)
}

func (w *WalletApplicationController) SaveAccount(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.SaveAccountRequest{}).(*vo.SaveAccountRequest)

	privateKey := request.PrivateKey
	account := AccountUtil.AccountFromPrivateKey(privateKey)
	w.blockchainNetCore.GetBlockchainCore().GetWallet().SaveAccount(account)
	var response vo.SaveAccountResponse

	success(rw, response)
}

func (w *WalletApplicationController) DeleteAccount(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.DeleteAccountRequest{}).(*vo.DeleteAccountRequest)

	address := request.Address
	w.blockchainNetCore.GetBlockchainCore().GetWallet().DeleteAccountByAddress(address)
	var response vo.DeleteAccountResponse

	success(rw, response)
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

	success(rw, response)
}
func (w *WalletApplicationController) AutomaticBuildTransaction(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.AutomaticBuildTransactionRequest{}).(*vo.AutomaticBuildTransactionRequest)

	response := w.walletApplicationService.AutomaticBuildTransaction(request)

	if response.BuildTransactionSuccess {
		success(rw, response)
		return
	} else {
		serviceUnavailable(rw)
		return
	}
}
func (w *WalletApplicationController) SubmitTransactionToBlockchainNetwork(rw http.ResponseWriter, req *http.Request) {
	request := GetRequest(req, vo.SubmitTransactionToBlockchainNetworkRequest{}).(*vo.SubmitTransactionToBlockchainNetworkRequest)

	response := w.walletApplicationService.SubmitTransactionToBlockchainNetwork(request)

	success(rw, response)
}
