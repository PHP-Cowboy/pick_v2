package request

type CainiaoCloudprintSingleCustomareaGetRequest struct {
	/*
	   这是商家用户id     */
	SellerId *int64 `json:"seller_id,omitempty" required:"false" `
}

func (s *CainiaoCloudprintSingleCustomareaGetRequest) SetSellerId(v int64) *CainiaoCloudprintSingleCustomareaGetRequest {
	s.SellerId = &v
	return s
}

func (req *CainiaoCloudprintSingleCustomareaGetRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.SellerId != nil {
		paramMap["seller_id"] = *req.SellerId
	}
	return paramMap
}

func (req *CainiaoCloudprintSingleCustomareaGetRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}
