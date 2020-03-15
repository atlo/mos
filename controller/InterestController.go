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

type InterestController struct {
	AuthAction map[string][]string;
}

func (ic *InterestController) Init() {
	ic.AuthAction = make(map[string][]string);
	ic.AuthAction["edit"] = []string{"interest/edit"};
	ic.AuthAction["save"] = []string{"interest/edit", "interest/new"};
	ic.AuthAction["new"] = []string{"interest/new"};
	ic.AuthAction["delete"] = []string{"interest/delete"};
	ic.AuthAction["list"] = []string{"interest/list"};
}

func (ic *InterestController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(ic.AuthAction["list"],session)) {
		var il list.InterestList;
		il.Init(ctx, session.GetActiveLang());
		pageInstance.Title = "List Interests"

		AdminContent := admin.Content{};
		AdminContent.Title = "Interests"
		AdminContent.SubTitle = "List Interests"

		AdminContent.Content = template.HTML(il.Render(il.GetToPage()))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true,pageInstance);
		return;
	}
}

func (ic *InterestController) EditAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(ic.AuthAction["edit"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var interestId = int64(id);
		var interest m.Interest;
		interest, err := interest.Get(interestId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "interest/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		var data map[string]interface{};
		if (!ctx.IsPost()) {
			data = map[string]interface{}{
				"id":         strconv.Itoa(int(interest.Id)),
				"name": interest.Name,
			};
		} else {
			data = map[string]interface{}{
				"id":         h.GetFormData(ctx, "id", false).(string),
				"name": h.GetFormData(ctx, "name", false).(string),
			};
		}

		var form = m.GetInterestForm(data, fmt.Sprintf("interest/edit/%v", data["id"].(string)));
		if (ctx.IsPost()) {
			succ, formErrors := ic.saveInterest(ctx, session, interest);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Interest save was successful.");
				Redirect(ctx, fmt.Sprintf("interest/edit/%v", data["id"].(string)), fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Interest - Edit"

		AdminContent := admin.Content{};
		AdminContent.Title = "Interest"
		AdminContent.SubTitle = fmt.Sprintf("Edit interest %v", interest.Name);
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "interest", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (ic *InterestController) saveInterest(ctx *fasthttp.RequestCtx, session *h.Session, interest m.Interest) (bool, map[string]error) {
	if (ctx.IsPost() && ((Ah.HasRights(ic.AuthAction["edit"],session) && interest.Id != 0) || (Ah.HasRights(ic.AuthAction["new"],session) && interest.Id == 0))) {
		var err error;
		var succ bool;
		var Validator = m.GetInterestFormValidator(ctx, interest);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		interest.Name = h.GetFormData(ctx, "name", false).(string);

		if (interest.Id > 0) {
			_, err = db.DbMap.Update(&interest);
		} else {
			err = db.DbMap.Insert(&interest);
		}
		h.Error(err, "", h.ERROR_LVL_ERROR)
		succ = err == nil;
		return succ, nil;
	} else {
		return false, nil;
	}
}

func (ic *InterestController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(ic.AuthAction["new"],session)) {
		var interest = m.NewEmptyInterest();
		var data map[string]interface{} = map[string]interface{}{};
		var dataKeys []string = []string{"id", "name"};
		for _,k := range dataKeys {
			var val string = "";
			if(ctx.IsPost()) {
				val = h.GetFormData(ctx, k, false).(string)
			}
			data[k] = val;
		}

		var form = m.GetInterestForm(data, "interest/new");
		if (ctx.IsPost()) {
			succ, formErrors := ic.saveInterest(ctx, session, interest);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Interest save was successful.");
				Redirect(ctx, "interest", fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Interest - New"

		AdminContent := admin.Content{};
		AdminContent.Title = "Interest"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (ic *InterestController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(ic.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var interestId = int64(id);
		var interest m.Interest;
		interest, err := interest.Get(interestId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "interest/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		interestName := interest.Name;
		count,err := db.DbMap.Delete(&interest);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil || count == 0){
			session.AddError(err.Error());
			session.AddError("An error occurred, could not delete interest.");
			status = fasthttp.StatusBadRequest;
		} else {
			session.AddSuccess(fmt.Sprintf("Interest %v has been deleted",interestName));
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "interest/index", status, true,pageInstance)
}
