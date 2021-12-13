package policy

import (
	"encoding/json"
	"ibtools_server/models"
)

func GetOssPolicy(path string) (string, error) {
	policy := new(models.OssPolicy)
	policy.Version = "1.0"
	statement := new(models.OssStatement)
	statement.Action = []string{"oss:PutObject", "oss:GetObject"}
	statement.Effect = "Allow"
	statement.Resource = "acs:oss:*:*:ibtools/projects/" + path
	policy.Statement = []models.OssStatement{*statement}
	if jsonbyte, err := json.Marshal(policy); err != nil {
		return "", err
	} else {
		return string(jsonbyte), nil
	}
}
