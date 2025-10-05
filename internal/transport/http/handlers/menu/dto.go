package menu_handler

type addCategoryReq struct {
	Position     int               `json:"position" example:"1"`
	Translations map[string]string `json:"translations" example:"{en:Category Name}"`
}
