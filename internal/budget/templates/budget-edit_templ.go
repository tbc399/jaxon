// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "fmt"
import "jaxon.app/jaxon/internal/budget/models/budgets"
import "jaxon.app/jaxon/internal/transaction/models"

func BudgetDetail(budgetView *budgets.BudgetView, transactions []models.TransactionView) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div id=\"main-content\" hx-swap-oob=\"outerHTML:#main-content\" class=\"flex flex-col w-full justify-center items-center h-full pt-8 px-32\"><div class=\"flex flex-row justify-start w-full mb-4\"><div class=\"flex flex-row items-center gap-x-2 w-full dark:text-slate-500 dark:hover:text-slate-400\"><svg viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" class=\"size-4 fill-none stroke-2 stroke-current\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18\"></path></svg> <a hx-get=\"/budgets/partial\" hx-push-url=\"/budgets\" hx-target=\"#main-content\" class=\"text-sm hover:cursor-pointer\">Back to budgets</a></div></div><div class=\"flex flex-row justify-start mb-4 w-full py-2\"><h1 class=\"font-medium text-xl dark:text-white\">budgetView.CategoryName</h1></div><div class=\"w-full p-1.5 rounded-2xl dark:bg-gray-800/50 ring-1 dark:ring-slate-700/50\"><form hx-put=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(string(templ.URL(fmt.Sprintf("/budgets/%s", budgetView.Id))))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/budget/templates/budget-edit.templ`, Line: 28, Col: 87}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "\" hx-swap=\"none\" class=\"w-full rounded-xl dark:bg-gray-900\"><div class=\"flex flex-col gap-y-5 w-full h-1/2 p-4\"><div class=\"flex flex-row justify-between mb-4\"><svg viewBox=\"0 0 24 24\" class=\"size-4 stroke-slate-500 stroke-[2.5px] fill-current\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M6 18 18 6M6 6l12 12\"></path></svg></div><div class=\"w-full flex justify-between gap-y-2\"><div class=\"flex flex-col\"><span class=\"text-slate-300 text-sm font-semibold\">Category</span><div class=\"text-sm font-normal text-slate-500\">What is this budget for?</div></div></div><div class=\"w-full flex justify-between gap-y-2\"><div class=\"flex flex-col\"><span class=\"text-slate-300 text-sm font-semibold\">Amount</span><div class=\"text-sm font-normal text-slate-500\">How much do you plan to spend?</div></div><div class=\"flex items-end\"><span class=\"text-teal-700 me-2\">$</span> <input id=\"budget-amount-edit\" type=\"text\" name=\"amount\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(budgetView.Amount))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/budget/templates/budget-edit.templ`, Line: 52, Col: 65}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\" class=\"rounded-lg bg-transparent dark:text-gray-400\" placeholder=\"Amount\"></div></div><div class=\"flex justify-end mt-5\"><button type=\"submit\" class=\"rounded-md bg-teal-800 hover:bg-teal-700 active:bg-teal-800 active:text-slate-200 py-1.5 px-2.5 text-xs font-semibold text-slate-200 hover:text-slate-100\">Create</button></div></div></form></div><!-- Transaction list --><div class=\"w-full mt-5\"><div class=\"flex flex-row justify-start w-full dark:text-slate-200\"><h1>Transactions</h1></div><div class=\"w-full\"></div></div><!-- Transaction list end --></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
