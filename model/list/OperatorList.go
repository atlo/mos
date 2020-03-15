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

type OperatorList struct {
	List m.List
}

func (ol *OperatorList) Init(ctx *fasthttp.RequestCtx, lang string) {
	var operator m.Operator;
	ol.List.Init(ctx, operator, lang)

	ol.List.AddSearchParam(m.SearchParam{"ID","id","number","id", nil});
	ol.List.AddSearchParam(m.SearchParam{"Name","name","text","name",nil});
	ol.List.AddSearchParam(m.SearchParam{"Evolution Date","evolution_date","string","evolution_date",nil});
	ol.List.AddSearchParam(m.SearchParam{"Registration Date","registration_date","string","registration_date",nil});
	ol.List.AddSearchParam(m.SearchParam{"Termination Date","termination_date","string","termination_date",nil});
}

func (ol *OperatorList) SetLimitParam(limitParam string) {
	ol.List.SetLimitParam(limitParam)
}

func (ol *OperatorList) SetPageParam(pageParam string) {
	ol.List.SetPageParam(pageParam);
}

func (ol *OperatorList) Render(elements []m.Operator) string {
	var headers []map[string]string;
	var rows []map[string]string;
	var options map[string]string;
	headers = []map[string]string{
		{"col": "id","title":"ID","order":"true"},
		{"col":"name","title":"Name"},
		{"col":"evolution_date","title":"Evolution Date"},
		{"col":"registration_date","title":"Registration Date"},
		{"col":"termination_date","title":"Termination Date"},
		{"col":"actions","title":"Actions"},
	}

	for _, o := range elements {
		var actions []string;
		actions = append(actions, h.Replace(
			`<a href="%link%">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("operator/edit",  []string{strconv.Itoa(int(o.Id))},true,"admin"), "[Edit]"},
		))

		/*actions = append(actions, h.Replace(
			`<a href="%link%" onclick="return window.confirm('Are you sure you want to delete the item?')">%title%</a>`,
			[]string{"%link%","%title%"},
			[]string{h.GetUrl("owner/delete",  []string{strconv.Itoa(int(o.Id))},true,"admin"), "[Delete]"},
		))*/

		rows = append(rows, map[string]string{
			"id":         strconv.Itoa(int(o.Id)),
			"name":      o.Name,
			"evolution_date": o.EvolutionDate,
			"registration_date": o.RegistrationDate,
			"termination_date": o.TerminationDate,
			"actions":strings.Join(actions,"&nbsp;&nbsp;"),
		});
	}

	options = map[string]string{
		"class": "table-striped table-bordered table-hover",
		"id":    "operator-list-table",
	};

	return ol.List.Render(headers, rows, options);
}

func (ol *OperatorList) GetAll() []m.Operator {
	var results []m.Operator
	var where string = ol.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v", ol.List.Table, where, ol.List.GetOrder(), ol.List.GetOrderDir());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}

func (ol *OperatorList) GetToPage() []m.Operator {
	var results []m.Operator
	var where string = ol.List.GetSqlParams();
	if(where != ""){
		where = fmt.Sprintf(" WHERE %v",where);
	}
	sql := fmt.Sprintf("SELECT * FROM %v%v ORDER BY %v %v LIMIT %v", ol.List.Table, where, ol.List.GetOrder(), ol.List.GetOrderDir(), ol.List.GetLimitString());
	h.PrintlnIf(sql,h.GetConfig().Mode.Debug);
	_, err := db.DbMap.Select(&results, sql);
	h.Error(err, "", h.ERROR_LVL_ERROR)
	return results;
}
