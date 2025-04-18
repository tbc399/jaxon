// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "jaxon.app/jaxon/internal/account/models/accounts"

func UploadPage(accounts []models.Account) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div _=\"on closeUploadModal remove me\"><!-- modal underlay --><div id=\"upload-file-modal-backdrop\" style=\"z-index: 79;\" class=\"hidden transition duration fixed inset-0 dark:bg-slate-800 bg-opacity-0 dark:bg-opacity-40 hs-overlay-backdrop\" _=\"on click trigger closeUploadModal\"></div><!-- end modal underlay --><!-- transaction upload modal --><div id=\"transaction-upload-file-modal\" class=\"hidden hs-overlay open size-full fixed top-0 start-0 z-[80] overflow-x-hidden overflow-y-auto pointer-events-none\" aria-hidden=\"false\" tabindex=\"-1\"><div id=\"main-content\" hx-swap-oob=\"outerHTML:#main-content\" class=\"flex flex-col w-full items-center justify-center h-full pt-8 px-32\"><div hx-get=\"/transactions\" hx-target=\"#main-content\" hx-push-url=\"/transactions\"><span>Back to transactions</span></div><div class=\"flex flex-row justify-start w-full py-2 text-lg font-normal text-gray-100\"><h1>Upload transactions</h1></div><form method=\"post\" action=\"/transactions/upload\" enctype=\"multipart/form-data\"><div class=\"flex flex-col pointer-events-auto max-w-full max-h-full h-full sm:max-w-lg sm:max-h-none sm:h-auto sm:border sm:rounded-xl sm:shadow-md sm:shadow-base-200 dark:bg-gray-900 sm:dark:border-gray-800\"><div class=\"flex flex-col py-3 px-4\"><div class=\"flex justify-end mb-2\"><button type=\"button\" class=\"flex items-center disabled:opacity-50 disabled:pointer-events-none\" data-hs-overlay=\"#trnsaction-upload-file-modal\" _=\"on click trigger closeUploadModal\"><svg class=\"flex-shrink-0 size-5 dark:text-slate-600 dark:hover:text-slate-500\" xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\"><path d=\"M18 6 6 18\"></path> <path d=\"m6 6 12 12\"></path></svg></button></div><div class=\"flex items-center justify-center w-full\"><div class=\"flex flex-col items-center justify-center pt-5 pb-6\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"50\" height=\"50\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"1.5\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"text-slate-500 lucide lucide-file-spreadsheet\"><path d=\"M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z\"></path> <path d=\"M14 2v4a2 2 0 0 0 2 2h4\"></path> <path d=\"M8 13h2\"></path> <path d=\"M14 13h2\"></path> <path d=\"M8 17h2\"></path> <path d=\"M14 17h2\"></path></svg><p class=\"mb-2 mt-3 text-sm text-slate-500 dark:text-slate-400\">Click to upload a csv file</p><p class=\"text-xs text-slate-500 dark:text-slate-400\">(MAX. 5MB)</p></div></div><div class=\"px-10\"><label><input name=\"file\" type=\"file\" class=\"block w-full text-sm text-gray-500\n                  file:me-4 file:py-2 file:px-4\n                  file:rounded-lg file:border-0\n                  file:text-sm file:font-semibold\n                  file:bg-teal-700 file:hover:bg-teal-600 file:text-slate-300\n                  p-1\n                  border rounded-xl dark:border-slate-800\n                  file:disabled:opacity-50 file:disabled:pointer-events-none\n                \"></label></div><div class=\"flex items-center justify-center w-full h-full\"><select data-hs-select=\"{\n                                  &#34;placeholder&#34;: &#34;Select an account...&#34;,\n                                  &#34;toggleTag&#34;: &#34;&lt;button type=\\&#34;button\\&#34; aria-expanded=\\&#34;false\\&#34;&gt;&lt;/button&gt;&#34;,\n                                  &#34;toggleClasses&#34;: &#34;hs-select-disabled:pointer-events-none hs-select-disabled:opacity-50 relative py-3 ps-4 pe-9 flex gap-x-2 text-nowrap w-full cursor-pointer bg-white border border-gray-200 rounded-lg text-start text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-neutral-900 dark:border-neutral-700 dark:text-neutral-400 dark:focus:outline-none dark:focus:ring-1 dark:focus:ring-neutral-600&#34;,\n                                  &#34;dropdownClasses&#34;: &#34;mt-2 z-50 w-full max-h-72 p-1 space-y-0.5 bg-white border border-gray-200 rounded-lg overflow-hidden overflow-y-auto [&amp;::-webkit-scrollbar]:w-2 [&amp;::-webkit-scrollbar-thumb]:rounded-full [&amp;::-webkit-scrollbar-track]:bg-gray-100 [&amp;::-webkit-scrollbar-thumb]:bg-gray-300 dark:[&amp;::-webkit-scrollbar-track]:bg-neutral-700 dark:[&amp;::-webkit-scrollbar-thumb]:bg-neutral-500 dark:bg-neutral-900 dark:border-neutral-700&#34;,\n                                  &#34;optionClasses&#34;: &#34;py-2 px-4 w-full text-sm text-gray-800 cursor-pointer hover:bg-gray-100 rounded-lg focus:outline-none focus:bg-gray-100 hs-select-disabled:pointer-events-none hs-select-disabled:opacity-50 dark:bg-neutral-900 dark:hover:bg-neutral-800 dark:text-neutral-200 dark:focus:bg-neutral-800&#34;,\n                                  &#34;optionTemplate&#34;: &#34;&lt;div class=\\&#34;flex justify-between items-center w-full\\&#34;&gt;&lt;span data-title&gt;&lt;/span&gt;&lt;span class=\\&#34;hidden hs-selected:block\\&#34;&gt;&lt;svg class=\\&#34;shrink-0 size-3.5 text-blue-600 dark:text-blue-500 \\&#34; xmlns=\\&#34;http:.w3.org/2000/svg\\&#34; width=\\&#34;24\\&#34; height=\\&#34;24\\&#34; viewBox=\\&#34;0 0 24 24\\&#34; fill=\\&#34;none\\&#34; stroke=\\&#34;currentColor\\&#34; stroke-width=\\&#34;2\\&#34; stroke-linecap=\\&#34;round\\&#34; stroke-linejoin=\\&#34;round\\&#34;&gt;&lt;polyline points=\\&#34;20 6 9 17 4 12\\&#34;/&gt;&lt;/svg&gt;&lt;/span&gt;&lt;/div&gt;&#34;,\n                                  &#34;extraMarkup&#34;: &#34;&lt;div class=\\&#34;absolute top-1/2 end-3 -translate-y-1/2\\&#34;&gt;&lt;svg class=\\&#34;shrink-0 size-3.5 text-gray-500 dark:text-neutral-500 \\&#34; xmlns=\\&#34;http://www.w3.org/2000/svg\\&#34; width=\\&#34;24\\&#34; height=\\&#34;24\\&#34; viewBox=\\&#34;0 0 24 24\\&#34; fill=\\&#34;none\\&#34; stroke=\\&#34;currentColor\\&#34; stroke-width=\\&#34;2\\&#34; stroke-linecap=\\&#34;round\\&#34; stroke-linejoin=\\&#34;round\\&#34;&gt;&lt;path d=\\&#34;m7 15 5 5 5-5\\&#34;/&gt;&lt;path d=\\&#34;m7 9 5-5 5 5\\&#34;/&gt;&lt;/svg&gt;&lt;/div&gt;&#34;\n                                }\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, account := range accounts {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<option value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(account.Id)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/transaction/templates/transactions-upload.templ`, Line: 112, Col: 36}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var3 string
			templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(account.Name)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/transaction/templates/transactions-upload.templ`, Line: 112, Col: 53}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</select></div><div class=\"px-10 mt-2 flex justify-end dark:text-slate-500 dark:hover:text-teal-700 dark:hover:underline dark:hover:underline-offset-1 text-sm font-bold\"><a href=\"#\">More options</a></div></div><div class=\"flex justify-end items-center gap-x-2 py-3 px-4 mt-auto sm:mt-0\"><button type=\"button\" class=\"py-2 px-3 inline-flex items-center gap-x-2 text-sm font-normal text-slate-400 \n              disabled:opacity-50 disabled:pointer-events-none dark:bg-slate-900 dark:border-gray-700 \n              dark:text-white dark:hover:bg-gray-800\" _=\"on click trigger closeUploadModal\">Close</button> <button type=\"submit\" class=\"py-2 px-3 inline-flex items-center gap-x-2 text-sm font-semibold rounded-md border border-transparent bg-primary text-slate-300 hover:bg-teal-600 disabled:opacity-50 disabled:pointer-events-none\">Import</button></div></div></form></div></div><!-- end transaction upload modal --></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
