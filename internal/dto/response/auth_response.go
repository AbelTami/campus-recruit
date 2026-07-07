package response

type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	UserInfo     UserInfo `json:"userInfo"`
}

type UserInfo struct {
	ID       uint64   `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Avatar   string   `json:"avatar"`
	Roles    []string `json:"roles"`
}

type RouteMeta struct {
	Title      string   `json:"title"`
	Icon       string   `json:"icon,omitempty"`
	Hidden     bool     `json:"hidden,omitempty"`
	AlwaysShow bool     `json:"alwaysShow,omitempty"`
	NoCache    bool     `json:"noCache,omitempty"`
	Affix      bool     `json:"affix,omitempty"`
	Permission []string `json:"permission,omitempty"`
}

type MenuTreeResponse struct {
	ID        uint64             `json:"id"`
	ParentID  *uint64            `json:"parentId"`
	Type      int16              `json:"type"`
	Path      string             `json:"path"`
	Name      string             `json:"name"`
	Title     string             `json:"title"`
	Redirect  string             `json:"redirect,omitempty"`
	Component string             `json:"component"`
	Status    int16              `json:"status"`
	Meta      RouteMeta          `json:"meta"`
	Children  []MenuTreeResponse `json:"children,omitempty"`
}
