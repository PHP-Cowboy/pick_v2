package domain

type CainiaoWaybillIiGetPackageInfoDto struct {
	/*
	   包裹id，用于拆合单场景（只能传入数字、字母和下划线；批量请求时值不得重复，大小写敏感，即123A,123a 不可当做不同ID，否则存在一定可能取号失败）     */
	Id *string `json:"id,omitempty" `

	/*
	   商品信息,数量限制为100     */
	Items *[]CainiaoWaybillIiGetItem `json:"items,omitempty" `

	/*
	   体积, 单位 ml defalutValue:0    */
	Volume *int64 `json:"volume,omitempty" `

	/*
	   重量,单位 g defalutValue:0    */
	Weight *int64 `json:"weight,omitempty" `

	/*
	   子母件模式中的总包裹数/总件数，用于打印当前包裹处于总件数的位置比如5-2，可以表示总包裹数为5，当前为第2个包裹，只有快运公司需要传入，其他的可以不用传入     */
	TotalPackagesCount *int64 `json:"total_packages_count,omitempty" `

	/*
	   大件快运中的包装方式描述     */
	PackagingDescription *string `json:"packaging_description,omitempty" `

	/*
	   大件快运中的货品描述，比如服装，家具。 顺丰取号必须传此参数     */
	GoodsDescription *string `json:"goods_description,omitempty" `

	/*
	   包裹长，单位厘米     */
	Length *int64 `json:"length,omitempty" `

	/*
	   包裹宽，单位厘米     */
	Width *int64 `json:"width,omitempty" `

	/*
	   包裹高，单位厘米     */
	Height *int64 `json:"height,omitempty" `

	/*
	   物品价值，单位元     */
	GoodValue *string `json:"good_value,omitempty" `
}

func (s *CainiaoWaybillIiGetPackageInfoDto) SetId(v string) *CainiaoWaybillIiGetPackageInfoDto {
	s.Id = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetItems(v []CainiaoWaybillIiGetItem) *CainiaoWaybillIiGetPackageInfoDto {
	s.Items = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetVolume(v int64) *CainiaoWaybillIiGetPackageInfoDto {
	s.Volume = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetWeight(v int64) *CainiaoWaybillIiGetPackageInfoDto {
	s.Weight = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetTotalPackagesCount(v int64) *CainiaoWaybillIiGetPackageInfoDto {
	s.TotalPackagesCount = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetPackagingDescription(v string) *CainiaoWaybillIiGetPackageInfoDto {
	s.PackagingDescription = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetGoodsDescription(v string) *CainiaoWaybillIiGetPackageInfoDto {
	s.GoodsDescription = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetLength(v int64) *CainiaoWaybillIiGetPackageInfoDto {
	s.Length = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetWidth(v int64) *CainiaoWaybillIiGetPackageInfoDto {
	s.Width = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetHeight(v int64) *CainiaoWaybillIiGetPackageInfoDto {
	s.Height = &v
	return s
}
func (s *CainiaoWaybillIiGetPackageInfoDto) SetGoodValue(v string) *CainiaoWaybillIiGetPackageInfoDto {
	s.GoodValue = &v
	return s
}
