package vo

/*
 @author king 409060350@qq.com
*/

type ActiveAutoSearchBlockRequest struct {
}

type ActiveAutoSearchBlockResponse struct {
	ActiveAutoSearchBlockSuccess bool `json:"activeAutoSearchBlockSuccess"`
}

type DeactiveAutoSearchBlockRequest struct {
}

type DeactiveAutoSearchBlockResponse struct {
	DeactiveAutoSearchBlockSuccess bool `json:"deactiveAutoSearchBlockSuccess"`
}

type IsAutoSearchBlockRequest struct {
}

type IsAutoSearchBlockResponse struct {
	AutoSearchBlock bool `json:"autoSearchBlock"`
}
