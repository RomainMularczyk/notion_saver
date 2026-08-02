package notionModels

type NotionIcon struct {
	Type            string                 `json:"type"`
	NativeIcon      *NotionNativeIcon      `json:"native_icon,omitempty"`
	CustomEmojiIcon *NotionCustomEmojiIcon `json:"custom_emoji_icon,omitempty"`
	EmojiIcon       *NotionEmojiIcon       `json:"emoji_icon,omitempty"`
}

type NotionNativeIcon struct {
	Icon struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"icon"`
}

type NotionCustomEmojiIcon struct {
	Icon struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"icon"`
}

type NotionEmojiIcon struct {
	Icon struct {
		Emoji string `json:"emoji"`
	} `json:"icon"`
}
