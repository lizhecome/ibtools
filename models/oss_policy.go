package models

type OssPolicy struct {
	Statement []OssStatement `json:"Statement"`
	Version   string         `json:"Version"`
}
type OssStatement struct {
	Action   []string `json:"Action"`
	Effect   string   `json:"Effect"`
	Resource string   `json:"Resource"`
}
