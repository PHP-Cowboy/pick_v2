package request

type CainiaoCloudprintIsvResourcesGetRequest struct {
	/*
	   isv资源类型，分为：TEMPLATE（表示模板），PRINT_ITEM（打印项），CUSTOM_AREA（预设自定义区）     */
	IsvResourceType *string `json:"isv_resource_type" required:"true" `
}

func (s *CainiaoCloudprintIsvResourcesGetRequest) SetIsvResourceType(v string) *CainiaoCloudprintIsvResourcesGetRequest {
	s.IsvResourceType = &v
	return s
}

func (req *CainiaoCloudprintIsvResourcesGetRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.IsvResourceType != nil {
		paramMap["isv_resource_type"] = *req.IsvResourceType
	}
	return paramMap
}

func (req *CainiaoCloudprintIsvResourcesGetRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}
