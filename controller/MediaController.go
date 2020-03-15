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

type MediaController struct {
	AuthAction map[string][]string;
}

func (mc *MediaController) Init() {
	mc.AuthAction = make(map[string][]string);
	mc.AuthAction["edit"] = []string{"media/edit"};
	mc.AuthAction["save"] = []string{"media/edit", "media/new"};
	mc.AuthAction["new"] = []string{"media/new"};
	mc.AuthAction["list"] = []string{"media/list"};
}

func (mc *MediaController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(mc.AuthAction["list"],session)) {
		var ml list.MediaList;
		ml.Init(ctx, session.GetActiveLang());
		pageInstance.Title = "List Medias"

		AdminContent := admin.Content{};
		AdminContent.Title = "Medias"
		AdminContent.SubTitle = "List Medias"

		AdminContent.Content = template.HTML(ml.Render(ml.GetToPage()))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true,pageInstance);
		return;
	}
}

func (mc *MediaController) EditAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(mc.AuthAction["edit"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var mediaId = int64(id);
		var media m.Media;
		media, err := media.Get(mediaId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "media/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		var data map[string]interface{};
		if (!ctx.IsPost()) {
			var news string = "0";
			if(media.News){
				news = "1";
			}
			data = map[string]interface{}{
				"id":         strconv.Itoa(int(media.Id)),
				"name": media.Name,
				"news": news,
				"media_type_id": strconv.Itoa(int(media.MediaTypeId)),
			};
		} else {
			data = map[string]interface{}{
				"id":         h.GetFormData(ctx, "id", false).(string),
				"name":         h.GetFormData(ctx, "name", false).(string),
				"news": h.GetFormData(ctx, "news", false).(string),
				"media_type_id": h.GetFormData(ctx, "media_type_id", false).(string),
			};
		}

		var form = m.GetMediaForm(data, fmt.Sprintf("media/edit/%v", data["id"].(string)));
		if (ctx.IsPost()) {
			succ, formErrors := mc.saveMedia(ctx, session, media);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Media save was successful.");
				Redirect(ctx, fmt.Sprintf("media/edit/%v", data["id"].(string)), fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Media - Edit"

		AdminContent := admin.Content{};
		AdminContent.Title = "Media"
		AdminContent.SubTitle = fmt.Sprintf("Edit media %v", media.Name);
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "media", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (mc *MediaController) saveMedia(ctx *fasthttp.RequestCtx, session *h.Session, media m.Media) (bool, map[string]error) {
	if (ctx.IsPost() && ((Ah.HasRights(mc.AuthAction["edit"],session) && media.Id != 0) || (Ah.HasRights(mc.AuthAction["new"],session) && media.Id == 0))) {
		var err error;
		var succ bool;
		var Validator = m.GetMediaFormValidator(ctx, media);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		media.Name = h.GetFormData(ctx, "name", false).(string);
		media.News = h.GetFormData(ctx, "news", false).(string) == "1";
		mediaTypeId,_ := strconv.Atoi(h.GetFormData(ctx, "media_type_id", false).(string));
		media.MediaTypeId = int64(mediaTypeId);

		if (media.Id > 0) {
			_, err = db.DbMap.Update(&media);
		} else {
			err = db.DbMap.Insert(&media);
		}
		h.Error(err, "", h.ERROR_LVL_ERROR)
		succ = err == nil;
		return succ, nil;
	} else {
		return false, nil;
	}
}

func (mc *MediaController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(mc.AuthAction["new"],session)) {
		var media = m.NewEmptyMedia();
		var data map[string]interface{} = map[string]interface{}{};
		var dataKeys []string = []string{"id", "name"};
		for _,k := range dataKeys {
			var val string = "";
			if(ctx.IsPost()) {
				val = h.GetFormData(ctx, k, false).(string)
			}
			data[k] = val;
		}
		data["hungarian"] = h.GetFormData(ctx, "hungarian", true).([]string);

		var form = m.GetMediaForm(data, "media/new");
		if (ctx.IsPost()) {
			succ, formErrors := mc.saveMedia(ctx, session, media);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Media save was successful.");
				Redirect(ctx, "media", fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Media - New"

		AdminContent := admin.Content{};
		AdminContent.Title = "Media"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (mc *MediaController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(mc.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var mediaId = int64(id);
		var media m.Media;
		media, err := media.Get(mediaId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "media/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		mediaName := media.Name;
		count,err := db.DbMap.Delete(&media);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil){
			session.AddError(err.Error());
			session.AddError("An error occured, could not delete the media.");
			status = fasthttp.StatusBadRequest;
		} else if(count == 1) {
			session.AddSuccess(fmt.Sprintf("Media %v has been deleted",mediaName));
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "media/index", status, true,pageInstance)
}
