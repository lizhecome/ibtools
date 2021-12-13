package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model        `json:"-"`
	OauthUserID       uint          `json:"-"`     //项目所有者ID
	Tilte             string        `json:"title"` //项目名称
	DDModelList       []DDModel     `json:"dd_model_list"`
	IsTemplate        int           `json:"is_template"`
	Users             []OauthUser   `gorm:"many2many:project_users;" json:"users"`
	PicUrl            string        `json:"pic_url"`
	Code              uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()" json:"code"`
	Type              string        `json:"type"`
	CompanyFullName   string        `json:"company_full_name"`   //公司全称
	CompanyShortName  string        `json:"company_short_name"`  //公司简称
	ScrcIndustryType1 string        `json:"scrc_industry_type1"` //证监会一级行业分类
	ScrcIndustryType2 string        `json:"scrc_industry_type2"` //证监会二级行业分类
	MainBusinessInfo  string        `json:"main_business_info"`  //主营业务简介
	MainProductsInfo  string        `json:"main_products_info"`  //主营产品简介
	ModifyTime        time.Time     `json:"modify_time"`         //修改时间
	ModifyUserID      uint          `json:"-"`
	ModifyUser        OauthUser     `json:"modify_user"`      //修改人
	CollectProgress   string        `json:"collect_progress"` //底稿收集进度
	ProjectFiles      []ProjectFile `json:"project_files"`    //项目文件
}

type DDModel struct {
	gorm.Model `json:"-"`
	Title      string           `json:"title"`
	ProjectID  uint             `json:"-"`
	Order      int              `json:"order"`
	Code       uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4()" json:"code"`
	DDItems    []DDItem         `gorm:"polymorphic:Owner;" json:"dd_Items"`
	Questions  []ReviewQuestion `json:"questions"`
	Path       string           `json:"path"` //路径
}

type DDItem struct {
	gorm.Model     `json:"-"`
	Order          int            `json:"order"`
	Title          string         `json:"title"`
	UserID         uint           `json:"-"`
	User           OauthUser      `json:"user"`
	BaseInfos      datatypes.JSON `gorm:"default:'{}'" json:"base_infos"`
	BaseInfoSchema datatypes.JSON `gorm:"default:'{}'" json:"base_infos_schema"`
	ReviewMethod   string         `json:"review_method"`
	Comments       []Comment      `json:"comments"`
	CollectFiles   []DDFile       `json:"collect_files"`
	FilePointer    string         `json:"file_pointer"`
	Status         string         `json:"status"`
	Events         []DDEvent      `json:"events"`
	Code           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"code"`
	OwnerID        uint           `json:"-"`
	OwnerType      string         `json:"owner_type"`
}

type ProjectFile struct {
	gorm.Model `json:"-"`
	ProjectID  uint      `json:"-"`
	FileName   string    `json:"file_name"`
	Status     string    `json:"status"`
	OSSPath    string    `json:"oss_path"`
	Code       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"code"`
}

type DDFile struct {
	gorm.Model `json:"-"`
	DDItemID   uint      `json:"-"`
	FileName   string    `json:"file_name"`
	Status     string    `json:"status"`
	OSSPath    string    `json:"oss_path"`
	Code       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"code"`
}

type DDEvent struct {
	gorm.Model        `json:"-"`
	DDItemID          uint      `json:"-"`
	FiledName         string    `json:"filed_name"`
	BeforeChangeValue string    `json:"before_change_value"`
	AfterChangeValue  string    `json:"after_change_value"`
	Code              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"code"`
}

type ReviewQuestion struct {
	gorm.Model          `json:"-"`
	DDModelID           uint      `json:"-"`
	Order               int       `json:"order"`
	QuestionName        string    `json:"question_name"`
	QuestionDescription string    `json:"question_description"`
	LawSupport          string    `json:"law_support"`
	Solution            string    `json:"solution"`
	AdditionDDItem      []DDItem  `gorm:"polymorphic:Owner;" json:"addition_dd_item"`
	Code                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"code"`
}

type ProjectUser struct {
	ProjectID   uint `gorm:"primaryKey"`
	OauthUserID uint `gorm:"primaryKey"`
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt
	RoleName    string
}
