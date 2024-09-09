// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Login() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Login - BoulderLog</title><link href=\"/static/css/tailwind.css\" rel=\"stylesheet\"><script src=\"/static/js/htmx.min.js\"></script></head><body class=\"bg-gray-100\"><div class=\"min-h-screen flex items-center justify-center\"><div class=\"bg-white p-8 rounded-lg shadow-md w-full max-w-md\"><h2 class=\"text-2xl font-bold mb-6 text-center\">Login to BoulderLog</h2><form hx-post=\"/auth/login\" hx-redirect=\"/\"><div class=\"mb-4\"><label for=\"username\" class=\"block text-gray-700 text-sm font-bold mb-2\">Username</label> <input type=\"text\" id=\"username\" name=\"username\" class=\"w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500\" required></div><div class=\"mb-6\"><label for=\"password\" class=\"block text-gray-700 text-sm font-bold mb-2\">Password</label> <input type=\"password\" id=\"password\" name=\"password\" class=\"w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500\" required></div><button type=\"submit\" class=\"w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50\">Login</button></form><div class=\"mt-4 text-center\"><a href=\"/\" class=\"text-blue-500 hover:text-blue-600\">Back to Home</a></div></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate