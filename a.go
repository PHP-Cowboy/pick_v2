package main

import "fmt"

type Group struct {
	Id       int
	ParentId int
	Name     string
	Children []Group
}

func main() {
	phpArray := []Group{
		{52, 30, "宁波组", []Group{}},
		{38, 24, "战队A组", []Group{}},
		{49, 29, "C组", []Group{}},
		{40, 24, "战队C组", []Group{}},
		{3, 1, "计划管理部", []Group{}},
		{7, 1, "财务部", []Group{}},
		{37, 23, "wayfair组", []Group{}},
		{16, 4, "研发一部", []Group{}},
		{25, 9, "第二事业部", []Group{}},
		{44, 26, "运营C组", []Group{}},
		{48, 29, "B组", []Group{}},
		{6, 1, "某某学院", []Group{}},
		{42, 26, "运营A组", []Group{}},
		{2, 1, "会议室", []Group{}},
		{19, 5, "技术开发部", []Group{}},
		{32, 11, "灯具组", []Group{}},
		{51, 30, "新战队", []Group{}},
		{10, 1, "供应链", []Group{}},
		{35, 14, "计划三组", []Group{}},
		{53, 30, "深圳组", []Group{}},
		{47, 29, "A组", []Group{}},
		{36, 23, "日本亚马逊组", []Group{}},
		{8, 1, "人力行政部", []Group{}},
		{5, 1, "IT中心", []Group{}},
		{24, 9, "第一事业部", []Group{}},
		{9, 1, "运营部", []Group{}},
		{13, 1, "总经办", []Group{}},
		{41, 25, "家具B组", []Group{}},
		{12, 1, "市场部", []Group{}},
		{21, 8, "人力资源部", []Group{}},
		{29, 10, "供应链一部", []Group{}},
		{17, 5, "架构部", []Group{}},
		{4, 1, "技术部", []Group{}},
		{22, 8, "行政部", []Group{}},
		{23, 9, "拓展部", []Group{}},
		{11, 1, "研发二部", []Group{}},
		{28, 10, "海外仓", []Group{}},
		{45, 27, "灯具品类设计组", []Group{}},
		{33, 14, "计划一组", []Group{}},
		{50, 29, "D组", []Group{}},
		{46, 27, "家具品类设计组", []Group{}},
		{18, 5, "信息服务部", []Group{}},
		{1, 0, "某某有限公司", []Group{}},
		{39, 24, "战队B组", []Group{}},
		{26, 9, "第三事业部", []Group{}},
		{34, 14, "计划二组", []Group{}},
		{14, 3, "计划部", []Group{}},
		{31, 11, "家具组", []Group{}},
		{27, 9, "设计组", []Group{}},
		{20, 5, "数据部", []Group{}},
		{43, 26, "运营B组", []Group{}},
		{30, 10, "供应链二部", []Group{}},
		{15, 4, "品质部", []Group{}},
		{52, 30, "宁波组", []Group{}},
		{38, 24, "战队A组", []Group{}},
		{49, 29, "C组", []Group{}},
		{40, 24, "战队C组", []Group{}},
		{3, 1, "计划管理部", []Group{}},
		{7, 1, "财务部", []Group{}},
		{37, 23, "wayfair组", []Group{}},
		{16, 4, "研发一部", []Group{}},
		{25, 9, "第二事业部", []Group{}},
		{44, 26, "运营C组", []Group{}},
		{48, 29, "B组", []Group{}},
		{6, 1, "某某学院", []Group{}},
		{42, 26, "运营A组", []Group{}},
		{2, 1, "会议室", []Group{}},
		{19, 5, "技术开发部", []Group{}},
		{32, 11, "灯具组", []Group{}},
		{51, 30, "新战队", []Group{}},
		{10, 1, "供应链", []Group{}},
		{35, 14, "计划三组", []Group{}},
		{53, 30, "深圳组", []Group{}},
		{47, 29, "A组", []Group{}},
		{36, 23, "日本亚马逊组", []Group{}},
		{8, 1, "人力行政部", []Group{}},
		{5, 1, "IT中心", []Group{}},
		{24, 9, "第一事业部", []Group{}},
		{9, 1, "运营部", []Group{}},
		{13, 1, "总经办", []Group{}},
		{41, 25, "家具B组", []Group{}},
		{12, 1, "市场部", []Group{}},
		{21, 8, "人力资源部", []Group{}},
		{29, 10, "供应链一部", []Group{}},
		{17, 5, "架构部", []Group{}},
		{4, 1, "技术部", []Group{}},
		{22, 8, "行政部", []Group{}},
		{23, 9, "拓展部", []Group{}},
		{11, 1, "研发二部", []Group{}},
		{28, 10, "海外仓", []Group{}},
		{45, 27, "灯具品类设计组", []Group{}},
		{33, 14, "计划一组", []Group{}},
		{50, 29, "D组", []Group{}},
		{46, 27, "家具品类设计组", []Group{}},
		{18, 5, "信息服务部", []Group{}},
		{1, 0, "某某有限公司", []Group{}},
		{39, 24, "战队B组", []Group{}},
		{26, 9, "第三事业部", []Group{}},
		{34, 14, "计划二组", []Group{}},
		{14, 3, "计划部", []Group{}},
		{31, 11, "家具组", []Group{}},
		{27, 9, "设计组", []Group{}},
		{20, 5, "数据部", []Group{}},
		{43, 26, "运营B组", []Group{}},
		{30, 10, "供应链二部", []Group{}},
		{15, 4, "品质部", []Group{}},
	}

	tree := buildTree(phpArray, 0)

	fmt.Println(tree)
}

func buildTree(groups []Group, parentId int) []Group {
	var tree []Group
	for _, group := range groups {
		if group.ParentId == parentId {
			children := buildTree(groups, group.Id)
			if children != nil {
				group.Children = children
			}
			tree = append(tree, group)
		}
	}
	return tree
}
