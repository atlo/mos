package controller

import (
	"github.com/valyala/fasthttp"
	"html/template"
	h "mos/helper"
	"mos/model/list"
	"mos/model/view/admin"
	"mos/model/view"
	"strconv"
	"fmt"
	"mos/model"
	"mos/db"
)

type OperatorInterestController struct {
	AuthAction map[string][]string;
}

func (oic *OperatorInterestController) Init() {
	oic.AuthAction = make(map[string][]string);
	oic.AuthAction["list"] = []string{"interest/list"};
	oic.AuthAction["delete"] = []string{"interest/delete"};
}

func (oic *OperatorInterestController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(oic.AuthAction["list"],session)) {
		var oil list.OperatorInterestList;
		oil.Init(ctx, session.GetActiveLang());
		pageInstance.Title = "List Operator Interests"

		AdminContent := admin.Content{};
		AdminContent.Title = "Operator Interests"
		AdminContent.SubTitle = "List Operator Interests"

		AdminContent.Content = template.HTML(oil.Render(oil.GetToPage()))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true,pageInstance);
		return;
	}
}

func (oc *OperatorInterestController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(oc.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var operatorInterestId = int64(id);
		var operatorInterest model.OperatorInterest;
		operatorInterest, err := operatorInterest.Get(operatorInterestId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "interest/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		count,err := db.DbMap.Delete(&operatorInterest);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil || count == 0){
			session.AddError("An error occurred, could not delete operator interest connection.");
			status = fasthttp.StatusBadRequest;
			return;
		} else {
			session.AddSuccess(fmt.Sprintf("Operator interest connection has been deleted."));
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "operatorinterest/index", status, true,pageInstance)
}
