package domain

type CainiaoCloudprintCmdprintRenderRenderConfig struct {
	/*
	   打印方向：normal-正常 reverse-翻转(旋转180°) defalutValue:normal    */
	Orientation *string `json:"orientation,omitempty" `

	/*
	   下联logo defalutValue:true    */
	NeedBottomLogo *bool `json:"need_bottom_logo,omitempty" `

	/*
	   中间logo defalutValue:true    */
	NeedMiddleLogo *bool `json:"need_middle_logo,omitempty" `

	/*
	   上联logo defalutValue:true    */
	NeedTopLogo *bool `json:"need_top_logo,omitempty" `
}

func (s *CainiaoCloudprintCmdprintRenderRenderConfig) SetOrientation(v string) *CainiaoCloudprintCmdprintRenderRenderConfig {
	s.Orientation = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderRenderConfig) SetNeedBottomLogo(v bool) *CainiaoCloudprintCmdprintRenderRenderConfig {
	s.NeedBottomLogo = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderRenderConfig) SetNeedMiddleLogo(v bool) *CainiaoCloudprintCmdprintRenderRenderConfig {
	s.NeedMiddleLogo = &v
	return s
}
func (s *CainiaoCloudprintCmdprintRenderRenderConfig) SetNeedTopLogo(v bool) *CainiaoCloudprintCmdprintRenderRenderConfig {
	s.NeedTopLogo = &v
	return s
}
