// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "jaxon.app/jaxon/internal/budget/models/categories"
import "fmt"
import "jaxon.app/jaxon/internal/templates"

func BudgetCreate(categories []categories.Category) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div id=\"main-content\" hx-swap-oob=\"outerHTML:#main-content\" class=\"flex flex-col h-full overflow-y-auto py-3 px-3\"><div class=\"w-[900px] pt-16\"><div id=\"create-budget-breadcrumb\" class=\"hidden flex flex-row\"><svg xmlns=\"http://www.w3.org/2000/svg\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-6\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M6.75 15.75 3 12m0 0 3.75-3.75M3 12h18\"></path></svg> budgets</div><div class=\"flex justify-between mb-7\"><h1 class=\"font-medium text-2xl text-slate-300\">Create a new budget</h1></div><form action=\"/budgets\" method=\"post\"><div class=\"flex flex-col gap-y-5 w-full h-1/2 p-4\"><div class=\"flex flex-row justify-between mb-4\"><span class=\"text-base text-slate-400\">Create a budget for <span class=\"font-bold text-base text-teal-800\">Groceries</span></span> <svg viewBox=\"0 0 24 24\" class=\"size-4 stroke-slate-500 stroke-[2.5px] fill-current\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M6 18 18 6M6 6l12 12\"></path></svg></div><div class=\"w-full flex justify-between gap-y-2\"><div class=\"flex flex-col\"><span class=\"text-slate-300 text-sm font-semibold\">Category</span><div class=\"text-sm font-normal text-slate-500\">What is this budget for?</div></div><div class=\"flex items-end\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templates.Dropdown(categories, "cat").Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<select id=\"create-budget-category\" name=\"category\" class=\"dropdown\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, category := range categories {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "<option id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("cat_%s", category.Id))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/budget/templates/budget-create.templ`, Line: 36, Col: 56}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(category.Id)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/budget/templates/budget-create.templ`, Line: 36, Col: 78}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(category.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/budget/templates/budget-create.templ`, Line: 36, Col: 96}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</select></div></div><div class=\"w-full flex justify-between gap-y-2\"><div class=\"flex flex-col\"><span class=\"text-slate-300 text-sm font-semibold\">Amount</span><div class=\"text-sm font-normal text-slate-500\">How much do you plan to spend?</div></div><div class=\"flex items-end\"><span class=\"text-teal-700 me-2\">$</span> <input id=\"budget-amount-edit\" type=\"text\" name=\"amount\" class=\"input input-sm input-bordered\" placeholder=\"Amount\"></div></div><div class=\"hidden w-full flex justify-between gap-y-2\"><div class=\"flex flex-col\"><span class=\"text-slate-300 text-sm font-semibold\">Frequency</span><div class=\"text-sm font-normal text-slate-500\">How often do you plan to spend?</div></div><div class=\"flex items-end\"><div id=\"toggle-count\" class=\"p-0.5 inline-block bg-gray-100 rounded-lg dark:bg-neutral-700\"><label for=\"toggle-count-monthly\" class=\"relative inline-block py-2 px-3\"><span class=\"inline-block relative z-10 text-sm font-medium text-gray-800 cursor-pointer dark:text-neutral-200\">Monthly</span> <input id=\"toggle-count-monthly\" name=\"frequency\" value=\"monthly\" type=\"radio\" class=\"absolute top-0 end-0 size-full border-transparent bg-transparent bg-none text-transparent rounded-lg cursor-pointer before:absolute before:inset-0 before:size-full before:rounded-lg focus:ring-offset-0 checked:before:bg-white checked:before:shadow-sm checked:bg-none focus:ring-transparent dark:checked:before:bg-neutral-800 dark:focus:ring-offset-transparent\" checked=\"\"></label> <label for=\"toggle-count-annual\" class=\"relative inline-block py-2 px-3\"><span class=\"inline-block relative z-10 text-sm font-medium text-gray-800 cursor-pointer dark:text-neutral-200\">Yearly</span> <input id=\"toggle-count-annual\" name=\"frequency\" value=\"yearly\" type=\"radio\" class=\"absolute top-0 end-0 size-full border-transparent bg-transparent bg-none text-transparent rounded-lg cursor-pointer before:absolute before:inset-0 before:size-full before:rounded-lg focus:ring-offset-0 checked:before:bg-white checked:before:shadow-sm checked:bg-none focus:ring-transparent dark:checked:before:bg-neutral-800 dark:focus:ring-offset-transparent\"></label></div></div><a class=\"text-xs hidden font-semibold text-teal-700 hover:underline hover:underline-offset-2\" href=\"#\">Split</a></div><div class=\"flex justify-end mt-5\"><button type=\"submit\" class=\"rounded-md bg-teal-800 hover:bg-teal-700 active:bg-teal-800 active:text-slate-200 py-1.5 px-2.5 text-xs font-semibold text-slate-200 hover:text-slate-100\">Create</button></div></div></form></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
