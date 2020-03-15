package list

import (
	"github.com/valyala/fasthttp"
	"mos/db"
	"fmt"
	h "mos/helper"
	m "mos/model"
	"strconv"
	"strings"
)

type MediaTypeList struct {
	List m.List
}

func (mtl *MediaTypeList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var mediaType m.MediaType;
	mtl.List.Init(ctx, mediaType, lang)

	mtl.List.AddSearchParam(m.SearchParam{"ID","id","number","id", nil});
	mtl.List.AddSearchParam(m.SearchParam{"Name","name","text","name",nil});
}

func (mtl *MediaTypeList) SetLimitParam(limitParam string) {
	mtl.List.SetLimitParam(limitParam)
}

func (mtl *MediaTypeList) SetPageParam(pageParam string) {
	mtl.List.SetPageParam(pageParam);
}

func (mtl *MediaTypeList) Render(elements []m.MediaType) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "id","title":"ID"},
		{"col":"name","title":"Name"},
		{"col":"actions","title":"Actions"},
	}

	for _, mt:= range elements {
		var actions []string;
		actions = append(actions, h.Replace(
			`<a href="%link%">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("mediatype/edit",  []string{strconv.Itoa(int(mt.Id))},true,"admin"), "[Edit]"},
		))

		actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("mediatype/delete",  []string{strconv.Itoa(int(mt.Id))},true,"admin"), "[Delete]"},
		))

		rows = append(rows, map[string]string{
			"id":         strconv.Itoa(int(mt.Id)),
			"name":      mt.Name,
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "mediatype-list-table",
	};
	return mtl.List.Render(headers, rows, options);
}

func (mtl *MediaTypeList) GetAll() []m.MediaType {
	var results []m.MediaType
	var where string = mtl.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", mtl.List.Table, where, mtl.List.GetOrder(), mtl.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (mtl *MediaTypeList) GetToPage() []m.MediaType {
	var results []m.MediaType
	var where string = mtl.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", mtl.List.Table, where, mtl.List.GetOrder(), mtl.List.GetOrderDir(), mtl.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
