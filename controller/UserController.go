package controller

import (
	"github.com/valyala/fasthttp"
	"html/template"
	h "mos/helper"
	m "mos/model"
	"strconv"
	adminview "mos/model/view/admin"
	"fmt"
	"mos/db"
	"mos/model/list"
	"mos/model/view"
)

type UserController struct {
	AuthAction map[string][]string;
}

func (u *UserController) Init() {
	u.AuthAction = make(map[string][]string);
	u.AuthAction["login"] = []string{"!@"};
	u.AuthAction["loginpost"] = []string{"!@"};
	u.AuthAction["logout"] = []string{"@"};
	u.AuthAction["edit"] = []string{"user/edit"};
	u.AuthAction["delete"] = []string{"user/delete"};
	u.AuthAction["save"] = []string{"user/edit", "user/new"};
	u.AuthAction["new"] = []string{"user/new"};
	u.AuthAction["list"] = []string{"user/list"};
	u.AuthAction["welcome"] = []string{"@"};
	u.AuthAction["switchlanguage"] = []string{"@"};
}

func (u *UserController) LoginAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(u.AuthAction["login"],session)) {
		pageInstance.Title = "User - Log in"

		AdminContent := adminview.Content{};

		AdminContent.Title = "User"
		AdminContent.SubTitle = "Log in";
		AdminContent.Content = template.HTML(h.GetScopeTemplateString("user/login.html", adminview.LoginForm{
			"POST",
			h.GetUrl("user/loginpost", nil,true,pageInstance.Scope),
		},pageInstance.Scope));

		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope),"",nil, false,0);
	} else {
		Redirect(ctx, "user/welcome", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (u *UserController) SwitchLanguageAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(u.AuthAction["switchlanguage"],session)) {
		var lc string = h.GetParamFromCtxPath(ctx,3,h.DefLang);
		if(session.GetActiveLang() != lc && h.Lang.IsAvailable(lc)){
			session.SetActiveLang(lc);
			session.AddSuccess(fmt.Sprintf("Changed store view to `%v`",lc));
		} else {
			session.AddError(fmt.Sprintf("You are already on language `%v`",lc));
		}
	}
	Redirect(ctx, string(ctx.Request.Header.Referer()), fasthttp.StatusOK, true, pageInstance)
}

func (u *UserController) LoginpostAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(u.AuthAction["loginpost"],session)) {
		var email = h.FormValue(ctx,"email");
		var password = h.FormValue(ctx,"password");
		var user m.User;
		var err error;
		user, err = user.GetUser(email, password);
		if err != nil {
			session.AddError(err.Error())
			Redirect(ctx, "user/login", fasthttp.StatusOK, true,pageInstance)
			return;
		}

		h.PrintlnIf("Logging in", h.GetConfig().Mode.Debug);
		session.Login(user.Id, user.SuperAdmin, user.GetRoles());
		status = fasthttp.StatusOK;
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "user/welcome", status, true,pageInstance)
}

func (u *UserController) LogoutAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(u.AuthAction["logout"],session)) {
		h.PrintlnIf("Logged in -> logout", h.GetConfig().Mode.Debug && session.Value(h.USER_SESSION_LOGGEDIN_KEY).(bool))
		h.PrintlnIf("Not logged in -> access denied", h.GetConfig().Mode.Debug && !session.Value(h.USER_SESSION_LOGGEDIN_KEY).(bool))
		session.Logout();
		Redirect(ctx, "user/login", fasthttp.StatusOK, true,pageInstance)
		return;
	} else {
		Redirect(ctx, "user/login", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (u *UserController) ListAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(u.AuthAction["list"],session)) {
		var ul list.UserList;
		ul.Init(ctx, session.GetActiveLang());

		pageInstance.Title = "List Users"

		AdminContent := adminview.Content{};
		AdminContent.Title = "Users"
		AdminContent.SubTitle = "List Users"
		AdminContent.Content = template.HTML(ul.Render(ul.GetToPage()))

		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope),"",nil, false,0)
	} else {
		Redirect(ctx, "user/login", fasthttp.StatusForbidden, true,pageInstance);
		return;
	}
}

func (u *UserController) WelcomeAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(u.AuthAction["welcome"],session)) {
		pageInstance.Title = "User - Welcome"

		AdminContent := adminview.Content{};
		AdminContent.Title = "User"
		AdminContent.SubTitle = "Welcome";
		AdminContent.Content = template.HTML(h.GetScopeTemplateString("user/welcome.html", nil,pageInstance.Scope))
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope),"",nil, false,0)
	} else {
		Redirect(ctx, "user/login", fasthttp.StatusForbidden, true,pageInstance);
		return;
	}
}

func (u *UserController) EditAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(u.AuthAction["edit"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3,""));
		var userId = int64(id);
		if (session.GetUserId() == userId) {
			session.AddError("You can not edit your own user.");
			Redirect(ctx, "user/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		var user m.User;
		User, err := user.Get(userId);
		if (err != nil) {
			session.AddError(err.Error());
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "user/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		if (User.SuperAdmin && !session.IsSuperAdmin()) {
			session.AddError("Nice Try, but you are not a superadmin, to edit another. :)");
			Redirect(ctx, "user/index", fasthttp.StatusForbidden, true,pageInstance);
			return;
		}

		var data map[string]interface{};
		if (!ctx.IsPost()) {
			data = map[string]interface{}{
				"id":              strconv.Itoa(int(User.Id)),
				"email":           User.Email,
				"password":        "",
				"password_verify": "",
				"status_id":       []string{strconv.Itoa(int(User.StatusId))},
				"role":            User.GetRoles(),
			};
		} else {
			data = map[string]interface{}{
				"id":              h.GetFormData(ctx, "id", false).(string),
				"email":           h.GetFormData(ctx, "email", false).(string),
				"password":        h.GetFormData(ctx, "password", false).(string),
				"password_verify": h.GetFormData(ctx, "password_verify", false).(string),
				"status_id":       h.GetFormData(ctx, "status_id", true).([]string),
				"role":            h.GetFormData(ctx, "role", true).([]string),
			};
		}

		var form = m.GetUserForm(data, fmt.Sprintf("user/edit/%v", data["id"].(string)));
		if (ctx.IsPost()) {
			succ, formErrors := u.saveUser(ctx, session, User);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("User save was successful.");
				Redirect(ctx,fmt.Sprintf("user/edit/%v",data["id"].(string)),fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "User - Edit"

		AdminContent := adminview.Content{};
		AdminContent.Title = "User"
		AdminContent.SubTitle = fmt.Sprintf("Edit user %v", User.Email);
		AdminContent.Content = template.HTML(form.Render())

		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope),"",nil, false,0)
	} else {
		Redirect(ctx, "user/index", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (u *UserController) DeleteAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	var status int;
	if (Ah.HasRights(u.AuthAction["delete"],session)) {
		//azért nem kell vizsgálni az errort, mert a request reguláris kifejezése csak akkor hozza ide, ha a végén \d van :)
		var id, _ = strconv.Atoi(h.GetParamFromCtxPath(ctx, 3,""));
		var userId = int64(id);
		if (session.GetUserId() == userId) {
			session.AddError("You can not delete your own user.");
			Redirect(ctx, "user/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		var user m.User;
		User, err := user.Get(userId);
		if (err != nil) {
			h.Error(err, "", h.ERROR_LVL_WARNING);
			Redirect(ctx, "user/index", fasthttp.StatusOK, true,pageInstance);
			return;
		}

		if (User.SuperAdmin && !session.IsSuperAdmin()) {
			session.AddError("Nice Try, but you are not a superadmin, to delete another. :)");
			Redirect(ctx, "user/index", fasthttp.StatusForbidden, true,pageInstance);
			return;
		}

		emailAddress := User.Email;
		count,err := db.DbMap.Delete(&User);
		h.Error(err,"",h.ERROR_LVL_WARNING);
		if(err != nil || count == 0){
			session.AddError("An error occurred, could not delete user.");
			status = fasthttp.StatusBadRequest;
		} else if(count == 1) {
			session.AddSuccess(fmt.Sprintf("User %v has been deleted",emailAddress));
			status = fasthttp.StatusOK;
		}
	} else {
		status = fasthttp.StatusForbidden;
	}
	Redirect(ctx, "user/index", status, true,pageInstance)
}

func (u *UserController) NewAction(ctx *fasthttp.RequestCtx, session *h.Session, pageInstance *view.Page) {
	if (Ah.HasRights(u.AuthAction["new"],session)) {
		var User = m.NewEmptyUser();
		var data map[string]interface{};
		if (!ctx.IsPost()) {
			data = map[string]interface{}{
				"id":              "",
				"email":           "",
				"password":        "",
				"password_verify": "",
				"status_id":       []string{string(m.STATUS_DEFAULT_VALUE)},
				"role":            []string{},
			};
		} else {
			data = map[string]interface{}{
				"id":              h.GetFormData(ctx, "id", false).(string),
				"email":           h.GetFormData(ctx, "email", false).(string),
				"password":        h.GetFormData(ctx, "password", false).(string),
				"password_verify": h.GetFormData(ctx, "password_verify", false).(string),
				"status_id":       h.GetFormData(ctx, "status_id", true).([]string),
				"role":            h.GetFormData(ctx, "role", true).([]string),
			};
		}

		var form = m.GetUserForm(data, "user/new");
		if (ctx.IsPost()) {
			succ, formErrors := u.saveUser(ctx, session, User);
			form.SetErrors(formErrors);
			if (succ) {
				session.AddSuccess("User save was successful.");
				Redirect(ctx,"user/index", fasthttp.StatusOK, true,pageInstance);
				return;
			}
		}

		pageInstance.Title = "User - New"

		AdminContent := adminview.Content{};
		AdminContent.Title = "User"
		AdminContent.SubTitle = "New";
		AdminContent.Content = template.HTML(form.Render())
		pageInstance.AddContent(h.GetScopeTemplateString("layout/content.html", AdminContent,pageInstance.Scope),"",nil, false,0)
	} else {
		Redirect(ctx, "user/index", fasthttp.StatusForbidden, true,pageInstance)
		return;
	}
}

func (u *UserController) saveUser(ctx *fasthttp.RequestCtx, session *h.Session, User m.User) (bool, map[string]error) {
	if (ctx.IsPost() && ((Ah.HasRights(u.AuthAction["edit"],session) && User.Id!=0) || (Ah.HasRights(u.AuthAction["new"],session) && User.Id==0))) {
		var err error;
		var succ bool;
		var Validator = m.GetUserFormValidator(ctx, User);
		succ, errors := Validator.Validate();
		if (!succ) {
			return false, errors;
		}

		User.Email = h.GetFormData(ctx, "email", false).(string);
		statusId, err := strconv.Atoi(h.GetFormData(ctx, "status_id", false).(string))
		h.Error(err, "", h.ERROR_LVL_WARNING)
		User.StatusId = int64(statusId);
		if (h.GetFormData(ctx, "password", false).(string) != "") {
			User.Password = h.GetFormData(ctx, "password", false).(string);
		}

		if (User.Id > 0) {
			_, err = db.DbMap.Update(&User);
		} else {
			err = db.DbMap.Insert(&User);
		}

		succ = err == nil;
		h.Error(err, "", h.ERROR_LVL_ERROR)

		h.PrintlnIf("Save successful", h.GetConfig().Mode.Debug && succ)
		h.PrintlnIf("Unsuccessful save", h.GetConfig().Mode.Debug && !succ)

		User.ModifyRoles(h.GetFormData(ctx, "role", true).([]string));
		return succ, nil;
	} else {
		return false, nil;
	}
}
