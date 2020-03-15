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

type OwnerController struct {
	AuthAction map[string][]string;
}

func (oc *OwnerController) Init() {
	oc.AuthAction = make(map[string][]string);
	oc.AuthAction["edit"] = []string{"owner/edit"};
	oc.AuthAction["save"] = []string{"owner/edit", "owner/new"};
	oc.AuthAction["new"] = []string{"owner/new"};
	oc.AuthAction["list"] = []string{"owner/list"};
}

func (oc *OwnerController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(oc.AuthAction["list"],session)) {
		var ol list.OwnerList;
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

func (oc *OwnerController) EditAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(oc.AuthAction["edit"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var ownerId = int64(id);
		var owner m.Owner;
		owner, err := owner.Get(ownerId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "owner/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		var data map[string]interface{};
		if (!ctx.IsPost()) {
			var hun string = "0";
			if(owner.Hungarian){
				hun = "1";
			}
			data = map[string]interface{}{
				"id":         strconv.Itoa(int(owner.Id)),
				"name": owner.Name,
				"hungarian": hun,
			};
		} else {
			data = map[string]interface{}{
				"id":         h.GetFormData(ctx, "id", false).(string),
				"name": h.GetFormData(ctx, "name", false).(string),
				"hungarian": h.GetFormData(ctx, "hungarian", false).(string),
			};
		}

		var form = m.GetOwnerForm(data, fmt.Sprintf("owner/edit/%v", data["id"].(string)));
		if (ctx.IsPost()) {
			succ, formErrors := oc.saveOwner(ctx, session, owner);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Owner save was successful.");
				Redirect(ctx, fmt.Sprintf("owner/edit/%v", data["id"].(string)), fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Owner - Edit"

		AdminContent := admin.Content{};
		AdminContent.Title = "Owner"
		AdminContent.SubTitle = fmt.Sprintf("Edit owner %v", owner.Name);
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "owner", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (oc *OwnerController) saveOwner(ctx *fasthttp.RequestCtx, session *h.Session, owner m.Owner) (bool, map[string]error) {
	if (ctx.IsPost() && ((Ah.HasRights(oc.AuthAction["edit"],session) && owner.Id != 0) || (Ah.HasRights(oc.AuthAction["new"],session) && owner.Id == 0))) {
		var err error;
		var succ bool;
		var Validator = m.GetOwnerFormValidator(ctx, owner);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		owner.Name = h.GetFormData(ctx, "name", false).(string);
		owner.Hungarian = h.GetFormData(ctx, "hungarian", false).(string) == "1";

		if (owner.Id > 0) {
			_, err = db.DbMap.Update(&owner);
		} else {
			err = db.DbMap.Insert(&owner);
		}
		h.Error(err, "", h.ERROR_LVL_ERROR)
		succ = err == nil;
		return succ, nil;
	} else {
		return false, nil;
	}
}

func (oc *OwnerController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(oc.AuthAction["new"],session)) {
		var owner = m.NewEmptyOwner();
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

		var form = m.GetOwnerForm(data, "owner/new");
		if (ctx.IsPost()) {
			succ, formErrors := oc.saveOwner(ctx, session, owner);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("Owner save was successful.");
				Redirect(ctx, "owner", fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "Owner - New"

		AdminContent := admin.Content{};
		AdminContent.Title = "Owner"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope), "", nil, false, 0)
	} else {
		Redirect(ctx, "", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (oc *OwnerController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(oc.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3, ""));
		var ownerId = int64(id);
		var owner m.Owner;
		owner, err := owner.Get(ownerId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "owner/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		ownerName := owner.Name;
		count,err := db.DbMap.Delete(&owner);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil || count == 0){
			session.AddError(err.Error());
			session.AddError("An error occurred, could not delete owner.");
			status = fasthttp.StatusBadRequest;
		} else {
			session.AddSuccess(fmt.Sprintf("Owner %v has been deleted",ownerName));
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "owner/index", status, true,pageInstance)
}
