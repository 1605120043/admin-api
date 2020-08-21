package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/shinmigo/pb/productpb"
	"goshop/api/pkg/validation"
	"goshop/api/service"
	"regexp"
	"strconv"
)

type ProductKind struct {
	*gin.Context
}

func NewProductKind(c *gin.Context) *ProductKind  {
	return &ProductKind{Context: c}
}

func (m *ProductKind) Index() (*productpb.ListKindRes, error)  {
	valid := validation.Validation{}
	page := m.DefaultQuery("page", "1")
	pageSize := m.DefaultQuery("page_size", "10")
	valid.Match(page, regexp.MustCompile(`^[0-9]{1,3}$`)).Message("页面的编号 不正确")
	valid.Match(pageSize, regexp.MustCompile(`^[0-9]{1,3}$`)).Message("页面的数量 不正确")
	if valid.HasError() {
		return nil, valid.GetError()
	}

	pNumber, _ := strconv.ParseUint(page, 10, 16)
	pSize, _ := strconv.ParseUint(pageSize, 10, 16)

	list, err := service.NewProductKind(m.Context).Index(pNumber, pSize)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *ProductKind) Add() error  {
	valid := validation.Validation{}
	name := m.PostForm("name")

	valid.Required(name).Message("名称不能为空")
	if valid.HasError() {
		return valid.GetError()
	}

	if err := service.NewProductKind(m.Context).Add(); err != nil {
		return err
	}

	return nil
}

func (m *ProductKind) Delete() error  {
	id := m.PostForm("kind_id")

	valid := validation.Validation{}
	valid.Required(id).Message("请选择要删除的商品类型！")
	if valid.HasError() {
		return valid.GetError()
	}

	if err := service.NewProductKind(m.Context).Delete(); err != nil {
		return err
	}

	return nil

}

func (m *ProductKind) Edit() error  {
	storeId := m.PostForm("store_id")
	name := m.PostForm("name")
	kindId := m.PostForm("kind_id")

	valid := validation.Validation{}
	valid.Required(kindId).Message("请提交要编辑的商品类型")
	valid.Match(kindId, regexp.MustCompile(`^[1-9][0-9]*$`)).Message("商品类型数据格式错误")
	valid.Required(storeId).Message("请提交商品类型所属门店！")
	valid.Match(storeId, regexp.MustCompile(`^[1-9][0-9]*$`)).Message("商品类型门店数据格式错误！")
	valid.Required(name).Message("请提交商品类型名称")
	if valid.HasError() {
		return valid.GetError()
	}

	if err := service.NewProductKind(m.Context).Edit(); err != nil {
		return err
	}

	return nil
}