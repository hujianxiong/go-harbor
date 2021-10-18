/*
Copyright 2020 The go-harbor Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
*/

package example

import (
	"fmt"
	"github.com/hujianxiong/go-harbor"
	"github.com/hujianxiong/go-harbor/pkg/model"
)

func Users(host, username, password string) error {
	clientSet, err := harbor.NewClientSet(host, username, password)
	if err != nil {
		return fmt.Errorf("get client set error:%v", err)
	}
	result, err := clientSet.User.Get("1")
	if err != nil || len(result.Username) == 0 {
		return fmt.Errorf("%v", err)
	}
	query := model.Query{
		PageSize: 2,
	}
	result1, err := clientSet.User.List(&query)
	if err != nil || len(*result1) == 0 {
		return fmt.Errorf("%v", err)
	}

	err = clientSet.User.Delete("3")
	if err != nil || len(*result1) == 0 {
		return fmt.Errorf("%v", err)
	}
	return err
}
