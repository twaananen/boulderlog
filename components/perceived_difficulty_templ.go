// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
)

func PerceivedDifficulty(grade string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex flex-col items-center relative\"><button hx-get=\"/log/grade\" hx-target=\"#main-content\" class=\"mb-4 bg-gray-300 hover:bg-gray-400 dark:bg-gray-700 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200 font-bold py-2 px-4 rounded inline-flex items-center\"><i class=\"fas fa-arrow-left mr-2\"></i> Back</button><h2 class=\"text-2xl font-bold mb-4 text-gray-900 dark:text-gray-100\">Boulder Grade: ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(grade)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/perceived_difficulty.templ`, Line: 16, Col: 93}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2><div class=\"space-y-4 w-full max-w-md\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for i := 1; i <= 8; i++ {
			var templ_7745c5c3_Var3 = []any{
				"w-full py-3 px-6 rounded font-bold text-white",
				templ.KV("bg-green-600 dark:bg-green-700", i == 1),
				templ.KV("bg-green-500 dark:bg-green-600", i == 2),
				templ.KV("bg-green-400 dark:bg-green-500", i == 3),
				templ.KV("bg-green-300 dark:bg-green-400", i == 4),
				templ.KV("bg-yellow-400 dark:bg-yellow-500", i == 5),
				templ.KV("bg-yellow-500 dark:bg-yellow-600", i == 6),
				templ.KV("bg-orange-500 dark:bg-orange-600", i == 7),
				templ.KV("bg-red-500 dark:bg-red-600", i == 8),
			}
			templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var3...)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button hx-post=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/log/submit/%s/%d", grade, i))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/perceived_difficulty.templ`, Line: 20, Col: 58}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#main-content\" class=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(templ.CSSClasses(templ_7745c5c3_Var3).String())
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/perceived_difficulty.templ`, Line: 1, Col: 0}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">Difficulty ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(i))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/perceived_difficulty.templ`, Line: 34, Col: 32}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><button class=\"absolute top-0 right-0 bg-blue-500 dark:bg-blue-600 text-white rounded-full w-8 h-8 flex items-center justify-center\" onclick=\"toggleInfoPopup()\"><i class=\"fas fa-info\"></i></button><div id=\"infoPopupOverlay\" class=\"fixed inset-0 bg-black bg-opacity-50 hidden\" onclick=\"toggleInfoPopup()\"></div><div id=\"infoPopup\" class=\"hidden fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 p-4 rounded shadow-lg w-3/4 max-w-xs z-10\"><h3 class=\"font-bold mb-2 text-gray-900 dark:text-gray-100\">Perceived Difficulty Scale</h3><ul class=\"list-disc pl-5 text-gray-700 dark:text-gray-300\"><li>1 - (Topped) flash very easily</li><li>2 - (Topped) flash with some difficulty</li><li>3 - (Topped) flash with a lot of difficulty</li><li>4 - (Topped) topped with a lot of difficulty and multiple attempts</li><li>5 - (Not topped) very close to topping, maybe next time</li><li>6 - (Not topped) could do all moves separately but not together</li><li>7 - (Not topped) some moves just could not be done</li><li>8 - (Not topped) couldn't do any moves</li></ul></div></div><script>\n\t\tfunction toggleInfoPopup() {\n\t\t\tconst popup = document.getElementById('infoPopup');\n\t\t\tconst overlay = document.getElementById('infoPopupOverlay');\n\t\t\tpopup.classList.toggle('hidden');\n\t\t\toverlay.classList.toggle('hidden');\n\t\t}\n\t</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
