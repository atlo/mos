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

type OperatorInterestList struct {
	List m.List
}

func (oil *OperatorInterestList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var operatorInterest m.OperatorInterest;
	oil.List.Init(ctx, operatorInterest, lang)

	var operator m.Operator;
	operatorList := operator.GetOptions(map[string]string{"value":"","label":"Anything"});

	var interest m.Interest;
	interestList := interest.GetOptions(map[string]string{"value":"","label":"Anything"});

	oil.List.AddSearchParam(m.SearchParam{"Operator", "operator_id", "select", "operator_id", map[string]interface{}{"options": operatorList}});
	oil.List.AddSearchParam(m.SearchParam{"Interest", "interest_id", "select", "interest_id", map[string]interface{}{"options": interestList}});
	oil.List.AddSearchParam(m.SearchParam{"Year","year","number","year",nil});
}

func (oil *OperatorInterestList) SetLimitParam(limitParam string) {
	oil.List.SetLimitParam(limitParam)
}

func (oil *OperatorInterestList) SetPageParam(pageParam string) {
	oil.List.SetPageParam(pageParam);
}

func (oil *OperatorInterestList) Render(elements []m.OperatorInterest) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "year","title":"Year"},
		{"col":"operator_id","title":"Operator"},
		{"col":"interest_id","title":"Interest"},
		{"col":"actions","title":"Actions"},
	}

	for _, oi:= range elements {
		operator := oi.GetOperator();
		interest := oi.GetInterest();
		var actions []string;

		actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("operatorinterest/delete",  []string{strconv.Itoa(int(oi.Id))},true,"admin"), "[Delete]"},
		))

		rows = append(rows, map[string]string{
			"year":      strconv.Itoa(oi.Year),
			"operator_id":      operator.Name,
			"interest_id":      interest.Name,
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "operatorinterest-list-table",
	};
	return oil.List.Render(headers, rows, options);
}

func (oil *OperatorInterestList) GetAll() []m.OperatorInterest {
	var results []m.OperatorInterest
	var where string = oil.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", oil.List.Table, where, oil.List.GetOrder(), oil.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (oil *OperatorInterestList) GetToPage() []m.OperatorInterest {
	var results []m.OperatorInterest
	var where string = oil.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", oil.List.Table, where, oil.List.GetOrder(), oil.List.GetOrderDir(), oil.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
