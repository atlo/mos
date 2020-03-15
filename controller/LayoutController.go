package controller

import (
	"github.com/valyala/fasthttp"
	h "mos/helper"
	"strings"
	adminview "mos/model/view/admin"
	"fmt"
	"mos/model/view"
)

type LayoutController struct {
	AuthAction map[string][]string;
	Layout     string;
}

func (l *LayoutController) Init() {
	l.AuthAction = make(map[string][]string);
}

func (l *LayoutController) PrependAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page, routeMap map[string]interface{}) {
	pageInstance.ContentType = "text/html; charset=utf-8;"
	pageInstance.Layout = "main.html";
	pageInstance.Scope = "frontend";
	if (strings.Index(string(ctx.Path()), "/"+h.GetConfig().AdminRouter) == 0) {
		pageInstance.Scope = "admin"
	}

	if (pageInstance.Scope == "admin") {
		l.prepareAdminView(ctx, session, pageInstance,routeMap);
	} else {
		l.prepareFrontendView(ctx, session, pageInstance, routeMap);
	}
}

func (l *LayoutController) prepareAdminView(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page, routeMap map[string]interface{}) {
	skipHeader, ok := routeMap["skip_header"];
	if (!ok) {
		skipHeader = false;
	}
	isAjax, ok := routeMap["is_ajax"];
	if (!ok) {
		isAjax = false;
	}

	h.PrintlnIf("template scope is admin", h.GetConfig().Mode.Debug)

	if (!isAjax.(bool) && !skipHeader.(bool)) {
		pageInstance.AddAdminScripts();
		pageInstance.AddAdminStylesheets();
		pageInstance.AddDefaultMetaData();
		pageInstance.AddContent(h.GetScopeTemplateString("layout/menu.html", h.GetMenu(session), pageInstance.Scope), "div", map[string]string{"id": "wrapper"}, false, 0)
		pageInstance.AddContent(h.GetScopeTemplateString("layout/messages.html", adminview.Messages{session.GetErrors(), session.GetSuccesses()}, pageInstance.Scope), "div", map[string]string{"id": "page-wrapper"}, false, 0)
	}
}

func (l *LayoutController) prepareFrontendView(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page, routeMap map[string]interface{}) {
	skipHeader, ok := routeMap["skip_header"];
	if (!ok) {
		skipHeader = false;
	}
	isAjax, ok := routeMap["is_ajax"];
	if (!ok) {
		isAjax = false;
	}

	pageInstance.AddDefaultMetaData();
	pageInstance.AddOgMetaData();
	if (!isAjax.(bool)) {
		//ide jönnek css-ek js-ek
		pageInstance.AddCss("/assets/css/main.css");
		pageInstance.BodyClass = "media-owner";
		pageInstance.AddScript("/assets/js/main.bundle.js","footer", "");
	}

	if (!isAjax.(bool) && !skipHeader.(bool)) {
		var cacheHeaderKeys []string = []string{session.GetActiveLang(),"header"};
		var filePath string = "layout/header.html";
		hasContent,headerContent := h.CacheStorage.GetString(filePath,cacheHeaderKeys);
		if(!hasContent){
			var header view.Header;
			header.Init(session);
			headerContent = h.GetScopeTemplateString(filePath,header,"frontend");
		}
		pageInstance.AddContent(headerContent,"div",map[string]string{"class":"page"},false,0);
	}
}

func (l *LayoutController) RenderAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page, routeMap map[string]interface{}) {
	l.PostpendAction(ctx, session, pageInstance, routeMap);
	h.PrintlnIf(fmt.Sprintf("Render layout %v", pageInstance.Scope), h.GetConfig().Mode.Debug)

	ctx.SetContentType(pageInstance.ContentType);
	independent,ok := routeMap["independent"];
	if(!ok || !independent.(bool)){
		_, err := ctx.WriteString(h.GetScopeTemplateString(fmt.Sprintf("layout/%v", pageInstance.Layout), pageInstance, pageInstance.Scope));
		h.Error(err, "", h.ERROR_LVL_ERROR)
	}
	return;
}

func (l *LayoutController) PostpendAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page, routeMap map[string]interface{}) {
	session.ClearMessages();
}
