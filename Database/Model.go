package Database

type ConfigFileInfo struct {
	Host     string `json:"host"`
	Port     int32  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}
type UserInfo struct {
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type PostInfo struct {
	Uuid          string `json:"uuid"`
	User_id       string `json:"user_id"`
	Post_ordered  int    `json:"post_ordered"`
	Comment_count int    `json:"comment_count"`
	Like_count    int    `json:"like_count"`
	Dislike_count int    `json:"dislike_count"`
	Image_path    string `json:"Image_path"`
}
type RecvMessage struct {
	Status  bool
	Message string
}
