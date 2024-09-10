// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Home() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>BoulderLog - Track Your Climbing Progress</title><link href=\"/static/css/tailwind.css\" rel=\"stylesheet\"><script src=\"/static/js/htmx.min.js\"></script></head><body class=\"bg-gray-100\"><header class=\"bg-blue-600 text-white p-4 fixed top-0 left-0 right-0 z-10\"><div class=\"container mx-auto flex justify-between items-center\"><h1 class=\"text-2xl font-bold\">BoulderLog</h1><div hx-get=\"/auth/status\" hx-trigger=\"load, every 5m\" hx-swap=\"outerHTML\"></div></div></header><div class=\"container mx-auto px-4 py-8 mt-16\"><h2 class=\"text-4xl font-bold text-center mb-4\">Elevate Your Climbing Game</h2><p class=\"text-xl text-center mb-8\">Track, Analyze, Conquer!</p><div class=\"text-center\"><a href=\"/login\" class=\"bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded\">Get Started</a></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
