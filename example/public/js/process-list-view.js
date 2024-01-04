import "./mithril.min.js"
import { Process } from './process-service.js'
import { NavigationMenu } from './navigation.js'
import { ProcessDetails } from './process-details-view.js'

export const ProcessList = {
    oninit: Process.loadList,
    view: function () {        
        var rowNum = 1
        return [
            m("main",
            m(NavigationMenu),
            m("h1", {
                class: "title"
            }, "Processes List"),
            m("table.process-list", [
                m("thead", [
                    m("tr", [
                        m("th", { scope: "col" }, "#"),
                        m("th", { scope: "col" }, "UUID"),
                        m("th", { scope: "col" }, "Created At"),
                        m("th", { scope: "col" }, "Current Status"),
                    ])
                ]),
                
                m("tbody", Process.list.map(function (p) {
                    console.log(p)
                    var currentStatusIndex = p.statuses.length - 1
                    var created_at = new Date(p.created_at)
                    return m("tr", [
                        m(ProcessDetails, {data: p}),
                        m("th", { scope: "col" }, rowNum++),
                        m("td", m("a.uuid",
                        { 
                            href: "#",
                            onclick: function(e) {                                
                                e.preventDefault()                                
                                var dlg = document.getElementById(`${p.uuid}-dlg`)
                                console.log(dlg)
                                dlg.setAttribute("open", "true")
                                return false
                            }
                        },
                         p.uuid)),
                        m("td.datetime", `${created_at.toLocaleDateString()} ${created_at.toLocaleTimeString()}`),
                        m("td", p.statuses[currentStatusIndex].name),
                    ])
                }))
            ]),
            )
        ]
    }
}