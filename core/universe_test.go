// Copyright 2024 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"os"
	"testing"
)

func TestNewUniverse(t *testing.T) {
	const testDBName = "universe_test.db"

	// 删除测试数据库文件，以确保测试从干净的状态开始
	os.Remove(testDBName)

	// 初始化 Universe 并获取 UDB 实例
	universe, err := NewUniverse(testDBName)
	if err != nil {
		t.Fatalf("NewUniverse failed: %v", err)
	}
	defer universe.DB.CloseDB()
	defer os.Remove(testDBName)

	// 检查数据库是否初始化正确
	if universe.DB == nil {
		t.Fatalf("Expected non-nil DB, got nil")
	}

	// 这里可以添加更多针对 Recv 方法的测试
}
