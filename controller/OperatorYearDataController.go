package controller

import (
	"github.com/valyala/fasthttp"
	"html/template"
	h "mos/helper"
	"mos/model/list"
	m "mos/model"
	"mos/db"
	"strconv"
	"mos/model/view/admin"
	"mos/model/view"
	"fmt"
)

type OperatorYearDataController struct {
	AuthAction map[string][]string;
}

func (opd *OperatorYearDataController) Init() {
	opd.AuthAction = make(map[string][]string);
	opd.AuthAction["edit"] = []string{"operator/edit"};
	opd.AuthAction["delete"] = []string{"operator/edit"};
	opd.AuthAction["save"] = []string{"operator/edit", "operator/new"};
	opd.AuthAction["new"] = []string{"operator/new"};
	opd.AuthAction["list"] = []string{"operator/list"};
}

func (opd *OperatorYearDataController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(opd.AuthAction["list"], session)) {
		var oydl list.OperatorYearDataList;
		oydl.Init(ctx, session.GetActiveLang());
		pageInstance.Title = "List Operator Data"

		AdminContent := admin.Content{};
		AdminContent.Title = "Operator Data"
		AdminContent.SubTitle = "List Operator Data"

		AdminContent.Content = template.HTML(oydl.Render(oydl.GetToPage()))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true, pageInstance);
		return;
	}
}

func (opd *OperatorYearDataController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(opd.AuthAction["new"], session)) {
		var operatorYearData m.OperatorYearData;
		var data map[string]interface{} = map[string]interface{}{};
		var dataKeys []string = []string{"year","operator_id","income_net","income_tax","income_operational","address"};
		for _, k := range dataKeys {
			var val string = "";
			if (ctx.IsPost()) {
				val = h.GetFormData(ctx, k, false).(string)
			}
			data[k] = val;
		}

		var form = m.GetOperatorYearDataForm(true, data, "operatordata/new");
		if (ctx.IsPost()) {
			succ, formErrors := opd.saveOperatorData(ctx, session, operatorYearData);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Operator data save was successful.");
				Redirect(ctx, "operatordata", fasthttp.StatusOK, true, pageInstance);
				return;
			}
		}

		pageInstance.Title = "Operator data - New"

		AdminContent := admin.Content{};
		AdminContent.Title = "operator data"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "", fasthttp.StatusForbidden, true, pageInstance)
		return;
	}
}

func (opd *OperatorYearDataController) saveOperatorData(ctx *fasthttp.RequestCtx, session *h.Session, operatorData m.OperatorYearData) (bool, map[string]error) {
	if (ctx.IsPost() && Ah.HasRights(opd.AuthAction["new"], session)) {
		var err error;
		var succ bool;

		var Validator = m.GetOperatorYearDataFormValidator(ctx, operatorData);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		var newModel bool = operatorData.OperatorId == 0;

		if(operatorData.OperatorId == 0) {
			operatorId, err := strconv.Atoi(h.GetFormData(ctx, "operator_id", false).(string));
			h.Error(err, "", h.ERROR_LVL_WARNING);
			operatorData.OperatorId =int64( operatorId);
		}

		if(operatorData.Year == 0) {
			year, err := strconv.Atoi(h.GetFormData(ctx, "year", false).(string));
			h.Error(err, "", h.ERROR_LVL_WARNING);
			operatorData.Year =year;
		}

		operatorData.Address = h.GetFormData(ctx, "address", false).(string);

		incomeNet, err := strconv.Atoi(h.GetFormData(ctx, "income_net", false).(string));
		h.Error(err, "", h.ERROR_LVL_NOTICE);
		operatorData.IncomeNet = int64(incomeNet);

		incomeTax, err := strconv.Atoi(h.GetFormData(ctx, "income_tax", false).(string));
		h.Error(err, "", h.ERROR_LVL_NOTICE);
		operatorData.IncomeTax = int64(incomeTax);

		incomeOperational, err := strconv.Atoi(h.GetFormData(ctx, "income_operational", false).(string));
		h.Error(err, "", h.ERROR_LVL_NOTICE);
		operatorData.IncomeOperational = int64(incomeOperational);


		if(newModel) {
			err = db.DbMap.Insert(&operatorData)
		} else {
			_,err = db.DbMap.Update(&operatorData);
		}

		h.Error(err, "", h.ERROR_LVL_ERROR);

		succ = err == nil;
		var errRet map[string]error = nil;
		if(!succ){
			errRet = map[string]error{"operator_id":err};
		}
		return succ, errRet;
	} else {
		return false, nil;
	}
}

func (oydc *OperatorYearDataController) EditAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(oydc.AuthAction["edit"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var operatorStrId string = h.GetParamFromCtxPath(ctx, 3, "");
		var yearStr string = h.GetParamFromCtxPath(ctx, 4, "");
		var operatorId, _ = strconv.Atoi(operatorStrId);
		var year, _ = strconv.Atoi(yearStr);

		var operatorYearData m.OperatorYearData;
		operatorYearData, err := operatorYearData.Get(int64(operatorId),year);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "operatordata/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		var data map[string]interface{};
		if (!ctx.IsPost()) {
			data = map[string]interface{}{
				"operator_id": operatorStrId,
				"year": yearStr,
				"address" : operatorYearData.Address,
				"income_net" : strconv.Itoa(int(operatorYearData.IncomeNet)),
				"income_tax" : strconv.Itoa(int(operatorYearData.IncomeTax)),
				"income_operational" : strconv.Itoa(int(operatorYearData.IncomeOperational)),
			};
		} else {
			data = map[string]interface{}{
				"operator_id": operatorStrId,
				"year": yearStr,
				"address" :  h.GetFormData(ctx, "address", false).(string),
				"income_net" : h.GetFormData(ctx, "income_net", false).(string),
				"income_tax" : h.GetFormData(ctx, "income_tax", false).(string),
				"income_operational" : h.GetFormData(ctx, "income_operational", false).(string),
			};
		}

		var form = m.GetOperatorYearDataForm(false, data, "operatordata/edit", data["operator_id"].(string),data["year"].(string));
		if (ctx.IsPost()) {
			succ, formErrors := oydc.saveOperatorData(ctx, session, operatorYearData);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Operator data save was successful.");
				Redirect(ctx, fmt.Sprintf("operatordata/edit/%v/%v", data["operator_id"].(string),data["year"].(string)), fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Operator Data - Edit"

		AdminContent := admin.Content{};
		AdminContent.Title = "Operator Data"
		AdminContent.SubTitle = "Edit operator data";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "operatordata", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (opd *OperatorYearDataController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(opd.AuthAction["delete"],session)) {
		var operatorStrId string = h.GetParamFromCtxPath(ctx, 3, "");
		var yearStr string = h.GetParamFromCtxPath(ctx, 4, "");
		var operatorId, _ = strconv.Atoi(operatorStrId);
		var year, _ = strconv.Atoi(yearStr);

		var operatorYearData m.OperatorYearData;
		operatorYearData, err := operatorYearData.Get(int64(operatorId),year);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "operatordata/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		count,err := db.DbMap.Delete(&operatorYearData);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil || count == 0){
			session.AddError(err.Error());
			session.AddError("An error occurred, could not delete operator data.");
			status = fasthttp.StatusBadRequest;
		} else {
			session.AddSuccess("Operator data has been deleted.");
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "operatordata/index", status, true,pageInstance);
}
