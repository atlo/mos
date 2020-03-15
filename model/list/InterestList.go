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

type InterestList struct {
	List m.List
}

func (il *InterestList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var interest m.Interest;
	il.List.Init(ctx, interest, lang)

	il.List.AddSearchParam(m.SearchParam{"ID","id","number","id", nil});
	il.List.AddSearchParam(m.SearchParam{"Name","name","text","name",nil});
}

func (il *InterestList) SetLimitParam(limitParam string) {
	il.List.SetLimitParam(limitParam)
}

func (il *InterestList) SetPageParam(pageParam string) {
	il.List.SetPageParam(pageParam);
}

func (il *InterestList) Render(elements []m.Interest) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "id","title":"ID","order":"true"},
		{"col":"name","title":"Name"},
		{"col":"actions","title":"Actions"},
	}

	for _, o := range elements {
		var actions []string;
		actions = append(actions, h.Replace(
			`<a href="%link%">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("interest/edit",  []string{strconv.Itoa(int(o.Id))},true,"admin"), "[Edit]"},
		))

		actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("interest/delete",  []string{strconv.Itoa(int(o.Id))},true,"admin"), "[Delete]"},
		))

		rows = append(rows, map[string]string{
			"id":         strconv.Itoa(int(o.Id)),
			"name":      o.Name,
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "interest-list-table",
	};
	return il.List.Render(headers, rows, options);
}

func (il *InterestList) GetAll() []m.Interest {
	var results []m.Interest
	var where string = il.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", il.List.Table, where, il.List.GetOrder(), il.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (il *InterestList) GetToPage() []m.Interest {
	var results []m.Interest
	var where string = il.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", il.List.Table, where, il.List.GetOrder(), il.List.GetOrderDir(), il.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
