package BlockTool

import (
	"helloworld-blockchain-go/core/Model"
	"testing"
)

func TestCalculateBlockHash(t *testing.T) {
	var Inputs0 []*Model.TransactionInput
	var Outputs0 []*Model.TransactionOutput
	Outputs0 = append(Outputs0, &Model.TransactionOutput{Value: 5000000000, OutputScript: &[]string{"01", "02", "00", "c80fe6d25b78f94c370e852226026f2d35833a9d", "03", "04"}})
	var t0 Model.Transaction
	t0.Inputs = Inputs0
	t0.Outputs = Outputs0

	var Inputs1 []*Model.TransactionInput
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "c5f03a4925331566e0b760fb3e317f7a37df42954a2851cae57497484a41440e", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"3045022100e8bb3846dfb317581c7518c62e4d0d332631f8796258c1af9dfb46926159f37d022058927e240d4e3d629838a25f49aa06f00c5b147372d0953fce778f2e56cd13b0",
		"00",
		"036642b7d2711330e00a383ff23f5efffa699bfab69157d2c45414f418f4b74f35"}})
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "e959b0604b0004436079b2e50d3ec08e04c0e2ba8ffa83a32ff98e711540c662", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"3045022100f54029fc7291d23c47462c92a5d1f4f5a6a4bf5cfdbf75671521a0658fd75ae2022059e8500750c199d3842507315f977a87564412cecd498e187c6f77a8b3a22067",
		"00",
		"02cba3c3281785973c4b6229784a5ad510286faf755a867e3215dbcc9ac95c1892"}})
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "8a438068122eff5a6ea5718f781f5481638352d9063d5b2b7c392242bb26e9d3", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"304402204318651fce3f47324c66e8c483714d2f6f2fd9b029b4c6f31b28f01b65e2035502203b1c6bcad1624a38be77e341078a1cf11bcc7e6d54afd09ccaae6202365d4870",
		"00",
		"020ff2680c2bd071666d74579a4568926bea504a2f145f70fda393c08a8d4dbf3e"}})
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "5360b4b5009ac051cd6868aebae09346680389db62f323cedeb6317c9f73bf5c", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"3045022100e53cd28257e239586972cf88247f63335c213b8c8f0a28566b94dc541e5b654c022010e87b9d05c27ece9d72b215e6d60795df6fb11e527dca811e564a2eb7d5f431",
		"00",
		"03a5fc19462ed85c37910c7b756397aa0e0078b95f58b2b035cd73d72ac9d56944"}})
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "b3b9fe544c376d9bd2f102bac1508b80e5aaec24d12c1814e9c1c36e909d3e0a", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"304402206e7a447019d6b62e21becc68cd1467c68551ea55701f5d3572bd5b4ff2b62987022063a9cf233bc010dca0f7b3849bc72a1cffb907b8ccc34d5e200049fa5b56da7d",
		"00",
		"03e50554a6260fbbae50a253fbe4324c362dc43d82737aaf6590f5dfc1ecf82c72"}})
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "23114d0bf66656ff04a59bb045f607432eccefeb23281196689c6025a555b3c7", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"3045022100ca5d0bbbd3d825f6e927f804c4e0f577448d6bd7961961adb957d8c144ea378202203e8cbe476b681f49fd9400390980af91860ac51ba5a3c846f139027e09922ee5",
		"00",
		"034d78a2480fc20d1d5766b1e2754ab25e4c9dd888033cbb756270baee2be22424"}})
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "7c4303fb7bafd09333530303f330af06ccd55cc3db8e90aac1518fbd65a87dec", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"304402200ef1c82381db3a828e347060c2d3ee1077bdae095b0a65468c08045d57a755f40220704c992b6999f2e2ac3d22157297ba32b24229814e6089cc8f602fe9f8f5143b",
		"00",
		"03018bbff04fac69b1b00d2c353a438a6627859961987bb1ec3bae397b548f4094"}})
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "c6312816df19956db4cd3fc15352144a6d368bd93e18131ca54608316cbb36d4", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"3045022100b9edd06bac8774420a9c49bd9f4ed9fc6e03b39ad5056dfc45e854046563f56d022024fb7b7867a253d34ccd75db3f7dc334d7334d87b3151a5c5ebf119ade12ad1d",
		"00",
		"03956cbc4d043929facdb5d4fd21735c97a0b1c12dc3f571179fcc60ab76fa05b2"}})
	Inputs1 = append(Inputs1, &Model.TransactionInput{UnspentTransactionOutput: &Model.TransactionOutput{TransactionHash: "6f563b71c2b5715cdbc1ddf1965ac872b98383b222c75ba05dac544372b9d0b6", TransactionOutputIndex: 1}, InputScript: &[]string{"00",
		"30450221008f99eca1b9939a4676dd6d53446acf5d48f5396cc9e88b748d2acca18497bc1402205935798e20c518101aab376703089a11a1857b164f3ba0df6d70018d2fda7084",
		"00",
		"02af31d8e6786161e4a9f7edaa192c5332a7b21953323888a7d17cef4247ced523"}})

	var Outputs1 []*Model.TransactionOutput
	Outputs1 = append(Outputs1, &Model.TransactionOutput{Value: 5000000000, OutputScript: &[]string{"01", "02", "00", "7217191657b39042d4ccd6f05a75a704f7fc66be", "03", "04"}})
	Outputs1 = append(Outputs1, &Model.TransactionOutput{Value: 15000000000, OutputScript: &[]string{"01", "02", "00", "80d0fc0a579329798a21cb009bb413d49b05ce79", "03", "04"}})
	Outputs1 = append(Outputs1, &Model.TransactionOutput{Value: 25000000000, OutputScript: &[]string{"01", "02", "00", "d6a5985a1b31f9036f8552c5a975cf9cd784b50b", "03", "04"}})

	var t1 Model.Transaction
	t1.Inputs = Inputs1
	t1.Outputs = Outputs1

	transactions := []*Model.Transaction{&t0, &t1}

	block := new(Model.Block)
	block.Timestamp = 1622928126312
	block.PreviousHash = "52ae3991310d52ff43e530b7c69c25abf7ce60667f77a3b19231da904f21648e"
	block.Nonce = "d63852382edbec3b79de31e58b34c8a703a2f3715efb7c5f2dd52c9ae27250b6"
	block.Transactions = transactions

	if "87bbd24191ebe280407698467e51ac4b2738841d14cb3f5fe3e8bd57201a9e01" != CalculateBlockMerkleTreeRoot(block) {
		t.Error("test failed")
	}

	if "4758aa769a5f0642cceb8fc7deb6a0102aa7faf28d2dff324b911b1a6650fbb4" != CalculateBlockHash(block) {
		t.Error("test failed")
	}
}
