package vo

/*
 @author x.king xdotking@gmail.com
*/

type ActiveAutoSearchBlockRequest struct {
}

type ActiveAutoSearchBlockResponse struct {
}

type DeactiveAutoSearchBlockRequest struct {
}

type DeactiveAutoSearchBlockResponse struct {
}

type IsAutoSearchBlockRequest struct {
}

type IsAutoSearchBlockResponse struct {
	AutoSearchBlock bool `json:"autoSearchBlock"`
}
