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

type MediaTypeController struct {
	AuthAction map[string][]string;
}

func (mt *MediaTypeController) Init() {
	mt.AuthAction = make(map[string][]string);
	mt.AuthAction["edit"] = []string{"mediatype/edit"};
	mt.AuthAction["save"] = []string{"mediatype/edit", "mediatype/new"};
	mt.AuthAction["new"] = []string{"mediatype/new"};
	mt.AuthAction["list"] = []string{"mediatype/list"};
	mt.AuthAction["delete"] = []string{"mediatype/delete"};
}

func (mt *MediaTypeController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (!Ah.HasRights(mt.AuthAction["list"], session)) {
		session.AddError("You do not have the required permissions to execute the request.");
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true, pageInstance)
		return;
	};
	var mtl list.MediaTypeList;
	mtl.Init(ctx, session.GetActiveLang());
	pageInstance.Title = "List Media Types"

	AdminContent := admin.Content{};
	AdminContent.Title = "Media Types"
	AdminContent.SubTitle = "List Media Types"

	AdminContent.Content = template.HTML(mtl.Render(mtl.GetToPage()))
	pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
}

func (mt *MediaTypeController) EditAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (!Ah.HasRights(mt.AuthAction["edit"], session)) {
		session.AddError("You do not have the required permissions to execute the request.");
		Redirect(ctx, "mediatype/index", fasthttp.StatusForbidden, true, pageInstance)
		return;
	};
	//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
	var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
	var mediaTypeId = int64(id);
	var mediaType m.MediaType;
	mediaType, err := mediaType.Get(mediaTypeId);
	if (err != nil) {
		session.AddError(err.Error());
		h.Error(err, "", h.ERROR_LVL_WARNING);
		Redirect(ctx, "mediatype/index", fasthttp.StatusOK, true, pageInstance);
		return;
	}

	var data map[string]interface{};
	if (!ctx.IsPost()) {
		data = map[string]interface{}{
			"id":   strconv.Itoa(int(mediaType.Id)),
			"name": mediaType.Name,
		};
	} else {
		data = map[string]interface{}{
			"id":   h.GetFormData(ctx, "id", false).(string),
			"name": h.GetFormData(ctx, "name", false).(string),
		};
	}

	var form = m.GetMediaTypeForm(data, fmt.Sprintf("mediatype/edit/%v", data["id"].(string)));
	if (ctx.IsPost()) {
		succ, formErrors := mt.saveMediaType(ctx, session, mediaType);
		form.SetErrors(formErrors);
		if (succ) {
			session.AddSuccess("Media Type save was successful.");
			Redirect(ctx, fmt.Sprintf("mediatype/edit/%v", data["id"].(string)), fasthttp.StatusOK, true, pageInstance);
			return;
		}
	}

	pageInstance.Title = "Media Type - Edit"

	AdminContent := admin.Content{};
	AdminContent.Title = "Media Type"
	AdminContent.SubTitle = fmt.Sprintf("Edit media type %v", mediaType.Name);
	AdminContent.Content = template.HTML(form.Render())
	pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
}

func (mt *MediaTypeController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (!Ah.HasRights(mt.AuthAction["new"], session)) {
		session.AddError("You do not have the required permissions to execute the request.");
		Redirect(ctx, "mediatype/index", fasthttp.StatusForbidden, true, pageInstance)
		return;
	};
	var mediaType = m.NewEmptyMediaType();
	var data map[string]interface{} = map[string]interface{}{};
	var dataKeys []string = []string{"id", "name"};
	for _, k := range dataKeys {
		var val string = "";
		if (ctx.IsPost()) {
			val = h.GetFormData(ctx, k, false).(string)
		}
		data[k] = val;
	}

	var form = m.GetMediaTypeForm(data, "mediatype/new");
	if (ctx.IsPost()) {
		succ, formErrors := mt.saveMediaType(ctx, session, mediaType);
		form.SetErrors(formErrors);
		if (succ) {
			session.AddSuccess("Media Type save was successful.");
			Redirect(ctx, "mediatype", fasthttp.StatusOK, true, pageInstance);
			return;
		}
	}

	pageInstance.Title = "Media Type - New"

	AdminContent := admin.Content{};
	AdminContent.Title = "Media Type"
	AdminContent.SubTitle = "New";
	AdminContent.Content = template.HTML(form.Render())
	pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent, pageInstance.Scope), "", nil, false, 0)
}

func (mt *MediaTypeController) saveMediaType(ctx *fasthttp.RequestCtx, session *h.Session, mediaType m.MediaType) (bool, map[string]error) {
	if (ctx.IsPost() && ((Ah.HasRights(mt.AuthAction["edit"], session) && mediaType.Id != 0) || (Ah.HasRights(mt.AuthAction["new"], session) && mediaType.Id == 0))) {
		var err error;
		var succ bool;
		var Validator = m.GetMediaTypeFormValidator(ctx, mediaType);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		mediaType.Name = h.GetFormData(ctx, "name", false).(string);
		if (mediaType.Id > 0) {
			_, err = db.DbMap.Update(&mediaType);
		} else {
			err = db.DbMap.Insert(&mediaType);
		}
		h.Error(err, "", h.ERROR_LVL_ERROR)
		succ = err == nil;
		return succ, nil;
	} else {
		return false, nil;
	}
}

func (mt *MediaTypeController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) () {
	if (!Ah.HasRights(mt.AuthAction["delete"], session)) {
		session.AddError("You do not have the required permissions to execute the request.");
		Redirect(ctx, "", fasthttp.StatusForbidden, true, pageInstance)
		return;
	};
	var status int;
	var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
	var mediaTypeId = int64(id);
	var mediaType m.MediaType;
	mediaType, err := mediaType.Get(mediaTypeId);
	if (err != nil) {
		session.AddError(err.Error());
		h.Error(err, "", h.ERROR_LVL_WARNING);
		Redirect(ctx, "mediatype/index", fasthttp.StatusOK, true, pageInstance);
		return;
	}

	intDeleted, err := db.DbMap.Delete(&mediaType);
	h.Error(err, "", h.ERROR_LVL_WARNING);
	if (err != nil) {
		session.AddError(err.Error());
	} else if (intDeleted == 0) {
		session.AddError("An error occured, could not remove item");
		status = fasthttp.StatusBadRequest;
	} else {
		session.AddSuccess("The item has been successfully removed.");
		status = fasthttp.StatusOK;
	}
	Redirect(ctx, "mediatype/index", status, true, pageInstance);
}
