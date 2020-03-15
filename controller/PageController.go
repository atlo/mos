package controller

import (
	"github.com/valyala/fasthttp"
	h "mos/helper"
	"mos/model/view"
	"mos/model/view/page"
)

type PageController struct {
	AuthAction map[string][]string
}

func (p *PageController) Init() {
	p.AuthAction = make(map[string][]string)
	p.AuthAction["index"] = []string{"*"}
	p.AuthAction["test"] = []string{"*"}
}

func (p *PageController) IndexAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(p.AuthAction["index"], session) && pageInstance.Scope != "admin") {
		var template string = "page/index.html";
		hasContent,content := h.CacheStorage.GetString(template,[]string{"index",session.GetActiveLang()});
		if(!hasContent){
			var indexPage page.Index;
			indexPage.Init(session);
			content = h.GetScopeTemplateString(template,indexPage,"frontend");
		}
		pageInstance.AddContent(content,"",nil,false,0);
	} else {
		Redirect(ctx,"user/login",fasthttp.StatusForbidden,true,pageInstance);
	}
	return;
}

func (p *PageController) TestAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(p.AuthAction["test"], session) && pageInstance.Scope != "admin") {

	} else {
		Redirect(ctx,"user/login",fasthttp.StatusForbidden,true,pageInstance);
	}
	return;
}