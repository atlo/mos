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

type MediaOwnerController struct {
	AuthAction map[string][]string;
}

func (mo *MediaOwnerController) Init() {
	mo.AuthAction = make(map[string][]string);
	mo.AuthAction["edit"] = []string{"media/edit"};
	mo.AuthAction["delete"] = []string{"media/edit"};
	mo.AuthAction["save"] = []string{"media/edit", "media/new"};
	mo.AuthAction["new"] = []string{"media/new"};
	mo.AuthAction["list"] = []string{"media/list"};
}

func (mo *MediaOwnerController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(mo.AuthAction["list"], session)) {
		var mol list.MediaOwnerList;
		mol.Init(ctx, session.GetActiveLang());
		pageInstance.Title = "List Media Owners"

		AdminContent := admin.Content{};
		AdminContent.Title = "Media Owners"
		AdminContent.SubTitle = "List Media Owners"

		AdminContent.Content = template.HTML(mol.Render(mol.GetToPage()))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true, pageInstance);
		return;
	}
}

func (mo *MediaOwnerController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(mo.AuthAction["new"], session)) {
		var mediaOwner = m.NewEmptyMediaOwner();
		var data map[string]interface{} = map[string]interface{}{};
		var dataKeys []string = []string{"id","media_id", "year","owner_id"};
		for _, k := range dataKeys {
			var val string = "";
			if (ctx.IsPost()) {
				val = h.GetFormData(ctx, k, false).(string)
			}
			data[k] = val;
		}

		var form = m.GetMediaOwnerForm(true, data, "mediaowner/new");
		if (ctx.IsPost()) {
			succ, formErrors := mo.saveMediaOwner(ctx, session, mediaOwner);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Media Owner save was successful.");
				Redirect(ctx, "mediaowner", fasthttp.StatusOK, true, pageInstance);
				return;
			}
		}

		pageInstance.Title = "Media Owner - New"

		AdminContent := admin.Content{};
		AdminContent.Title = "Media Owner"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "", fasthttp.StatusForbidden, true, pageInstance)
		return;
	}
}

func (mo *MediaOwnerController) saveMediaOwner(ctx *fasthttp.RequestCtx, session *h.Session, mediaOwner m.MediaOwner) (bool, map[string]error) {
	if (ctx.IsPost() && Ah.HasRights(mo.AuthAction["new"], session)) {
		var err error;
		var succ bool;

		var Validator = m.GetMediaOwnerFormValidator(ctx, mediaOwner);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		mediaId,err := strconv.Atoi(h.GetFormData(ctx, "media_id", false).(string));
		h.Error(err,"",h.ERROR_LVL_WARNING);

		ownerId,err := strconv.Atoi(h.GetFormData(ctx, "owner_id", false).(string));
		h.Error(err,"",h.ERROR_LVL_WARNING);

		year,err := strconv.Atoi(h.GetFormData(ctx, "year", false).(string));
		h.Error(err,"",h.ERROR_LVL_WARNING);

		mediaOwner.MediaId = int64(mediaId);
		mediaOwner.OwnerId = int64(ownerId);
		mediaOwner.Year = year;

		err = db.DbMap.Insert(&mediaOwner);
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

func (mo *MediaOwnerController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(mo.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var moId = int64(id);
		var mediaOwner m.MediaOwner;
		mediaOwner, err := mediaOwner.Get(moId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "mediaowner/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		count,err := db.DbMap.Delete(&mediaOwner);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil){
			session.AddError("An error occured, could not delete media owner connection.");
			status = fasthttp.StatusBadRequest;
		} else if(count == 1) {
			session.AddSuccess("Media Owner connection has been deleted");
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "mediaowner/index", status, true,pageInstance)
}
