package domain

type CainiaoCloudprintCmdprintRenderRenderDocument struct {
	/*
	   包含的区域列表。对于有自定义区的文档，content会包含两条，即第一条是标准模板区域内容、第二条是自定义区域内容     */
	Contents *[]CainiaoCloudprintCmdprintRenderRenderContent `json:"contents,omitempty" `
}

func (s *CainiaoCloudprintCmdprintRenderRenderDocument) SetContents(v []CainiaoCloudprintCmdprintRenderRenderContent) *CainiaoCloudprintCmdprintRenderRenderDocument {
	s.Contents = &v
	return s
}
