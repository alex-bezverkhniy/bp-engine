import "./mithril.min.js"
import { Process } from './process-service.js'
import { NavigationMenu } from './navigation.js'

export const ProcessList = {
    oninit: Process.loadList,
    view: function () {
        let formatter = new Intl.DateTimeFormat('en');
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
                        m("th", { scope: "col" }, "#"),
                        m("td", m("a.uuid", p.uuid)),
                        m("td.datetime", `${created_at.toLocaleDateString()} ${created_at.toLocaleTimeString()}`),
                        m("td", p.statuses[currentStatusIndex].name),
                    ])
                }))
                /* m("tbody", [ 
                    m("tr", [
                        m("th", {scope:"col"}, "#"),
                        m("td",  "1fd93604-54f5-4bb3-bbea-fcaa8c074c81"),
                        m("td",  "created"),
                    ])
                ]) */
            ]),
            )
        ]
    }
}