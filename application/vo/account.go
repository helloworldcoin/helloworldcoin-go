package vo

/*
 @author x.king xdotking@gmail.com
*/

type AccountVo struct {
	PrivateKey    string `json:"privateKey"`
	PublicKey     string `json:"publicKey"`
	PublicKeyHash string `json:"publicKeyHash"`
	Address       string `json:"address"`
}
type AccountVo2 struct {
	PrivateKey string `json:"privateKey"`
	Address    string `json:"address"`
	Value      uint64 `json:"value"`
}
type CreateAccountRequest struct {
}
type CreateAccountResponse struct {
	Account *AccountVo `json:"account"`
}
type CreateAndSaveAccountRequest struct {
}
type CreateAndSaveAccountResponse struct {
	Account *AccountVo `json:"account"`
}
type DeleteAccountRequest struct {
	Address string `json:"address"`
}
type DeleteAccountResponse struct {
}
type QueryAllAccountsRequest struct {
}
type QueryAllAccountsResponse struct {
	Balance  uint64        `json:"balance"`
	Accounts []*AccountVo2 `json:"accounts"`
}
type SaveAccountRequest struct {
	PrivateKey string `json:"privateKey"`
}
type SaveAccountResponse struct {
}
