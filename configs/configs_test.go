/*
@Time : 2018/9/7 上午11:59 
@Author : Mingjian Song
@File : configs_test
@Software: 深圳超链科技
*/

package configs

import (
	"testing"
	"fmt"
)

func TestConfig_Load(t *testing.T) {
	cfg := Init()
	cfg.Load("../configs.json")
	fmt.Println(cfg)
}
