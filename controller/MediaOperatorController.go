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
)

type MediaOperatorController struct {
	AuthAction map[string][]string;
}

func (mo *MediaOperatorController) Init() {
	mo.AuthAction = make(map[string][]string);
	mo.AuthAction["edit"] = []string{"media/edit"};
	mo.AuthAction["delete"] = []string{"media/edit"};
	mo.AuthAction["save"] = []string{"media/edit", "media/new"};
	mo.AuthAction["new"] = []string{"media/new"};
	mo.AuthAction["list"] = []string{"media/list"};
}

func (mo *MediaOperatorController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(mo.AuthAction["list"], session)) {
		var mol list.MediaOperatorList;
		mol.Init(ctx, session.GetActiveLang());
		pageInstance.Title = "List Media Operators"

		AdminContent := admin.Content{};
		AdminContent.Title = "Media Operators"
		AdminContent.SubTitle = "List Media Operators"

		AdminContent.Content = template.HTML(mol.Render(mol.GetToPage()))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true, pageInstance);
		return;
	}
}

func (mo *MediaOperatorController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(mo.AuthAction["new"], session)) {
		var mediaOperator = m.NewEmptyMediaOperator();
		var data map[string]interface{} = map[string]interface{}{};
		var dataKeys []string = []string{"id","media_id", "year","operator_id"};
		for _, k := range dataKeys {
			var val string = "";
			if (ctx.IsPost()) {
				val = h.GetFormData(ctx, k, false).(string)
			}
			data[k] = val;
		}

		var form = m.GetMediaOperatorForm(true, data, "mediaoperator/new");
		if (ctx.IsPost()) {
			succ, formErrors := mo.saveMediaOperator(ctx, session, mediaOperator);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Media Operator save was successful.");
				Redirect(ctx, "mediaoperator", fasthttp.StatusOK, true, pageInstance);
				return;
			}
		}

		pageInstance.Title = "Media Operator - New"

		AdminContent := admin.Content{};
		AdminContent.Title = "Media Operator"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "", fasthttp.StatusForbidden, true, pageInstance)
		return;
	}
}

func (mo *MediaOperatorController) saveMediaOperator(ctx *fasthttp.RequestCtx, session *h.Session, mediaOperator m.MediaOperator) (bool, map[string]error) {
	if (ctx.IsPost() && Ah.HasRights(mo.AuthAction["new"], session)) {
		var err error;
		var succ bool;

		var Validator = m.GetMediaOperatorFormValidator(ctx, mediaOperator);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		mediaId,err := strconv.Atoi(h.GetFormData(ctx, "media_id", false).(string));
		h.Error(err,"",h.ERROR_LVL_WARNING);

		operatorId,err := strconv.Atoi(h.GetFormData(ctx, "operator_id", false).(string));
		h.Error(err,"",h.ERROR_LVL_WARNING);

		year,err := strconv.Atoi(h.GetFormData(ctx, "year", false).(string));
		h.Error(err,"",h.ERROR_LVL_WARNING);

		mediaOperator.MediaId = int64(mediaId);
		mediaOperator.OperatorId = int64(operatorId);
		mediaOperator.Year = year;

		err = db.DbMap.Insert(&mediaOperator);
		h.Error(err, "", h.ERROR_LVL_ERROR);

		succ = err == nil;
		var errRet map[string]error = nil;
		if(!succ){
			errRet = map[string]error{"media_id":err};
		}
		return succ, errRet;
	} else {
		return false, nil;
	}
}

func (mo *MediaOperatorController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(mo.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var moId = int64(id);
		var mediaOperator m.MediaOperator;
		mediaOperator, err := mediaOperator.Get(moId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "mediaoperator/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		count,err := db.DbMap.Delete(&mediaOperator);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil){
			session.AddError(err.Error());
			session.AddError("Could not delete media operator connection.");
			Redirect(ctx,"mediaoperator/index",fasthttp.StatusBadRequest,true,pageInstance);
			return;
		} else if(count == 1) {
			session.AddSuccess("Media Operator connection has been deleted");
			Redirect(ctx,"mediaoperator/index",fasthttp.StatusOK,true,pageInstance);
			return;
		}
	} else {
		Redirect(ctx, "mediaoperator/index", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}
