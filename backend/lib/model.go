package lib

type DatabaseConfig struct {
	Type     string `json:"type"` // 数据库类型
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"` // 数据库名
	Migrate  bool   `json:"migrate"`
}
