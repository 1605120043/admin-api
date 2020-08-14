package controller

import (
	"goshop/api/filter"
)

var productParamFilter *filter.ProdcutParam

type ProductParam struct {
	Base
}

func (m *ProductParam) Initialise() {
	productParamFilter = filter.NewProductParam(m.Context)
}

// 商品参数列表
func (m *ProductParam) Index() {
	str, err := productParamFilter.Index()
	if err != nil {
		m.SetResponse(str, err)
		return
	}
	
	m.SetResponse(str)
}

// 添加商品参数
func (m *ProductParam) Add() {
	if err := productParamFilter.Add(); err != nil {
		m.SetResponse(nil, err)
		return
	}
	m.SetResponse()
}

// 编辑商品参数
func (m *ProductParam) Edit() {
	if err := productParamFilter.Edit(); err != nil {
		m.SetResponse(nil, err)
		return
	}
	m.SetResponse()
}
