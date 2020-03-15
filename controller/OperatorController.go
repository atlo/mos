package controller

import (
	"github.com/valyala/fasthttp"
	"html/template"
	h "mos/helper"
	"mos/model/list"
	m "mos/model"
	"mos/db"
	"strconv"
	"fmt"
	"mos/model/view/admin"
	"mos/model/view"
)

type OperatorController struct {
	AuthAction map[string][]string;
}

func (oc *OperatorController) Init() {
	oc.AuthAction = make(map[string][]string);
	oc.AuthAction["edit"] = []string{"operator/edit"};
	oc.AuthAction["save"] = []string{"operator/edit", "operator/new"};
	oc.AuthAction["list"] = []string{"operator/list"};
}

func (oc *OperatorController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(oc.AuthAction["list"],session)) {
		var ol list.OperatorList;
		ol.Init(ctx, session.GetActiveLang());
		pageInstance.Title = "List Owners"

		AdminContent := admin.Content{};
		AdminContent.Title = "Owners"
		AdminContent.SubTitle = "List Owners"

		AdminContent.Content = template.HTML(ol.Render(ol.GetToPage()))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true,pageInstance);
		return;
	}
}

func (oc *OperatorController) EditAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(oc.AuthAction["edit"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var operatorId = int64(id);
		var operator m.Operator;
		operator, err := operator.Get(operatorId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "operator/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		var data map[string]interface{};
		if (!ctx.IsPost()) {
			data = map[string]interface{}{
				"id":         strconv.Itoa(int(operator.Id)),
				"name": operator.Name,
				"evolution_date": operator.EvolutionDate,
				"registration_date": operator.RegistrationDate,
				"termination_date": operator.TerminationDate,
			};
		} else {
			data = map[string]interface{}{
				"id":         h.GetFormData(ctx, "id", false).(string),
				"name": h.GetFormData(ctx, "name", false).(string),
				"evolution_date":  h.GetFormData(ctx, "evolution_date", false).(string),
				"registration_date":  h.GetFormData(ctx, "registration_date", false).(string),
				"termination_date":  h.GetFormData(ctx, "termination_date", false).(string),
			};
		}

		var form = m.GetOperatorForm(data, fmt.Sprintf("operator/edit/%v", data["id"].(string)));
		if (ctx.IsPost()) {
			succ, formErrors := oc.saveOperator(ctx, session, operator);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Operator save was successful.");
				Redirect(ctx, fmt.Sprintf("operator/edit/%v", data["id"].(string)), fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Operator - Edit"

		AdminContent := admin.Content{};
		AdminContent.Title = "Operator"
		AdminContent.SubTitle = fmt.Sprintf("Edit operator %v", operator.Name);
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "operator/index", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (oc *OperatorController) saveOperator(ctx *fasthttp.RequestCtx, session *h.Session, operator m.Operator) (bool, map[string]error) {
	if (ctx.IsPost() && ((Ah.HasRights(oc.AuthAction["edit"],session) && operator.Id != 0) || (Ah.HasRights(oc.AuthAction["new"],session) && operator.Id == 0))) {
		var err error;
		var succ bool;
		var Validator = m.GetOperatorFormValidator(ctx, operator);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		operator.Name = h.GetFormData(ctx, "name", false).(string);
		operator.EvolutionDate = h.GetFormData(ctx, "evolution_date", false).(string);
		operator.RegistrationDate = h.GetFormData(ctx, "registration_date", false).(string);
		operator.TerminationDate = h.GetFormData(ctx, "termination_date", false).(string);

		if (operator.Id > 0) {
			_, err = db.DbMap.Update(&operator);
		} else {
			err = db.DbMap.Insert(&operator);
		}
		h.Error(err, "", h.ERROR_LVL_ERROR)
		succ = err == nil;
		return succ, nil;
	} else {
		return false, nil;
	}
}

func (oc *OperatorController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(oc.AuthAction["new"],session)) {
		var operator = m.NewEmptyOperator();
		var data map[string]interface{} = map[string]interface{}{};
		var dataKeys []string = []string{"id", "name", "evolution_date","registration_date","termination_date"};
		for _,k := range dataKeys {
			var val string = "";
			if(ctx.IsPost()) {
				val = h.GetFormData(ctx, k, false).(string)
			}
			data[k] = val;
		}

		var form = m.GetOperatorForm(data, "operator/new");
		if (ctx.IsPost()) {
			succ, formErrors := oc.saveOperator(ctx, session, operator);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Operator save was successful.");
				Redirect(ctx, "operator", fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Operator - New"

		AdminContent := admin.Content{};
		AdminContent.Title = "Operator"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "operator/index", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (oc *OperatorController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(oc.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var operatorId = int64(id);
		var operator m.Operator;
		operator, err := operator.Get(operatorId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "operator/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		operatorName := operator.Name;
		count,err := db.DbMap.Delete(&operator);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil || count == 0){
			session.AddError("An error occurred, could not delete operator.");
			status = fasthttp.StatusBadRequest;
			return;
		} else {
			session.AddSuccess(fmt.Sprintf("Operator %v has been deleted",operatorName));
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "operator/index", status, true,pageInstance)
}
