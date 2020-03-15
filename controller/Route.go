package controller

import (
	"regexp"
	"strings"
	"sort"
	h "mos/helper"
	"github.com/valyala/fasthttp"
	"fmt"
	"mos/model/view"
)

var Routes []map[string]map[string]interface{}

var AccessC AccessController;
var UserC UserController;
var PageC PageController;
var ConfigC ConfigController;
var BlockC BlockController;
var MediaTypeC MediaTypeController;
var MediaOwnerC MediaOwnerController;
var MediaOperatorC MediaOperatorController;
var OperatorYearDataC OperatorYearDataController;
var MediaC MediaController;
var OwnerC OwnerController;
var InterestC InterestController;
var OperatorInterestC OperatorInterestController;
var OperatorC OperatorController;
var FeedC FeedController;
var LayoutC LayoutController;

var Ah h.AuthHelper;

func Redirect(ctx *fasthttp.RequestCtx, route string, status int, includeScope bool, page *view.Page) {
	page.Redirected = true;
	var url string;
	if(strings.Contains(route,"http://") || strings.Contains(route,"https://")){
		url = route;
	} else {
		url = h.GetUrl(route, nil, includeScope,page.Scope)
	}
	h.PrintlnIf(fmt.Sprintf("Redirecting to %v",url),h.GetConfig().Mode.Debug);
	ctx.Redirect(url, status);
}

func Route(ctx *fasthttp.RequestCtx) {
	var Log = h.SetLog();
	var p view.Page;

	var page *view.Page = p.Instantiates();
	var session = h.SessionGet(&ctx.Request.Header);
	var firstInRequest bool = true;
	var hadMach bool = false;
	var staticCompile = regexp.MustCompile("^/(vendor|assets|images|frontend)/?");
	var staticHandler = fasthttp.FSHandler("", 0);

	h.Lang.SetLanguage(ctx,session);
	// Set for local frontend development
	// ctx.Response.Header.Set("Access-Control-Allow-Origin", "*");

	if (staticCompile.MatchString(string(ctx.Path()))) {
		staticHandler(ctx);
	} else {
		if(AccessC.CheckBan(ctx,session,page)) {
			ctx.Response.Header.SetStatusCode(fasthttp.StatusTooManyRequests);
			ctx.WriteString("Too many request.");
			return;
		} else {
			for _, Route := range Routes {
				for strMatch, routeMap := range Route {
					var arrMatch = strings.Split(strMatch, "|");
					var mustC = regexp.MustCompile(arrMatch[1])
					var methods = strings.Split(arrMatch[0], ",");
					sort.Strings(methods)
					var methodI = sort.SearchStrings(methods, string(ctx.Method()));
					if (mustC.MatchString(string(ctx.Path())) && methodI < len(methods) && methods[methodI] == string(ctx.Method())) {
						hadMach = true;
						if (firstInRequest) {
							fmt.Println("--------------------------------------------------------------------------------------------------------------")
							h.PrintlnIf("Layout prepend running", h.GetConfig().Mode.Debug);
							LayoutC.PrependAction(ctx, session, page, routeMap);
							firstInRequest = false;
						}
						h.PrintlnIf(fmt.Sprintf("REGEXP FOUND: \"%v\" in path \"%v\"\n", arrMatch[1], string(ctx.Path())), h.GetConfig().Mode.Debug)
						ctx.SetStatusCode(fasthttp.StatusOK);
						funcToHandle := routeMap["func"].(func(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page));
						funcToHandle(ctx, session, page);
						if (!page.Redirected) {
							LayoutC.RenderAction(ctx, session, page, routeMap);
						} else {
							h.PrintlnIf(fmt.Sprintf("No render %v -> redirect", string(ctx.Path())), h.GetConfig().Mode.Debug);
							page.Redirected = false;
						}
						session.Send(&ctx.Response.Header, h.Duration);
					}
				}
				if (hadMach) {
					break;
				}
			}
		}
	}
	defer Log.Close();
}

func dispatchRoutes() {
	//just admin call without controller or action (login)
	AddRoute(fmt.Sprintf("GET|/%v/?$", h.GetConfig().AdminRouter),UserC.LoginAction,map[string]interface{}{});

	//Forbidden access default routes
	AddRoute(fmt.Sprintf("GET|/(%v/)?access/forbidden/?$", h.GetConfig().AdminRouter),AccessC.ForbiddenAction,map[string]interface{}{})

	//ADMIN REQUESTS
	adminDispatch();

	//FRONTEND REQUESTS
	frontendDispatch();
}

func adminDispatch(){
	emptyMap := map[string]interface{}{};
	//user login, logout, loginpost
	AddRoute(fmt.Sprintf("GET|^/%v/user/login$", h.GetConfig().AdminRouter), UserC.LoginAction,emptyMap)
	AddRoute(fmt.Sprintf("POST|^/%v/user/loginpost$", h.GetConfig().AdminRouter), UserC.LoginpostAction,emptyMap)
	AddRoute(fmt.Sprintf("GET|^/%v/user/welcome$", h.GetConfig().AdminRouter), UserC.WelcomeAction,emptyMap)
	AddRoute(fmt.Sprintf("POST|^/%v/user/logout$", h.GetConfig().AdminRouter), UserC.LogoutAction,emptyMap)

	//user useractions
	AddRoute(fmt.Sprintf("GET|^/%v/user/?(index)?$", h.GetConfig().AdminRouter),UserC.ListAction,emptyMap);
	AddRoute(fmt.Sprintf("GET,POST|^/%v/user/edit/(\\d)+$", h.GetConfig().AdminRouter),UserC.EditAction,emptyMap);
	AddRoute(fmt.Sprintf("GET|^/%v/user/delete/(\\d)+$", h.GetConfig().AdminRouter),UserC.DeleteAction,emptyMap);
	AddRoute(fmt.Sprintf("GET,POST|^/%v/user/new$", h.GetConfig().AdminRouter),UserC.NewAction,emptyMap);
	AddRoute(fmt.Sprintf("GET|^/%v/user/switchlanguage/([a-z])+$", h.GetConfig().AdminRouter),UserC.SwitchLanguageAction,emptyMap);

	//block useractions
	AddRoute(fmt.Sprintf("GET|^/%v/block/?(index)?$", h.GetConfig().AdminRouter), BlockC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/block/edit/(\\d)+$", h.GetConfig().AdminRouter), BlockC.EditAction, emptyMap)
	AddRoute(fmt.Sprintf("GET|^/%v/block/delete/(\\d)+$", h.GetConfig().AdminRouter), BlockC.DeleteAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/block/new$", h.GetConfig().AdminRouter), BlockC.NewAction, emptyMap)

	//media type useractions
	AddRoute(fmt.Sprintf("GET|^/%v/mediatype/?(index)?$", h.GetConfig().AdminRouter), MediaTypeC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/mediatype/edit/(\\d)+$", h.GetConfig().AdminRouter), MediaTypeC.EditAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/mediatype/new$", h.GetConfig().AdminRouter), MediaTypeC.NewAction, emptyMap)
	AddRoute(fmt.Sprintf("GET|^/%v/mediatype/delete/(\\d)+$", h.GetConfig().AdminRouter), MediaTypeC.DeleteAction, emptyMap)

	//media owner useractions
	AddRoute(fmt.Sprintf("GET|^/%v/mediaowner/?(index)?$", h.GetConfig().AdminRouter), MediaOwnerC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/mediaowner/new$", h.GetConfig().AdminRouter), MediaOwnerC.NewAction, emptyMap)
	AddRoute(fmt.Sprintf("GET|^/%v/mediaowner/delete/(\\d)+$", h.GetConfig().AdminRouter), MediaOwnerC.DeleteAction, emptyMap)

	//media operator useractions
	AddRoute(fmt.Sprintf("GET|^/%v/mediaoperator/?(index)?$", h.GetConfig().AdminRouter), MediaOperatorC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/mediaoperator/new$", h.GetConfig().AdminRouter), MediaOperatorC.NewAction, emptyMap)
	AddRoute(fmt.Sprintf("GET|^/%v/mediaoperator/delete/(\\d)+$", h.GetConfig().AdminRouter), MediaOperatorC.DeleteAction, emptyMap)

	//media useractions
	AddRoute(fmt.Sprintf("GET|^/%v/media/?(index)?$", h.GetConfig().AdminRouter), MediaC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/media/edit/(\\d)+$", h.GetConfig().AdminRouter), MediaC.EditAction, emptyMap)

	//owner useractions
	AddRoute(fmt.Sprintf("GET|^/%v/owner/?(index)?$", h.GetConfig().AdminRouter), OwnerC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/owner/edit/(\\d)+$", h.GetConfig().AdminRouter), OwnerC.EditAction, emptyMap)

	//interest useractions
	AddRoute(fmt.Sprintf("GET|^/%v/interest/?(index)?$", h.GetConfig().AdminRouter), InterestC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/interest/edit/(\\d)+$", h.GetConfig().AdminRouter), InterestC.EditAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/interest/new$", h.GetConfig().AdminRouter), InterestC.NewAction, emptyMap)
	AddRoute(fmt.Sprintf("GET|^/%v/interest/delete/(\\d)+$", h.GetConfig().AdminRouter), InterestC.DeleteAction, emptyMap)

	//operator useractions
	AddRoute(fmt.Sprintf("GET|^/%v/operator/?(index)?$", h.GetConfig().AdminRouter), OperatorC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/operator/edit/(\\d)+$", h.GetConfig().AdminRouter), OperatorC.EditAction, emptyMap)

	//operatorinterest useractions
	AddRoute(fmt.Sprintf("GET|^/%v/operatorinterest/?(index)?$", h.GetConfig().AdminRouter), OperatorInterestC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET|^/%v/operatorinterest/delete/(\\d)+$", h.GetConfig().AdminRouter), OperatorInterestC.DeleteAction, emptyMap)

	//operator yeardata useractions
	AddRoute(fmt.Sprintf("GET|^/%v/operatordata/?(index)?$", h.GetConfig().AdminRouter), OperatorYearDataC.ListAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/operatordata/new$", h.GetConfig().AdminRouter), OperatorYearDataC.NewAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/operatordata/edit/(\\d)+/(\\d{4})+$", h.GetConfig().AdminRouter), OperatorYearDataC.EditAction, emptyMap)
	AddRoute(fmt.Sprintf("GET,POST|^/%v/operatordata/delete/(\\d)+/(\\d{4})+$", h.GetConfig().AdminRouter), OperatorYearDataC.DeleteAction, emptyMap)

	//config useraction
	AddRoute(fmt.Sprintf("GET,POST|^/%v/config/?(index)?$", h.GetConfig().AdminRouter),ConfigC.IndexAction,emptyMap)
}

func frontendDispatch(){
	//AddRoute("GET|^/test$?",PageC.TestAction,map[string]interface{}{})
	//utolsó, ennek kell alul lennie, minden ide fut
	AddRoute("GET|^/media$",FeedC.MediaAction,map[string]interface{}{"skip_header":true,"is_ajax":true})
	AddRoute("GET|^/owners$",FeedC.OwnersAction,map[string]interface{}{"skip_header":true,"is_ajax":true})
	AddRoute("GET|^/interests$",FeedC.InterestsAction,map[string]interface{}{"skip_header":true,"is_ajax":true})
	AddRoute("GET|^/operators$",FeedC.OperatorsAction,map[string]interface{}{"skip_header":true,"is_ajax":true})
	AddRoute("GET|^/connections$",FeedC.ConnectionsAction,map[string]interface{}{"skip_header":true,"is_ajax":true})
	AddRoute("GET|^/?",PageC.IndexAction,map[string]interface{}{})
}

func AddRoute(path string, toCall func(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page), options map[string]interface{}){
	RouteOptions := map[string]interface{}{
		"func": toCall,
	};
	for k,v := range options{
		RouteOptions[k] = v;
	}
	Routes = append(Routes, map[string]map[string]interface{}{
		path: RouteOptions,
	});
}

func InitControllers() {
	UserC.Init();
	AccessC.Init();
	LayoutC.Init();
	PageC.Init();
	BlockC.Init();
	MediaTypeC.Init();
	MediaOwnerC.Init();
	MediaOperatorC.Init();
	OperatorYearDataC.Init();
	MediaC.Init();
	ConfigC.Init();
	OwnerC.Init();
	OperatorC.Init();
	FeedC.Init();
	InterestC.Init();
	OperatorInterestC.Init();

	dispatchRoutes();
}
