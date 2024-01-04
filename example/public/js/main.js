import "./mithril.min.js"
import { ProcessList } from './process-list-view.js'
import { ProcessCreateForm } from './process-create-view.js'
import { ProcessChangeForm } from './process-change-view.js'


var root = document.body

//m.mount(root, ProcessList)
m.route(root, "/list", {
    "/list": ProcessList,
    "/change": ProcessChangeForm,
    "/create": ProcessCreateForm,
})