package WalletApplicationApi

//生成账户(公钥、私钥、地址)
const CREATE_ACCOUNT = "/Api/WalletApplication/CreateAccount"

//生成账户(公钥、私钥、地址)并保存
const CREATE_AND_SAVE_ACCOUNT = "/Api/WalletApplication/CreateAndSaveAccount"

//新增账户
const SAVE_ACCOUNT = "/Api/WalletApplication/SaveAccount"

//删除账户
const DELETE_ACCOUNT = "/Api/WalletApplication/DeleteAccount"

//查询所有的账户
const QUERY_ALL_ACCOUNTS = "/Api/WalletApplication/QueryAllAccounts"

//构建交易
const AUTO_BUILD_TRANSACTION = "/Api/WalletApplication/AutoBuildTransaction"

//提交交易到区块链网络
const SUBMIT_TRANSACTION_TO_BLOCKCHIAIN_NEWWORK = "/Api/WalletApplication/SubmitTransactionToBlockchainNetwork"
