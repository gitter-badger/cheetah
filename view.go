package cheetah

type View struct {
	Title       string // page's title
	Keywords    string // page's keywords
	Description string // page's description
	HeaderCss   []*CssAsset
	HeaderJs    []*JsAsset
	FooterCss   []*CssAsset
	FooterJs    []*JsAsset
}

func NewView(title, keywords, description string) *View {
	return &View{
		Title       :title,
		Keywords    :keywords,
		Description :description,
		HeaderCss :make([]*CssAsset, 0),
		HeaderJs    :make([]*JsAsset, 0),
		FooterCss   :make([]*CssAsset, 0),
		FooterJs    :make([]*JsAsset, 0),
	}
}