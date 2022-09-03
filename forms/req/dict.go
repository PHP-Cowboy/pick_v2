package req

type CreateDictTypeForm struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type DictTypeListForm struct {
	Paging
	Code string `json:"code"`
	Name string `json:"name"`
}

type ChangeDictTypeForm struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type DeleteDictTypeForm struct {
	Code string `json:"code"`
}

type CreateDictForm struct {
	TypeCode string `json:"type_code" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Value    string `json:"value"`
	IsEdit   int    `json:"is_edit"`
}

type DictListForm struct {
	Code string `json:"code"`
}

type ChangeDictForm struct {
	TypeCode string `json:"type_code"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

type DeleteDictForm struct {
	TypeCode string `json:"type_code"`
	Code     string `json:"code"`
}
