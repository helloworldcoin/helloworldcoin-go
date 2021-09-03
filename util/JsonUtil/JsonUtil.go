package JsonUtil

import (
	"encoding/json"
	"helloworld-blockchain-go/application/vo"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore/po"
)

func ToString(object interface{}) string {
	jsonString, _ := json.Marshal(object)
	return string(jsonString)
}

func ToObject(jsonString string, object interface{}) interface{} {
	_0001, ok := object.(dto.GetBlockRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0001)
		return &_0001
	}
	_1001, ok := object.(dto.GetBlockResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1001)
		return &_1001
	}
	_0002, ok := object.(dto.PostBlockRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0002)
		return &_0002
	}
	_1002, ok := object.(dto.PostBlockResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1002)
		return &_1002
	}
	_0003, ok := object.(dto.GetUnconfirmedTransactionsRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0003)
		return &_0003
	}
	_1003, ok := object.(dto.GetUnconfirmedTransactionsResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1003)
		return &_1003
	}
	_0004, ok := object.(dto.GetBlockchainHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0004)
		return &_0004
	}
	_1004, ok := object.(dto.GetBlockchainHeightResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1004)
		return &_1004
	}
	_0005, ok := object.(dto.PostTransactionRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0005)
		return &_0005
	}
	_1005, ok := object.(dto.PostTransactionResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1005)
		return &_1005
	}
	_0006, ok := object.(dto.PostBlockchainHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0006)
		return &_0006
	}
	_1006, ok := object.(dto.PostBlockchainHeightResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1006)
		return &_1006
	}
	_0007, ok := object.(dto.PingRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0007)
		return &_0007
	}
	_1007, ok := object.(dto.PingResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1007)
		return &_1007
	}
	_0008, ok := object.(dto.GetNodesRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0008)
		return &_0008
	}
	_1008, ok := object.(dto.GetNodesResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1008)
		return &_1008
	}
	_0009, ok := object.(dto.BlockDto)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0009)
		return &_0009
	}
	_1009, ok := object.(dto.BlockDto)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1009)
		return &_1009
	}

	_0010, ok := object.(po.NodePo)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0010)
		return &_0010
	}

	_0100, ok := object.(vo.ActiveAutoSearchNodeRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0100)
		return &_0100
	}
	_0101, ok := object.(vo.AddNodeRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0101)
		return &_0101
	}
	_0102, ok := object.(vo.DeactiveAutoSearchNodeRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0102)
		return &_0102
	}
	_0103, ok := object.(vo.DeleteNodeRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0103)
		return &_0103
	}
	_0104, ok := object.(vo.IsAutoSearchNodeRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0104)
		return &_0104
	}
	_0105, ok := object.(vo.NodeVo)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0105)
		return &_0105
	}
	_0106, ok := object.(vo.QueryAllNodesRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0106)
		return &_0106
	}
	_0107, ok := object.(vo.QueryBlockchainHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0107)
		return &_0107
	}
	_0108, ok := object.(vo.UpdateNodeRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0108)
		return &_0108
	}

	_0150, ok := object.(vo.DeleteBlocksRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0150)
		return &_0150
	}
	_0151, ok := object.(vo.QueryBlockByBlockHashRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0151)
		return &_0151
	}
	_0152, ok := object.(vo.QueryBlockByBlockHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0152)
		return &_0152
	}
	_0153, ok := object.(vo.QueryTop10BlocksRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0153)
		return &_0153
	}
	_0200, ok := object.(vo.SetMaxBlockHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0200)
		return &_0200
	}

	_0250, ok := object.(vo.CreateAccountRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0250)
		return &_0250
	}
	_0251, ok := object.(vo.CreateAndSaveAccountRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0251)
		return &_0251
	}
	_0252, ok := object.(vo.DeleteAccountRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0252)
		return &_0252
	}
	_0253, ok := object.(vo.QueryAllAccountsRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0253)
		return &_0253
	}
	_0254, ok := object.(vo.SaveAccountRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0254)
		return &_0254
	}
	_0300, ok := object.(vo.AutoBuildTransactionRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0300)
		return &_0300
	}

	_0350, ok := object.(vo.QueryTransactionByTransactionHashRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0350)
		return &_0350
	}
	_0351, ok := object.(vo.QueryTransactionOutputByAddressRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0351)
		return &_0351
	}
	_0352, ok := object.(vo.QueryTransactionOutputByTransactionOutputIdRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0352)
		return &_0352
	}
	_0353, ok := object.(vo.QueryTransactionsByBlockHashTransactionHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0353)
		return &_0353
	}
	_0354, ok := object.(vo.QueryUnconfirmedTransactionByTransactionHashRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0354)
		return &_0354
	}
	_0355, ok := object.(vo.QueryUnconfirmedTransactionsRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0355)
		return &_0355
	}
	_0356, ok := object.(vo.SubmitTransactionToBlockchainNetworkRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0356)
		return &_0356
	}
	panic("JsonUtil.ToObject can not recognize object type")
}
