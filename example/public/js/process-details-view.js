import "./mithril.min.js"

export const ProcessDetails = {
    view: function (vnode) {
        var p = vnode.attrs.data
        var id = `${p.uuid}-dlg`
        var currentStatusIndex = p.statuses.length - 1
        var created_at = new Date(p.created_at)
        created_at = `${created_at.toLocaleDateString()} ${created_at.toLocaleTimeString()}`
        var changed_at = new Date(p.changed_at)
        changed_at = `${changed_at.toLocaleDateString()} ${changed_at.toLocaleTimeString()}`
        return m(`dialog[id='${id}']`,
            m("article",
                [
                    m("a.close[href='#close'][aria-label='Close'][data-target='modal-example'][onclick='toggleModal(event)']",
                    ),
                    m("h3",
                        "Process Info"
                    ),
                    m("table",
                        [
                            m("thead",
                                m("tr",
                                    [
                                        m("th[scope='col']",
                                            "#"
                                        ),
                                        m("th[scope='col']",
                                            "Value"
                                        )
                                    ]
                                )
                            ),
                            m("tbody",
                                m("tr",
                                    [
                                        m("th[scope='col']", m("kbd", "UUID")),
                                        m("td", p.uuid)
                                    ]
                                ),
                                m("tr",
                                    [
                                        m("th[scope='col']", m("kbd", "Current Status")),
                                        m("td", p.statuses[currentStatusIndex].name)
                                    ]
                                ),
                                m("tr",
                                    [
                                        m("th[scope='col']", m("kbd", "Created At")),
                                        m("td", created_at)
                                    ]
                                ),
                                m("tr",
                                    [
                                        m("th[scope='col']", m("kbd", "Changed At")),
                                        m("td", changed_at)
                                    ]
                                ),
                                m("tr",
                                    [
                                        m("th[scope='col']", m("kbd", "Payload")),                                        
                                        m("td", m("pre.json-payload", JSON.stringify(p.payload, undefined, 2)))
                                    ]
                                ),
                                m("tr",
                                    [
                                        m("th[scope='col']", m("kbd", "Statuses")),
                                        m("td", m("pre.json-payload", JSON.stringify(p.statuses, undefined, 2)))
                                    ]
                                ),
                            )
                        ]
                    ),
                    m("footer",
                        [
                            m("a.secondary[href='#cancel'][role='button'][data-target='modal-example']",
                                {
                                    onclick: function (e) {
                                        e.preventDefault()
                                        var dlg = document.getElementById(id)
                                        dlg.removeAttribute("open")
                                        return false
                                    }
                                },
                                " Close "
                            ),
                        ]
                    )
                ]
            )
        )
    }
}