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

type OperatorYearDataList struct {
	List m.List
}

func (oydl *OperatorYearDataList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var operatorYearData m.OperatorYearData;
	oydl.List.Init(ctx, operatorYearData, lang)

	var operator m.Operator;
	opaeratorList := operator.GetOptions(map[string]string{"value":"","label":"Anything"});

	oydl.List.AddSearchParam(m.SearchParam{"Operator", "operator_id", "select", "operator_id", map[string]interface{}{"options": opaeratorList}});
	oydl.List.AddSearchParam(m.SearchParam{"Year","year","number","year",nil});
	oydl.List.AddSearchParam(m.SearchParam{"Address","address","text","address",nil});
	oydl.List.AddSearchParam(m.SearchParam{"Income without tax","income_net","number_range","income_net",nil});
	oydl.List.AddSearchParam(m.SearchParam{"Income with tax","income_tax","number_range","income_tax",nil});
	oydl.List.AddSearchParam(m.SearchParam{"Income operational","income_operational","number_range","income_tax",nil});
}

func (oydl *OperatorYearDataList) SetLimitParam(limitParam string) {
	oydl.List.SetLimitParam(limitParam)
}

func (oydl *OperatorYearDataList) SetPageParam(pageParam string) {
	oydl.List.SetPageParam(pageParam);
}

func (oydl *OperatorYearDataList) Render(elements []m.OperatorYearData) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "year","title":"Year"},
		{"col":"operator_id","title":"Operator"},
		{"col":"address","title":"Address"},
		{"col":"income_net","title":"Income without tax"},
		{"col":"income_tax","title":"Income with tax"},
		{"col":"income_operational","title":"Income operational"},
		{"col":"actions","title":"Actions"},
	}

	for _, oyd:= range elements {
		operator := oyd.GetOperator();
		var actions []string;

		actions = append(actions, h.Replace(
			`<a href="%link%">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("operatordata/edit",  []string{strconv.Itoa(int(oyd.OperatorId)),strconv.Itoa(int(oyd.Year))},true,"admin"), "[Edit]"},
		))

		actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("operatordata/delete",  []string{strconv.Itoa(int(oyd.OperatorId)),strconv.Itoa(int(oyd.Year))},true,"admin"), "[Delete]"},
		))

		rows = append(rows, map[string]string{
			"year":      strconv.Itoa(oyd.Year),
			"operator_id":      operator.Name,
			"address":      oyd.Address,
			"income_net":      strconv.Itoa(int(oyd.IncomeNet)),
			"income_tax":      strconv.Itoa(int(oyd.IncomeTax)),
			"income_operational":      strconv.Itoa(int(oyd.IncomeOperational)),
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "operatoryeardata-list-table",
	};
	return oydl.List.Render(headers, rows, options);
}

func (oydl *OperatorYearDataList) GetAll() []m.OperatorYearData {
	var results []m.OperatorYearData
	var where string = oydl.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", oydl.List.Table, where, oydl.List.GetOrder(), oydl.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (oydl *OperatorYearDataList) GetToPage() []m.OperatorYearData {
	var results []m.OperatorYearData
	var where string = oydl.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", oydl.List.Table, where, oydl.List.GetOrder(), oydl.List.GetOrderDir(), oydl.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
