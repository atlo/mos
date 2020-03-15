package controller

import (
	"github.com/valyala/fasthttp"
	h "mos/helper"
	"mos/model/view"
	"mos/model"
	"fmt"
	"mos/db"
	"encoding/json"
	"database/sql"
	"strconv"
	"strings"
)

type FeedController struct {
	AuthAction map[string][]string
}

func (f *FeedController) Init() {
	f.AuthAction = make(map[string][]string)
	f.AuthAction["media"] = []string{"*"}
	f.AuthAction["interests"] = []string{"*"}
	f.AuthAction["operators"] = []string{"*"}
	f.AuthAction["owners"] = []string{"*"}
	f.AuthAction["connections"] = []string{"*"}
}

func (f FeedController) MediaAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(f.AuthAction["media"], session)) {
		pageInstance.Layout = "ajax.html"
		pageInstance.ContentType = "application/json;charset=utf-8";

		type MediaJson struct{
			Id int64 `json:"_id"`
			Name string `json:"name"`
			Type string `json:"type"`
			News bool `json:"news_non_news"`
		}

		var results []MediaJson;

		var media model.Media;
		var mediaType model.MediaType;
		var query string = fmt.Sprintf("SELECT id,name,(SELECT name FROM %v WHERE id = `%v`.`media_type_id`) as type,news FROM %v",mediaType.GetTable(),media.GetTable(),media.GetTable());
		var err error;

		_,err = db.DbMap.Select(&results,query);
		h.Error(err,"",h.ERROR_LVL_ERROR);

		bytesJson,err := json.Marshal(results)
		h.Error(err,"",h.ERROR_LVL_ERROR);

		pageInstance.AddContent(string(bytesJson),"",nil,false,0);
	}
	return;
}

func (f FeedController) InterestsAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(f.AuthAction["interests"], session)) {
		pageInstance.Layout = "ajax.html"
		pageInstance.ContentType = "application/json;charset=utf-8";

		var results []model.Interest;

		var interest model.Interest;
		var query string = fmt.Sprintf("SELECT * FROM %s",interest.GetTable());
		var err error;

		_,err = db.DbMap.Select(&results,query);
		h.Error(err,"",h.ERROR_LVL_ERROR);

		bytesJson,err := json.Marshal(results)
		h.Error(err,"",h.ERROR_LVL_ERROR);

		pageInstance.AddContent(string(bytesJson),"",nil,false,0);
	}
	return;
}

func (f FeedController) OwnersAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(f.AuthAction["owners"], session)) {
		pageInstance.Layout = "ajax.html"
		pageInstance.ContentType = "application/json;charset=utf-8";

		type OwnerJson struct{
			Id int64 `json:"_id"`
			Name string `json:"name"`
			Hungarian bool `json:"hun_non_hun"`
		}

		var results []OwnerJson;
		var owner model.Owner;
		var query string = fmt.Sprintf("SELECT id,name,hungarian FROM %v",owner.GetTable());
		var err error;

		_,err = db.DbMap.Select(&results,query);
		h.Error(err,"",h.ERROR_LVL_ERROR);

		bytesJson,err := json.Marshal(results)
		h.Error(err,"",h.ERROR_LVL_ERROR);

		pageInstance.AddContent(string(bytesJson),"",nil,false,0);
	} else {
		Redirect(ctx,"user/login",fasthttp.StatusForbidden,true,pageInstance);
	}
	return;
}

func (f FeedController) OperatorsAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(f.AuthAction["operators"], session)) {
		pageInstance.Layout = "ajax.html"
		pageInstance.ContentType = "application/json;charset=utf-8";

		type OperatorDataJson struct {
			OperatorId int64 `json:"-"`
			Year int `json:"year"`
			Address string `json:"address"`
			NettoProfit int `json:"netto_profit"`
			TaxedProfit int `json:"taxed_profit"`
			OperatingProfit int `json:"operating_profit"`
			Interest []int64 `json:"interests"`
		}

		type OperatorJson struct{
			Id int64 `json:"_id"`
			Name string `json:"name"`
			Data map[int]OperatorDataJson `json:"data"`
		}

		var results []OperatorJson;
		var operator model.Operator;
		var operatorYearData model.OperatorYearData;
		var query string = fmt.Sprintf("SELECT id,name FROM %v",operator.GetTable());
		var err error;
		_,err = db.DbMap.Select(&results,query);
		h.Error(err,"",h.ERROR_LVL_ERROR);

		for i,o := range results {
			var query string = fmt.Sprintf("SELECT year,address,income_net,income_tax,income_operational FROM %v WHERE operator_id = ?",operatorYearData.GetTable());
			var err error;
			var rows *sql.Rows;
			var resultsYears map[int]OperatorDataJson = make(map[int]OperatorDataJson);
			var resultYear OperatorDataJson;

			resultYear.OperatorId = o.Id;

			rows,err = db.DbMap.Query(query,o.Id);
			h.Error(err,"",h.ERROR_LVL_ERROR);
			for rows.Next(){
				err = rows.Scan(&resultYear.Year,&resultYear.Address,&resultYear.NettoProfit,&resultYear.TaxedProfit,&resultYear.OperatingProfit);
				h.Error(err,"",h.ERROR_LVL_WARNING);

				var interest model.Interest;
				var operatorInterest model.OperatorInterest;
				var interestQuery string = fmt.Sprintf("SELECT %s FROM %s WHERE %s IN (" +
					"SELECT %s FROM %s WHERE %s = ? AND %s = ?" +
				")","id",interest.GetTable(),interest.GetPrimaryKey()[0],"interest_id",operatorInterest.GetTable(),"operator_id","year")

				resultYear.Interest = nil;
				_,err := db.DbMap.Select(&resultYear.Interest,interestQuery,o.Id,resultYear.Year);
				h.Error(err,"",h.ERROR_LVL_WARNING);

				resultsYears[resultYear.Year] = resultYear;
			}
			rows.Close();

			results[i].Data = resultsYears;
		}

		bytesJson,err := json.Marshal(results)
		h.Error(err,"",h.ERROR_LVL_ERROR);

		pageInstance.AddContent(string(bytesJson),"",nil,false,0);
	} else {
		Redirect(ctx,"user/login",fasthttp.StatusForbidden,true,pageInstance);
	}
	return;
}

func (f FeedController) ConnectionsAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(f.AuthAction["connections"], session)) {
		pageInstance.Layout = "ajax.html"
		pageInstance.ContentType = "application/json;charset=utf-8";

		type ConnectionEntityYearJson struct {
			OperatorId int64 `json:"_id"`
			Data map[string]interface{} `json:"data"`
		}

		var results []ConnectionEntityYearJson;
		var mediaOwner model.MediaOwner;
		var operator model.Operator;
		var mediaOperator model.MediaOperator;

		var operators []model.Operator = operator.GetAll();

		for _,o := range operators{
			var jsonOp ConnectionEntityYearJson;
			jsonOp.Data = make(map[string]interface{});
			jsonOp.OperatorId = o.Id;

			queryMedOp := fmt.Sprintf("SELECT GROUP_CONCAT(media_id),year FROM %v WHERE operator_id = ? GROUP BY `year`", mediaOperator.GetTable());
			rowsMedOp,err := db.DbMap.Query(queryMedOp,o.Id);
			for rowsMedOp.Next(){
				var mediaIdStr string;
				var ownerIdStr string;
				var year int;

				err = rowsMedOp.Scan(&mediaIdStr,&year);
				h.Error(err,"",h.ERROR_LVL_WARNING);

				if(mediaIdStr != ""){
					jsonOp.Data[strconv.Itoa(year)] = make(map[string]interface{});
					jsonOp.Data[strconv.Itoa(year)].(map[string]interface{})["mediaIds"] = strings.Split(mediaIdStr,",");

					queryMediaOpOw := fmt.Sprintf("SELECT GROUP_CONCAT(owner_id) FROM %v WHERE media_id IN (%v) AND `year` = ? GROUP BY `year`", mediaOwner.GetTable(), mediaIdStr);
					rowsMedOw := db.DbMap.QueryRow(queryMediaOpOw,year);
					err = rowsMedOw.Scan(&ownerIdStr);
					h.Error(err,"",h.ERROR_LVL_WARNING)
					if(ownerIdStr != ""){
						jsonOp.Data[strconv.Itoa(year)].(map[string]interface{})["ownerIds"] = strings.Split(ownerIdStr,",");
					}
				}
			}

			results = append(results, jsonOp);
		}
		bytesJson,err := json.Marshal(results);
		h.Error(err,"",h.ERROR_LVL_ERROR);
		pageInstance.AddContent(string(bytesJson),"",nil,false,0);
	} else {
		Redirect(ctx,"user/login",fasthttp.StatusForbidden,true,pageInstance);
	}
	return;
}