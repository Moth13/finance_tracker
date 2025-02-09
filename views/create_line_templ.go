// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.793
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func CreateLine() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"fixed inset-0 flex items-center justify-center bg-gray-800 bg-opacity-50\" id=\"modal\"><div class=\"bg-white p-6 rounded-lg shadow-lg w-96\"><form><div class=\"modal-body\"><div class=\"mb-4\"><label class=\"block text-gray-700\">Titre</label> <input type=\"text\" name=\"title\" class=\"w-full p-2 border rounded\" required></div><div class=\"mb-4\"><label class=\"block text-gray-700\">Valeur (€)</label> <input type=\"number\" step=\"0.01\" id=\"amount\" class=\"w-full p-2 border rounded\" required></div><div class=\"mb-4\"><label class=\"block text-gray-700\">Date</label> <input type=\"date\" name=\"due_date\" class=\"w-full p-2 border rounded\" required></div><div class=\"mb-4\"><label class=\"block text-gray-700\">Statut</label> <input type=\"checkbox\" name=\"checked\" value=\"true\"></div><div class=\"mb-4\"><label class=\"block text-gray-700\">Compte</label> <input type=\"text\" name=\"account_name\" class=\"w-full p-2 border rounded\" required></div><div class=\"mb-4\"><label class=\"block text-gray-700\">Mois</label> <input type=\"text\" name=\"month_name\" class=\"w-full p-2 border rounded\" required></div><div class=\"mb-4\"><label class=\"block text-gray-700\">Catégorie</label> <input type=\"text\" name=\"category_name\" class=\"w-full p-2 border rounded\" required></div><div class=\"mb-4\"><label class=\"block text-gray-700\">Description</label> <textarea name=\"description\" class=\"w-full p-2 border rounded\"></textarea></div></div><div class=\"modal-footer\"><button type=\"button\" class=\"bg-gray-500 text-white px-4 py-2 rounded\" hx-get=\"/\" hx-target=\"body\">Annuler</button> <button type=\"submit\" class=\"bg-blue-500 text-white px-4 py-2 rounded\" hx-post=\"/new_line\" hx-target=\"body\">Ajouter</button></div></form></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
